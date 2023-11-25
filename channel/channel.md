# panic的三种情况

1. 关闭已经关闭的channel
2. 关闭未初始化的channel
3. 向已经关闭的channel写入数据

# 阻塞总结

1. 无缓冲管道
    1. 数据要发送，没有接收者
    2. 数据要接受，没有发送者

2. 有缓冲管道
    1. 数据要发送，没有接收者
    2. 数据要接受，没有发送者
    3. channel容量已满，发送者阻塞在写入
    4. channel容量为空，接受者阻塞在读取

# 阻塞而产生死锁
* 接受者与发送者没有一一对应，数据发送，没有接收者
```go
func main() {
    c := make(chan int)
    c <- 1
    fmt.Println(<-c)
}
```
* 接受者与发送者没有一一对应，数据接受，没有发送者
```go
func main() {
    c := make(chan int)
    fmt.Println(<-c)
    c <- 1
}
```
* 向未初始化nil channel中写入数据阻塞在`fmt.Println("look this 1")`，从而导致读nil管道引起死锁
```go
func main() {
	var ch chan int

	go func() {
		fmt.Println("look this 1")
		ch <- 10
		fmt.Println("look this 2")
	}()
	time.Sleep(time.Second)
	r := <-ch
	fmt.Println(r)
	fmt.Println()
}
```

* 向未初始化nil channel中读取数据阻塞在`fmt.Println("look this 1")`，从而导致写nil管道引起死锁
```go
func main() {
	var ch chan int

	go func() {
		fmt.Println("look this 1")
		r := <-ch
		fmt.Println("look this 2", r)
	}()
	time.Sleep(time.Second)
	ch <- 1
	fmt.Println()
}
```

* 没有发送者，但是还要从管道读取数据，引发阻塞，发生死锁
```go
func main() {
	var ch = make(chan int)
	go func() {
		// defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	for {
		x, ok := <-ch
		if ok {
			fmt.Println("接收到的值：", x)
		} else {
			break
		}
	}
}
```

* 发送者goroutine已经关闭，接受者goroutine发生阻塞
```go
func main() {
	var ch = make(chan int)
	go func() {
		defer fmt.Println("end")
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	go func() {
		for {
			x := <-ch
			fmt.Println("接收到的值：", x)
		}
	}()

	time.Sleep(5 * time.Second)
}
```



* 有缓冲区，但是缓冲区已满，发送数据阻塞，引起死锁
```go
func main() {
    c := make(chan int ,5)
    for i := 0; i < 10; i++ {
       c <- i
    }
}
```

* 有缓冲区，但是缓冲区为空，读取数据阻塞，引起死锁
```go
func main() {
	c := make(chan int, 5)
	recv := <-c
	fmt.Println(recv)
}
```

* 读取一个关闭的管道，不会引发死锁，而是返回该类型的0值
```go
func main() {
	var ch = make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	for {
		x, ok := <-ch
		if ok {
			fmt.Println("接收到的值：", x)
		} else {
			fmt.Println("end:", ok)
			break
		}
	}
}
```

# 管道使用总结

1. 无缓冲管道，接受者与发送者要一一对应，不能放在一个goroutine
2. 如果在一个goroutine中在启动了子goroutine，并通过管道通信，要先启动该子goroutine（可能是发送者或者接受），否则该goroutine会因为没有发送者（发送者）会导致死锁
3. 在发送者处关闭管道，在接受者处做好对管道关闭的判断

# 管道应用场景
## 信号量模式
### 有缓冲通道的实现
```go
package main

import (
	"fmt"
)

var N = 10

func doSomething(i int, f float64) float64 {
	return float64(i) + f
}
func main() {
	data := make([]float64, N)
	res := make([]float64, N)
	sem := make(chan struct{}, N)

	data = []float64{0.0, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9}

	for i, datum := range data {
		go func(i int, datum float64) {

			res[i] = doSomething(i, datum)
			sem <- struct{}{}
		}(i, datum)
	}

	for i := 0; i < N; i++ {
		<-sem
	}
	close(sem)
	fmt.Println(res)

}
```

### 无缓冲通道的实现
```go
package main

import (
	"fmt"
)

var N = 10

func doSomething(i int, f float64) float64 {
	return float64(i) + f
}
func main() {
	data := make([]float64, N)
	res := make([]float64, N)
	sem := make(chan struct{})

	data = []float64{0.0, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9}

	for i, datum := range data {
		go func(i int, datum float64) {

			res[i] = doSomething(i, datum)
			sem <- struct{}{}
		}(i, datum)
		<-sem
	}
	close(sem)
	fmt.Println(res)
}
```

