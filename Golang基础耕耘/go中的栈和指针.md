## 首先go中的函数的参数都是 “值传递”。

那么如果在参数中加入 指针，就成了共享变量，但其中具体又是如何传递与赋值的呢？

函数栈调用完毕后是否需要再清理对应的栈内存呢？


- Functions execute within the scope of frame boundaries that provide an individual memory space for each respective function.

- When a function is called, there is a transition that takes place between two frames.

- The benefit of passing data “by value” is readability.

- The stack is important because it provides the physical memory space for the frame boundaries that are given to each individual function.

- All stack memory below the active frame is invalid but memory from the active frame and above is valid.

- Making a function call means the goroutine needs to frame a new section of memory on the stack.

- It’s during each function call, when the frame is taken, that the stack memory for that frame is wiped clean.

- Pointers serve one purpose, to share a value with a function so the function can read and write to that value even though the value does not exist directly inside its own frame.

- For every type that is declared, either by you or the language itself, you get for free a compliment pointer type you can use for sharing.

- The pointer variable allows indirect memory access outside of the function’s frame that is using it.

- Pointer variables are not special because they are variables like any other variable. They have a memory allocation and they hold a value.


1. 函数运行在为其分配的单独帧空间内

2. 当发生函数调用的时候，会发生两个帧之间的转移

3. 传值的优势是可读性

4. 栈是非常重要的，因为它为每个函数的帧边界提供了运行的空间

5. 所有的栈中，其下方是无效的，上方是有效的

6. 调用一个函数意味着协程需要在栈中分配一块新的空间

7. 指针存在的目的就是为了用于和函数共享值，所以函数可以操作位于其帧之外的变量

8. 声明任何一种变量，不管是自定义的还是 go 中自带的，都有一种对应的指针类型

9. 指针变量允许函数操作位于其帧外的内存

10. 指针变量并不特殊，和其他变量一样，也有内存分配，也含有一个值


## 以上结论 可以直接观看 参考文章如下：

> 文章是我精选出来的，这里就不copy了，直接看链接即可。

### 翻译版本：
### [go 中的栈和指针](https://juejin.im/post/5e9193116fb9a03c840d66ef#heading-3)


### 英文原文版：（William KennedyMay 18, 2017）
### [Language Mechanics On Stacks And Pointers](https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-stacks-and-pointers.html)


--- 

看完文章就直接来几张图就很清晰了：

### 无指针的参数版本：

![](https://www.ardanlabs.com/images/goinggo/80_figure1.png)

![](https://www.ardanlabs.com/images/goinggo/80_figure2.png)

![](https://www.ardanlabs.com/images/goinggo/80_figure3.png)

![](https://www.ardanlabs.com/images/goinggo/80_figure4.png)

### 带指针的参数版本：

![](https://www.ardanlabs.com/images/goinggo/80_figure5.png)

![](https://www.ardanlabs.com/images/goinggo/80_figure6.png)


