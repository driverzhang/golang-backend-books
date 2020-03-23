# gin实现操作日志中间件

> 在大部分后台管理中都会出现操作日志模块的需求，一般情况你会想到创建一个操作日志表 opLog 然后就是在每个 api 的地放调用插入数据即可。

如上所述，这样的设计确实十分灵活可以让每个写api的开发人员自由的操作内容，但是也比较费时间和精力。

## 解决答方案

利用gin框架的use挂起中间件，将操作日志进行统一抽象的被动调用。

先给出如下的表结构字段设计，

这里我是放在mongodb中的：

```go
// 后台管理员操作日志表
type OperationLog struct {
    Id            primitive.ObjectID `form:"_id" json:"_id" bson:"_id"`            // 操作id
    OperationId   int64              `json:"operationId" bson:"operationId"`       // 操作员Id
    OperationName string             `json:"operation_name" bson:"operation_name"` // 操作员name
    Ip            string             `json:"ip" bson:"ip"`                         // 操作地址Ip
    Operation     string             `json:"operation" bson:"operation"`           // 操作 增删改
    Business      string             `json:"business" bson:"business"`             // 业务名称
    Tables        string             `json:"tables" bson:"tables"`                 // 直接相关表（一个）
    Data          interface{}        `json:"data" bson:"data"`                     // 新数据的保存 只保留主键id
    CreatedAt     int64              `json:"created_at" bson:"created_at"`
}
```

具体如下：

* 全部获取都需要 ctx \*gin.Context
* 利用 获取 userId 方法获取 operationId \(一般都是从你的token产生出提取userId即可\)
* ip := ctx.ClientIP\(\) 直接调用获取客户端ip地址
* ctx.Request.Method 获取api实际的方法对应 新增，编辑，删除操作即可
* 开发任务只需要维护一个表面+业务 和 HandlerFunc 映射数据即可

> 根据每个api 调用的 HandlerFunc 都是唯一不同的，即可获取 Tables、Business

* 最后将代理抽象利用 ctx.Set\(\) 和 ctx.Get\(\) 获取 相关数据 id 即可。

```go
// 用与 handle 函数与 表名及其业务描述的映射
type HandleTableAndDesc map[string]string // value 值填入 {最直接关联的唯一表名},{业务描述}

var HandleTableName = HandleTableAndDesc{
    "InsertWithdrawalOrder": "withdrawal,生成提现订单",
    "UpdateStatus":          "withdrawal,提现操作状态",
    "TaskShelf":             "task_base,任务上架",
    "UploadIcon":            "photo,上传任务图片",
    "StaffFreeze":           "staffs,冻结解冻管理员",
    "UserFreeze":            "cend_user,冻结解冻用户",
}
```

## 具体代码封装：

```go
var Engine = gin.Default()

func init() {
    swagger.RegSwagger(Engine)
    Engine.Use(operationlog.Operation(contracts.HandleTableName))

    Engine.POST("/stat/withdrawal", controllers.InsertWithdrawalOrder)
    Engine.PUT("/stat/withdrawal", controllers.UpdateStatus)          
    Engine.GET("/stat/withdrawal", controllers.GetWithdrawalList)  
}
```

具体的 operationlog 包 是核心代码如下：

