package channel

import (
    "fmt"
    "strconv"
    "sync"
)

type Task struct {
    TaskID   int
    TaskName string
}

type Pool struct {
    lock  sync.Mutex
    Tasks []*Task
}

func process(t *Task) {
    go func(t *Task) {
        t.TaskName = "XiangliZhen" + strconv.Itoa(t.TaskID)
    }(t)
}

func Worker(p *Pool) {
    task := &Task{}
    for {
        if len(p.Tasks) != 0 {
            p.lock.Lock()
            task = p.Tasks[0]
            p.Tasks = p.Tasks[1:]
            p.lock.Unlock()
            process(task)
        } else {
            break
        }
    }
}

func testTask1() {
    tasks := []*Task{
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
    pool := &Pool{
        lock:  sync.Mutex{},
        Tasks: tasks,
    }
    Worker(pool)

    for i, task := range tasks {
        fmt.Printf("%d: %+v\n", i, task)
    }
}
