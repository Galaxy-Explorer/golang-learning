package channel

import (
    "fmt"
    "strconv"
    "sync"
)

type Task1 struct {
    TaskID   int
    TaskName string
}

type Pool1 struct {
    lock  sync.Mutex
    Tasks []*Task1
}

func process1(t *Task1) {
    go func(t *Task1) {
        t.TaskName = "XiangliZhen" + strconv.Itoa(t.TaskID)
    }(t)
}

func Worker(p *Pool1) {
    task := &Task1{}
    for {
        if len(p.Tasks) != 0 {
            p.lock.Lock()
            task = p.Tasks[0]
            p.Tasks = p.Tasks[1:]
            p.lock.Unlock()
            process1(task)
        } else {
            break
        }
    }
}

func testTask1() {
    tasks := []*Task1{
        {
            TaskID:   1,
            TaskName: "",
        },
        {
            TaskID:   2,
            TaskName: "",
        },
        {
            TaskID:   3,
            TaskName: "",
        },
    }
    pool := &Pool1{
        lock:  sync.Mutex{},
        Tasks: tasks,
    }
    Worker(pool)

    for i, task := range tasks {
        fmt.Printf("%d: %+v\n", i, task)
    }
}