```go
func Operation(handlerTableName map[string]string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "" {
            c.Next() // 注意 next()方法的作用是跳过该调用链去直接后面的中间件以及api路由
        }
        // 1. 获取当前用户ID
        ok, userId := xxx.GetUserInfo(c)
        if !ok {
            userId = ""
        }
        // 2. 获取表名 这里会把前面 "加载的包名.HandlerName"
        handlerName := strings.Split(c.HandlerName(), ".")[1]
        if _, ok := handlerTableName[handlerName]; !ok {
            c.Next() // 如果映射中没有找到表关系，说明没有手动加入
        }

        tableDesc := strings.Split(handlerTableName[handlerName], ",")
        if len(tableDesc) != 2 {
            return // 存在映射，但是不完全则返回
        }

        tableName := tableDesc[0] // 0表名
        desc := tableDesc[1]      // 1业务描述
        // 3. 获取Ip
        ip := c.ClientIP()
        opMethod := &OperationMethod{
            C:      c,
            Table:  tableName,
            Desc:   desc,
            Ip:     ip,
            UserId: userId,
        }
        // 4. 根据method方法名确定操作行为
        switch c.Request.Method {
        case "PUT":
            err := opMethod.putOperation()
            if err != nil {
                c.Next()
            }
        case "POST":
            err := opMethod.postOperation()
            if err != nil {
                c.Next()
            }
        case "DELETE":
            err := opMethod.deleteOperation()
            if err != nil {
                c.Next()
            }
        default:
            return
        }
    }
}

type OperationMethod struct {
    C      *gin.Context
    Table  string `json:"table"`
    Desc   string `json:"desc"`
    Ip     string `json:"ip"`
    UserId string `json:"userId"`
}

type bodyId struct {
    Id interface{} `json:"_id" bson:"_id"`
}

func (o *OperationMethod) putOperation() (err error) {
    // 更新操作要拿到操作id的话注意一定要统计使用的业务Id 对应 mongo 为 _id
    // 获取body中的数据
    var bodyBytes []byte
    if o.C.Request.Body != nil {
        bodyBytes, _ = ioutil.ReadAll(o.C.Request.Body)
    }
    o.C.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 返回body中的原值

    bodyId := &bodyId{}
    var putId interface{}
    _ = json.Unmarshal(bodyBytes, bodyId)
    if bodyId.Id == nil {
        o.C.Next()
        putId, _ = o.C.Get("updateId")
    } else {
        putId = bodyId.Id
    }

    // 插入操作日志表中
    userId, _ := strconv.Atoi(o.UserId)
    staff := entities.Staff{}
    staffData, err := staff.GetNameById(uint64(userId))
    if err != nil {
        return err
    }

    operationLog := entities.OperationLog{
        OperationId:   int64(userId),
        Operation:     "PUT",
        OperationName: staffData.Name,
        Ip:            o.Ip,
        Business:      o.Desc,
        Tables:        o.Table,
        Data:          putId,
    }

    err = operationLog.InsertLog()
    if err != nil {
        return err
    }

    return
}

func (o *OperationMethod) postOperation() (err error) {
    // 新建操作 需要记录新插入的数据
    // 获取插入数据的Id
    o.C.Next()                         // 直接跳过调用链 让其开始走api
    insertId, _ := o.C.Get("insertId") // 若不存在 也应该先插入操作记录 只是业务操作Insert没有返回id
    userId, _ := strconv.Atoi(o.UserId)
    staff := entities.Staff{}
    staffData, err := staff.GetNameById(uint64(userId))
    if err != nil {
        return
    }

    operationLog := entities.OperationLog{
        OperationId:   int64(userId),
        Operation:     "POST",
        OperationName: staffData.Name,
        Ip:            o.Ip,
        Business:      o.Desc,
        Tables:        o.Table,
        Data:          insertId,
    }

    err = operationLog.InsertLog()
    if err != nil {
        return
    }

    return
}

func (o *OperationMethod) deleteOperation() (err error) {
    mongo.Collection(o.Table)
    // 删除操作 需要对删除的旧数据进行记录
    var ids []string
    if len(o.C.QueryArray("ids")) != 0 {
        ids = o.C.QueryArray("ids")
        ids = make([]string, len(ids))
    } else if len(o.C.QueryArray("_id")) != 0 {
        ids = o.C.QueryArray("_id")
        ids = make([]string, len(ids))
    }

    idsAll := "" // 对应敏感业务 软删除可以查询到删除的具体数据
    for _, v := range ids {
        if len(ids) == 1 {
            idsAll = idsAll + v
        } else {
            idsAll = idsAll + v + ","
        }
    }

    userId, _ := strconv.Atoi(o.UserId)
    staff := entities.Staff{}
    staffData, err := staff.GetNameById(uint64(userId))
    if err != nil {
        return
    }

    operationLog := entities.OperationLog{
        OperationId:   int64(userId),
        Operation:     "DELETE",
        OperationName: staffData.Name,
        Ip:            o.Ip,
        Business:      o.Desc,
        Tables:        o.Table,
        Data:          idsAll,
    }

    err = operationLog.InsertLog()
    if err != nil {
        return
    }

    return
}
```

为了获取 更新和新增的相关id，需要在调用具体api时，返回对应的id。然后利用gin.Context 调用 ctx.Set\("insertId", id\) 进行存储，这样方便中间件调用ctx.Get\("insertId"\)，进行获取数据即可。

