package channel

import "fmt"

type item int
type container struct {
    items []item
}

func (c *container) Iter() <-chan item {
    ch := make(chan item)
    go func() {
        defer close(ch)
        for i := 0; i < len(c.items); i++ {
            ch <- c.items[i]
        }
    }()
    return ch
}

func Iter() {
    C := &container{
        items: []item{1, 2, 3, 4, 5},
    }
    for x := range C.Iter() {
        fmt.Println("success:", x)
    }
}
