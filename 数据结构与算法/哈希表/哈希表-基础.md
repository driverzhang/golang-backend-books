# 哈希表-基础

> 哈希表又可被称作散列表或者 Hash Table

## 哈希表与数组：

#### 散列表用的是数组支持按照下标随机访问数据的特性，所以散列表其实就是数组的一种扩展，由数组演化而来。可以说，如果没有数组，就没有散列表。

具体的对应 Go语言中的数据结构就是 map

而哈希表最重要的核心就是 哈希函数 与 哈希碰撞

## 1. 哈希函数

![](https://raw.githubusercontent.com/driverzhang/image-storage/master/20200321185537.png)

哈希函数，首先它是一种函数，可以设为 hash\(key\) 其中key表示我们的元素键值，而 hash\(key\) 表示经过哈希函数计算出的结果值。

简单的说 将一串数值编号 利用 哈希函数 计算或映射出对应的 一个索引值。

而被计算出来的 该索引值 被称为 哈希值。

哈希函数应该能够将不同键能够地映射到不同的索引上，这要求哈希函数输出范围大于输入范围，但是由于键的数量会远远大于映射的范围，所以在实际使用时，这个理想的结果是不可能实现的。

一般在哈希函数计算出来的值要求尽量均匀的，结果不均匀的哈希函数会造成更多的冲突并导致更差的读写性能。

但是基本无法避免这种冲突被称为 哈希冲突。

### 构造哈希函数三点基本要求：

* 散列函数计算得到的散列值是一个非负整数；
* 如果 key1 = key2，那 hash\(key1\) == hash\(key2\)；
* 如果 key1 ≠ key2，那 hash\(key1\) ≠ hash\(key2\)。

第三点理解起来可能会有问题，这个要求看起来合情合理，但是在真实的情况下，要想找到一个不同的 key 对应的散列值都不一样的散列函数，几乎是不可能的。即便像业界著名的 [MD5](https://zh.wikipedia.org/wiki/MD5)、 [SHA](https://zh.wikipedia.org/wiki/SHA%E5%AE%B6%E6%97%8F)、 [CRC](https://zh.wikipedia.org/wiki/%E5%BE%AA%E7%92%B0%E5%86%97%E9%A4%98%E6%A0%A1%E9%A9%97) 等哈希算法，也无法完全避免这种 哈希冲突问题。

而且正常情况下数组的存储空间是有限的，所以早晚会出现重叠冲突的时候。

## 2. 哈希冲突

#### 开放寻址方法

按照线性探测的去插入：

![](https://static001.geekbang.org/resource/image/5c/d5/5c31a3127cbc00f0c63409bbe1fbd0d5.jpg)

按照线性探测的去查找：

![](https://static001.geekbang.org/resource/image/91/ff/9126b0d33476777e7371b96e676e90ff.jpg)

#### 链表法

链表法是一种更加常用的散列冲突解决办法，相比开放寻址法，它要简单很多。

在散列表中，每个“桶（bucket）”或者“槽（slot）”会对应一条链表，所有散列值相同的元素我们都放到相同槽位对应的链表中。

![](https://static001.geekbang.org/resource/image/a4/7f/a4b77d593e4cb76acb2b0689294ec17f.jpg)

### 参看链接：

[散列表（上）：Word文档中的单词拼写检查功能是如何实现的？](https://time.geekbang.org/column/article/64233)

[Go中的哈希表](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-hashmap/)

