package channel

import (
    "fmt"
    "strconv"
)

type task struct {
    taskID   int
    taskName string
}

func sendWork(n int) chan *task {
    chTask := make(chan *task)
    go func(n int) {
        defer close(chTask)
        for i := 0; i < n; i++ {
            chTask <- &task{
                taskID:   i,
                taskName: "",
            }
        }
    }(n)

    return chTask
}

func process(t *task) {
    t.taskName = "XiangliZhen" + strconv.Itoa(t.taskID)
}

func worker(in, out chan *task) {
    go func(in, out chan *task) {
        for {
            // 对比打印素数，为啥这块可以检测到通道关闭
            t, ok := <-in
            if ok {
                process(t)
                out <- t
            } else {
                close(out)
                break
            }
        }
    }(in, out)
}

//var hashMap = map[int]string{1: "xianglizhen", 2: "lee", 3: "july"}
//
//func processChannel(in <-chan int, out chan<- string) {
//    for inValue := range in {
//        result := hashMap[inValue]
//        out <- result
//    }
//}
//func Selector() {
//    sendChan := make(chan int)
//    receiveChan := make(chan string)
//    go processChannel(sendChan, receiveChan)
//    sendChan <- 1
//    res := <-receiveChan
//    fmt.Println(res)
//}

func testTask2() {
    pending := sendWork(100)
    done := make(chan *task)
    worker(pending, done)
    for t := range done {
        fmt.Println(t.taskID, t.taskName)
    }
}
