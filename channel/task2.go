package channel

import (
    "fmt"
    "strconv"
)

type task2 struct {
    taskID   int
    taskName string
}

func sendWork2(n int) chan *task2 {
    chTask := make(chan *task2)
    go func(n int) {
        defer close(chTask)
        for i := 0; i < n; i++ {
            chTask <- &task2{
                taskID:   i,
                taskName: "",
            }
        }
    }(n)

    return chTask
}

func process2(t *task2) {
    t.taskName = "XiangliZhen" + strconv.Itoa(t.taskID)
}

func worker2(in, out chan *task2) {
    go func(in, out chan *task2) {
        for {
            // 对比打印素数，为啥这块可以检测到通道关闭
            t, ok := <-in
            if ok {
                process2(t)
                out <- t
            } else {
                close(out)
                break
            }
        }
    }(in, out)
}

func testTask2() {
    pending := sendWork2(100)
    done := make(chan *task2)
    worker2(pending, done)
    for t := range done {
        fmt.Println(t.taskID, t.taskName)
    }
}