### 实现互斥锁
```go
package main

import (
	"fmt"
)

type semaphore chan struct{}

func (s semaphore) lock() {
	s <- struct{}{}
}

func (s semaphore) unlock() {
	<-s
}

func main() {
	test := 0
	sem := make(semaphore)
	for i := 0; i < 10; i++ {
		go func(i int) {
			sem.lock()
			test = i
		}(i)
		sem.unlock()
	}
	close(sem)
	fmt.Println(test)
}

```

### 管道工厂
```go
package main

import (
	"fmt"
	"time"
)

var ch chan int

func pump() chan int {
	ch = make(chan int)
	go func() {
		defer close(ch)
		for i := 0; ; i++ {
			ch <- i
		}
	}()
	return ch
}

func suck(ch chan int) {
	for {
		receiver, ok := <-ch
		if ok {
			fmt.Println("success:", receiver)
		} else {
			break
		}
	}
    // 下面这种写法也可以，一直读到管道关闭
	//for v := range ch {
	//	fmt.Println("success:", v)
	//}
}

func main() {
	stream := pump()
	go suck(stream)
	time.Sleep(time.Millisecond)
}
```


### 通道迭代模式
```go
package main

import "fmt"

type item int
type container struct {
	items []item
}

func (c *container) Iter() <-chan item {
	ch := make(chan item)
	go func() {
		defer close(ch)
		for i := 0; i < len(c.items); i++ {
			ch <- c.items[i]
		}
	}()
	return ch
}

func main() {
	C := &container{
		items: []item{1, 2, 3, 4, 5},
	}
	for x := range C.Iter() {
		fmt.Println("success:", x)
	}
}

```

### 生产消费者模式
```go
package main

import "fmt"

type request struct {
	id   int
	body string
}

var ch = make(chan *request)

func Consume(r chan *request) {
	for x := range r {
		fmt.Println(x.id, x.body)
	}
}

func Produce() chan *request {
	r := &request{}
	go func(r *request) {
		defer close(ch)
		for i := 0; i < 10; i++ {
			r = &request{
				id:   i,
				body: "XiangliZhen" + string(i),
			}
			ch <- r
		}
	}(r)
	return ch
}

func main() {
	Consume(Produce())
}

```

### 管道选择器模式
#### 根据管道输入值，输出对应的值
GO指导的例子（很赞）
```go
package main

import "fmt"

var hashMap = map[int]string{1: "xianglizhen", 2: "lee", 3: "july"}

func processChannel(in <-chan int, out chan<- string) {
	for inValue := range in {
		result := hashMap[inValue]
		out <- result
	}
}
func main() {
	sendChan := make(chan int)
	receiveChan := make(chan string)
	go processChannel(sendChan, receiveChan)
	sendChan <- 1
	res := <-receiveChan
	fmt.Println(res)
}
```

#### 打印素数1
```go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main
import (
    "fmt"
)
// Send the sequence 2, 3, 4, ... to returned channel
func generate() chan int {
    ch := make(chan int)
    go func() {
        for i := 2; ; i++ {
            ch <- i
        }
    }()
    return ch
}
// Filter out input values divisible by 'prime', send rest to returned channel
func filter(in chan int, prime int) chan int {
    out := make(chan int)
    go func() {
        for {
            if i := <-in; i%prime != 0 {
                out <- i
            }
        }
    }()
    return out
}
func sieve() chan int {
    out := make(chan int)
    go func() {
        ch := generate()
        for {
            prime := <-ch
            ch = filter(ch, prime)
            out <- prime
        }
    }()
    return out
}
func main() {
    primes := sieve()
    for {
        fmt.Println(<-primes)
    }
}
```

#### 打印素数1
做了一点优点，在sieve传入MAX，输出小于MAX的所有素数
```go
package main

import "fmt"

// Send the sequence 2, 3, 4, ... to returned channel
func generate() chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 2; ; i++ {
            ch <- i
        }
    }()
    return ch
}
// Filter out input values divisible by 'prime', send rest to returned channel
func filter(in chan int, prime int) chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            if i := <-in; i%prime != 0 {
                out <- i
            }
        }
    }()
    return out
}
func sieve(n int) chan int {
    out := make(chan int)
    stop := false
    go func(stop bool) {
        defer close(out)
        ch := generate()
        for !stop {
            prime := <-ch
            if prime > n {
                stop = true
            }
            ch = filter(ch, prime)
            if prime < n {
                out <- prime
            }

        }
    }(stop)
    return out
}
func main() {
    primes := sieve(100)
    for x := range primes {
        fmt.Println(x)
    }
}
```