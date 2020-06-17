# Factory method pattern

> 创建型模式 之 工厂方法模式

工厂方法模式（Factory method pattern）使用的频率非常高，在我们日常的开发中总能见到它的身影。

其定义为：

Define an interface for creating an object,but let subclasses decide which class to
instantiate.Factory Method lets a class defer instantiation to subclasses.（

定义一个用于创建对象的
接口，让子类决定实例化哪一个类。工厂方法使一个类的实例化延迟到其子类。）


## 工厂方法模式的优点

- 首先，良好的封装性，代码结构清晰。

- 工厂方法模式的扩展性非常优秀。

- 屏蔽产品类。这一特点非常重要，产品类的实现如何变化，调用者都不需要关
  心，它只需要关心产品的接口，只要接口保持不变，系统中的上层模块就不要发生变化。
  
- 工厂方法模式是典型的解耦框架。高层模块值需要知道产品的抽象类，其他的实
  现类都不用关心