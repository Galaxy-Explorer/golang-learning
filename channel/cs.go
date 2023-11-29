package channel

import (
    "fmt"
)

type Request struct {
    a, b   int
    replyc chan int // reply channel inside the Request
}

type binOp func(a, b int) int

func run(op binOp, req *Request) {
    req.replyc <- op(req.a, req.b)
}

func server(op binOp, service chan *Request, quit chan bool) {
    for {
        select {
        case req := <-service:
            go run(op, req)
        case <-quit:
            return // 跳出当前循环

        }
    }
}

func startServer(op binOp) (service chan *Request, quit chan bool) {
    // 抛出了一个接受者，用于接受request请求
    service = make(chan *Request)
    quit = make(chan bool)
    go server(op, service, quit)
    return service, quit
}

func testCS() {
    // adder 用于接受request请求，并起了一个server groutine
    adder, quit := startServer(func(a, b int) int { return a + b })

    const N = 100
    var reqs [N]Request

    for i := 0; i < N; i++ {
        req := &reqs[i]
        req.a = i
        req.b = i + N
        req.replyc = make(chan int)
        adder <- req // adder is a channel of requests，server抛出了一个addr管道用于接受请求，
    }
    // checks:
    for i := 0; i < N; i++ {
        res := <-reqs[i].replyc
        fmt.Println("Request: ", res, "is ok!")
    }
    quit <- true
    fmt.Println("done")
}
