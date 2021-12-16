# Golang

### 关键字
#### new和make区别
    make关键字的作用是创建slice、map和channel等内置的数据结构， new的作用是为类型申请一片内存空间，并返回指向这片内存的指针。

#### defer顺序
    FILO后进先出

#### cap和len的区别
    cap表示最大容量，len表示当前长度。

### channel
#### 对空channel进行读写会发生什么
    阻塞当前goroutine。
#### 对已经关闭对channel进行读写会发生什么
    如果channel中还有未读取完的值，则读取剩下的值。否则返回默认零值。
#### 如何判断channel已经关闭
    通过返回的第二个字段(true or false)判断
#### 怎么循环读取channel的值
    通过for range语法读取
#### channel 有缓冲 与 无缓冲 的区别，主要用于什么场景
    提高并发
#### 如何判断chan是否满了
    select或者使用cap判断
#### 未初始化的channel
    读取未初始化的channel都会发生阻塞
#### 如何同时监听多个channel,或者说多路复用
    select监听多个channel

### map
#### 如何实现集合，value采用什么类型，为什么？
    可以使用布尔类型，不过会占用一定的内存空间。最好使用空结构体，因为空结构体不占用内存。
#### 如何实现顺序读实现
    使用切片获取key,排序后再进行读取。
#### 对一个map进行遍历，元素的顺序一样吗？
    不一样



### slice 
#### 怎么声明空数组，有什么区别
```go
package main
import "fmt"

func main()  {
   a := [...]int{}
   b := [0]int{}
   var c [0]int
   var d = [...]int{}
   fmt.Println(a == b, b == c, c == d) // true true true
}
```
#### 数组和切片的区别
    数组是值类型，切片是引用类型。切片可以动态伸缩。
#### 切片append会发生内存重新分配吗
    append后长度超过切片的cap会重新分配内存，分配大小为原来两倍的内存。
#### 截取切片的一部分赋给一个新的切片，修改这个新的切片其中的一个元素，那么原切片会被修改吗？对新的切片进行append操作，原切片会发生什么？
```go
// append不导致原切片cap发生改变，则会改变原切片值
package main

import "fmt"

func main()  {
   a := make([]int, 10)
   for i := 0; i < 5; i++ {
      a[i] = i+1
   }
   b := a[2:4]
   b[0] = 10
   fmt.Println(cap(b))
   b = append(b, 9,29)
   fmt.Println(a) // [1 2 10 4 9 29 0 0 0 0]
   fmt.Println(b) // [10 4 9 29]
}
// append 导致原切cap发生改变，新切片指向新数组
//a := make([]int, 5)
//for i := 0; i < 5; i++ {
//a[i] = i+1
//}
//b := a[2:4]
//b[0] = 10
//fmt.Println(cap(b))
//b = append(b, 9,29)
//fmt.Println(a) // [1 2 10 4 5]
//fmt.Println(b) // [10 4 9 29]
```

### string
#### 单引号、双引号、反引号有什么区别
    单引号表示byte或者rune类型， 双引号表示字符串类型支持转义序列， 反引号表示字符串字面量，不支持转义序列。
#### string 和 []byte 进行转换时会发生内存拷贝吗？
    不会， string类型底层包含[]byte类型。
#### 字符串拼接相关，内存拷贝问题，如何高效拼接字符串

### struct
#### 空结构体
    空结构体不占用内存，一般可以用来实现集合、控制协程并发和仅包含方法的结构体。
#### 结构体方法采用指针和非指针区别
    golang参数传递都是值传递方式, 要改变结构体内部的值需要传指针。对于接口来说，*T包含的方法集合包括接收者是*T和T的方法，而T包含的方法集合只包含接收者是T的方法

