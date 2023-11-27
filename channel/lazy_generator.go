package channel

import (
    "fmt"
)

var resume = make(chan int)

func integers() chan int {
    yield := make(chan int)
    go func(count int) {
        for {
            yield <- count
            count++
        }

    }(0)
    return yield
}

func generateInteger() int {
    return <-resume
}

func testLazyGen() {
    resume = integers()
    fmt.Println(generateInteger()) //=> 0
    fmt.Println(generateInteger()) //=> 1
    fmt.Println(generateInteger()) //=> 2
    fmt.Println(generateInteger()) //=> 2
}
