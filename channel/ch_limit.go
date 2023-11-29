package channel

import (
    "fmt"
    "time"
)

const timeout = 6

type taskWithCh struct {
    id          int
    executeTime int
    timeout     int
    result      chan string // 此处必须要用缓冲为1的管道来接受值
}

func RunWithCh(task *taskWithCh) {
    defer close(task.result)
    chRun := make(chan string)
    go runWithCh(task.id, task.executeTime, chRun)
    select {
    case re := <-chRun:
        task.result <- re
    case <-time.After(time.Duration(task.timeout) * time.Second):
        re := fmt.Sprintf("task id %d , timeout", task.id)
        task.result <- re
    }
}
func runWithCh(taskID, executeTime int, chRun chan string) {
    defer close(chRun)
    time.Sleep(time.Duration(executeTime) * time.Second)
    chRun <- fmt.Sprintf("task id %d , sleep %d second", taskID, executeTime)
    return
}
func testChLimit() {
    tasks := []*taskWithCh{
        {
            id:          1,
            executeTime: 2,
            timeout:     timeout,
            result:      make(chan string, 1),
        },
        {
            id:          2,
            executeTime: 4,
            timeout:     timeout,
            result:      make(chan string, 1),
        },
        {
            id:          3,
            executeTime: 5,
            timeout:     timeout,
            result:      make(chan string, 1),
        },
    }
    chLimit := make(chan bool, 1)
    limitFunc := func(chLimit chan bool, task *taskWithCh) {
        // 如果result使用了无缓冲管道，会导致Run阻塞在此，等待从<-task.result读出数据
        // chLimit写入的数据，在此无法被接受，如果len(tasks) == chLimit，不会形成死锁，chLimit < len(tasks), 形成死锁
        RunWithCh(task)
        <-chLimit
    }
    startTime := time.Now()
    fmt.Println("Multirun start")
    for _, task := range tasks {
        chLimit <- true
        go limitFunc(chLimit, task)
    }
    for _, task := range tasks {
        fmt.Println(<-task.result)
    }
    close(chLimit)

    fmt.Printf("Multissh finished. Process time %s. Number of task is %d", time.Now().Sub(startTime), len(tasks))
}