### 协程
#### GMP模型以及数量关系
[Golang调度器GMP原理与调度全分析](https://learnku.com/articles/41728)
1. GMP模型
![GMP模型](./statics/GMP模型.jpeg)
    - G: goroutine协程
    - P: processor处理器(包含了G)
    - M: 系统级线程, M关联一个KSE实体(即操作系统内核可调度的最小单位)
2. 调度器策略
![G调度流程](./statics/G调度流程.jpeg)
    - 复用线程: 避免频繁创建、销毁线程。
    - work stealing机制  
         当本线程无可运行的 G 时，尝试从其他线程绑定的 P 偷取 G，而不是销毁线程。
    - hand off机制  
         当本线程因为 G 进行系统调用阻塞时，线程释放绑定的 P，把 P 转移给其他空闲的线程执行。
#### 一个协程只会在一个processor上运行吗
      不会
#### 如何准确等待所有协程的结束
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			fmt.Println(a)
		}(i)
	}
	wg.Wait()
}
```

#### 如何等待任意一个任务返回
[等待任意一个任务返回源码](./fcfs/fcfs.go)
#### 生产者消费者模型
```go

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Producer(factor int, out chan <- int) {
	for i := 0; ; i++ {
		out <- factor * i
	}
}
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	// 简单的生产者/消费者模型
	// 成果队列
	ch := make(chan int, 10)

	// 两个生产者
	go Producer(3, ch)
	go Producer(5, ch)

	// 一个消费者
	go Consumer(ch)

	// ctrl + c退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit by (%v)\n", <-sig)
}
```
#### 如何控制多个协程(任务)同时结束
#### 单例模式
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// 单例模式
	var once sync.Once

	onceBody := func() {
		fmt.Println("Only once")
	}

	done := make(chan struct{}, 10)

	for i := 0; i < cap(done); i++ {
		once.Do(onceBody)
		done <- struct{}{}
	}

	for i := 0; i < cap(done); i++ {

	}
}
```
### 发布订阅模型
[发布者订阅者模型代码](./pubsub/pubsub.go)
```go
package main

import (
	"fmt"
	"interview/pubsub"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	p := pubsub.NewPublisher(10, 100*time.Millisecond)
	defer p.Close()
	all := p.Subscribe()
	filtrate := p.SubscribeTopic(func(v interface{}) bool {
		if num, ok := v.(int); ok {
			return num > 10
		}
		return false
	})
	go p.Publish(12)
	go p.Publish(4)

	go func() {
		for v := range all {
			fmt.Printf("all: %v\n", v)
		}
	}()

	go func() {
		for v := range filtrate {
			fmt.Printf("more than 10 : %v\n", v)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("quit", <-sig)
}
```
### GC
    - 标记-清除法  
        措施：stw - 标记 - 清除 - stw
        缺点：gc期间暂停业务运行，可能会出现业务卡顿现象
    - 三色标记法 + 堆插入屏障（满足强三色不变式） + 堆删除屏障（满足弱三色不变式）
        不使用STW，此时如果出现黑色对象引用白色对象， 而白色对象的上游灰色对象被删除；会出现对象被误删除。  
        解决方式：满足强/弱三色不变式之一， 堆上添加对象时，使用插入屏障（添加对象标记为灰色， 也就不会出现黑色对象引用白色对象的情况）， 堆上删除对象时，使用删除屏障（如果被删除对象时白色或者灰色， 标记该对象为灰色）；  
        缺点是：由于栈上没有使用添加/删除屏障， 灰色标记队列为空时， 回收对象之前， 使用rescan栈上的对象， 然后在进行回收；另外删除屏障会导致回收精度低， 会出现延迟回收；
    - 三色标记法 + 混合写屏障机制（满足变形的弱三色不变式）  
      具体操作：
        - GC开始时将栈上的对象全部扫描并标记为黑色（栈无需二次扫描）
        - GC期间， 任何在栈上创建的新对象都标记为黑色
        - 被删除的对象标记为灰色
        - 被添加的对象标记为灰色  
      优点：栈上无需二次扫描， 对象的回收是延迟回收。
### 逃逸分析
    编译期做逃逸分析，首先可以减少垃圾回收的压力，变量如果没有逃逸到heap上，内存分配相对更加容易，函数运行结束的同时直接回收变量资源，  
    其次逃逸分析过后可以知道那些变量可以分配在栈上，栈上分配更快，性能更好，  
    再者可以做同步消除，定义变量的函数如果有同步锁，而运行时只有一个线程访问，逃逸分析后的机器码，会去掉同步锁。
### HTTPS原理分析
[https原理分析](https://juejin.cn/post/6844903830916694030)

### Cookie、Session、Token
[jwt跨域认证解决方案](https://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html)
### 操作系统I/O模型
- 同步阻塞I/O：
  1. 进程启动 IO 的 read 调用开始，用户线程进入阻塞状态；
  2. 系统内核收到调用开始准备数据，这时候，数据还没有到达内核缓冲区，这时候内核就要等待；
  3. 内核一直等到完整的数据后，就会将数据从内核缓冲区复制到用户的进程缓冲区，然后等内核返回结果； 
  4. 当内核返回后，用户线程才会解除阻塞状态，重新开始运行；
- 同步非阻塞I/O：
  1. 发起一个非阻塞 socket 的 read 读操作的系统调用，在内核数据没有准备好的情况下，用户线程发起 IO 请求时，立即返回。为了读取到最终的数据，用户线程需要不断地发起 IO 系统调用；
  2. 内核数据到达后，用户线程发起系统调用，用户线程阻塞。内核开始复制数据，它会将数据从内核缓冲区复制到进程缓冲区，然后返回结果；
  3. 用户线程读取到数据后，才会接触阻塞状态，重新运行。所以，用户进程需要经过多次尝试，才能保证最终真正读取到数据，然后继续执行；
- 同步多路复用I/O：
  1. 首先是选择器注册。将需要 read 操作的目标 socket 网络连接，提前注册到 select/epoll 选择器中，然后才可以开启整个 IO 多路复用模型的轮询；
  2. 通过选择器的查询方法，查询所有注册过的 socket 连接的就绪状态，通过该查询，内核返回一个就绪的 socket 列表。当任何一个 socket 中的数据准备好了，内核缓冲区就有数据了，内核就将该 socket 加入到就绪的列表中。注：当用户进程调用了 select 查询方法，那么整个线程就会被阻塞；
  3. 用户线程获取了就绪的 socket 列表后，根据其中的 socket 连接发起 read 系统调用，用户线程阻塞，内核开始将数据从内核缓冲区复制到用户缓冲区；
  4. 复制完成后，内核返回结果，用户线程解除阻塞状态，用户线程读到数据，继续执行；
- 同步信号驱动I/O：
- 异步I/O：
  1. 用户线程发起 read 调用后，立即可以做其他的了，用户线程不阻塞；
  2. 内核开始准备数据，等数据准备好了，内核就将数据从内核缓冲区复制到用户缓冲区；
  3. 完成后，内核会给用户线程发送一个信号，或者调用户线程注册的回调接口，通知用户线程 read 操作已完成；
  4. 用户线程读取数据，继续后续操作；