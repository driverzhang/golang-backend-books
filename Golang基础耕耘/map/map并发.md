# goroutine 与 map 并发的采坑事件

## 1. goroutine 与map 的并发读写操作

> 在Go 1.6之前， 内置的map类型是部分goroutine安全的，并发的读没有问题，并发的写可能有问题。自go 1.6之后， 并发地读写map会报错，这在一些知名的开源库中都存在这个问题，所以go 1.9之前的解决方案是额外绑定一个锁，封装成一个新的struct或者单独使用锁都可以。

> 因为map为引用类型，所以即使函数传值调用，参数副本依然指向映射m, 所以多个goroutine并发写同一个映射m， 写过多线程程序的同学都知道，对于共享变量，资源，并发读写会产生竞争的， 故共享资源遭到破坏


---

### 1. 有并发问题的map

官方的[Why are map operations not defined to be atomic? ](https://golang.org/doc/faq#atomic_maps)]已经提到内建的map不是线程(goroutine)安全的。


> ... and in those cases where it did, the map was probably part of some larger data structure or computation that was already synchronized. 

我们来看一下代码吧：

一个goroutine一直读，一个goroutine一只写同一个键值，即即使读写的键不相同，而且map也没有"扩容"等操作，代码还是会报错。

```go

func main() {

	m := make(map[int]int)
	go func() {
		for {
			_ = m[1]
		}
	}()
	go func() {
		for {
			m[2] = 2
		}
	}()
	select {}
}


```

然后你会发现 运行不起：

    fatal error: concurrent map read and map write



有时候数据竞争不是很容易发现，你可以输入：

    go run --race main.go

进行查看。


---


### 2. Go 1.9之前的解决方案

你可以用互斥锁 sync.Mutex 也可以用读写锁 sync.RWMutex(性能好些)。

```go

var counter = struct{
    sync.RWMutex
    m map[string]int
}{m: make(map[string]int)}


```

读数据的时候很方便的加锁：

```go

counter.RLock()
n := counter.m["some_key"]
counter.RUnlock()
fmt.Println("some_key:", n)


```

写数据的时候:


```go

counter.Lock()
counter.m["some_key"]++
counter.Unlock()

```

当然你也可以单独使用 读写锁进行读写加锁操作，只是如果有多个的情况下就没有嵌入结构体那么方便操作了。


---


### 3. 终极解决方案 sync.Map

> 可以说，上面的解决方案相当简洁，并且利用读写锁而不是Mutex可以进一步减少读写的时候因为锁带来的性能。


- Store
- LoadOrStore
- Load
- Delete
- Range


Store(key, value interface{})

> 存 key,value 存储一个设置的键值。

LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
>  返回键的现有值(如果存在)，否则存储并返回给定的值，如果是读取则返回true，如果是存储返回false。

Load(key interface{}) (value interface{}, ok bool)
> 读取存储在map中的值，如果没有值，则返回nil。OK的结果表示是否在map中找到值。


Delete(key interface{})
> 删除key,及其value

Range(f func(key, value interface{}) bool)
> 循环读取map中的值.遍历所有的key,value

```go

package main
 
import (
	"fmt"
	"sync"
)
 
func main() {
	var m sync.Map
 
	//Store
	m.Store(1,"a")
	m.Store(2,"b")
 
	//LoadOrStore
	//若key不存在，则存入key和value，返回false和输入的value
	v,ok := m.LoadOrStore("1","aaa")
	fmt.Println(ok,v) //false aaa
 
	//若key已存在，则返回true和key对应的value，不会修改原来的value
	v,ok = m.LoadOrStore(1,"aaa")
	fmt.Println(ok,v) //false aaa
 
	//Load
	v,ok = m.Load(1)
	if ok{
		fmt.Println("it's an existing key,value is ",v)
	} else {
		fmt.Println("it's an unknown key")
	}
 
	//Range
	//遍历sync.Map, 要求输入一个func作为参数
	f := func(k, v interface{}) bool {
		//这个函数的入参、出参的类型都已经固定，不能修改
		//可以在函数体内编写自己的代码，调用map中的k,v
 
			fmt.Println(k,v)
			return true
		}
	m.Range(f)
 
	//Delete
	m.Delete(1)
	fmt.Println(m.Load(1))
 
}

```



关于 其源码分析 [可参考此文章](https://colobu.com/2017/07/11/dive-into-sync-Map/)