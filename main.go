package main

import (
    "fmt"
    "strconv"
    "sync"
)

type Task struct {
    TaskID   int
    TaskName string
}

func generateTask(n int) chan *Task {
    task := make(chan *Task)
    go func(n int) {
        for i := 0; i < n; i++ {
            task <- &Task{
                TaskID:   i,
                TaskName: "",
            }
        }
    }(n)

    return task
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

func main() {
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
