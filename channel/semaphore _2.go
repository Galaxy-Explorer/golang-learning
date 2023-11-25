package channel

import (
    "fmt"
)

var N2 = 10

func doSomething2(i int, f float64) float64 {
    return float64(i) + f
}
func Semaphpre2() {
    data := make([]float64, N2)
    res := make([]float64, N2)
    sem := make(chan struct{})

    data = []float64{0.0, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9}

    for i, datum := range data {
        go func(i int, datum float64) {

            res[i] = doSomething2(i, datum)
            sem <- struct{}{}
        }(i, datum)
        <-sem
    }
    close(sem)
    fmt.Println(res)
}
