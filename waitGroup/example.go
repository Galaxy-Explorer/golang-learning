package waitGroup

import (
    "fmt"
    "sync"
)

func testWaitGroup() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            //wg.Add(1) ，看似add 与 done 是成对出现的，但是存在问题如下：
            //有可能出现 Wait 方法先于 Add 方法执行，此时由于计数器值为 0，Wait 方法会被直接放行
            defer wg.Done()
            fmt.Println("i=", i)
        }(i)
    }
    wg.Wait()
}
