package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("test")
    go func() {
        var testP *int
        var test = 0
        testP = &test
        time.AfterFunc(
            time.Second, func() {
                fmt.Println("3", *testP)
                if *testP == 0 {
                    *testP = 1
                }
                fmt.Println("4")
            })

        if *testP == 0 {
            fmt.Println("1")
            *testP = 2
        }
        fmt.Println("before stop test:", *testP)
        //timer.Stop()
        fmt.Println("after stop test:", *testP)
    }()
    time.Sleep(time.Second * 3)
}
