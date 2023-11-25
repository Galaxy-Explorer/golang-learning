package channel

import "fmt"

type request struct {
    id   int
    body string
}

var order = make(chan *request)

func Consume(r chan *request) {
    for x := range r {
        fmt.Println(x.id, x.body)
    }
}

func Produce() chan *request {
    r := &request{}
    go func(r *request) {
        defer close(order)
        for i := 0; i < 10; i++ {
            r = &request{
                id:   i,
                body: "XiangliZhen" + string(i),
            }
            order <- r
        }
    }(r)
    return order
}

func Consumer() {
    Consume(Produce())
}
