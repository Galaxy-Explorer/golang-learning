package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    test := func() {}
    fmt.Println(test)
    ctx, test := context.WithTimeout(context.Background(), time.Second*10)
    fmt.Println(ctx.Deadline())
    fmt.Println(test)
}
