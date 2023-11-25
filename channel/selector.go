package channel

import "fmt"

var hashMap = map[int]string{1: "xianglizhen", 2: "lee", 3: "july"}

func processChannel(in <-chan int, out chan<- string) {
    for inValue := range in {
        result := hashMap[inValue]
        out <- result
    }
}
func Selector() {
    sendChan := make(chan int)
    receiveChan := make(chan string)
    go processChannel(sendChan, receiveChan)
    sendChan <- 1
    res := <-receiveChan
    fmt.Println(res)
}
