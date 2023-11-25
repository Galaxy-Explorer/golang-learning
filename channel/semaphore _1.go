package channel

import (
    "fmt"
)

var N1 = 10

func doSomething1(i int, f float64) float64 {
    return float64(i) + f
}
func Semaphore1() {
    data := make([]float64, N1)
    res := make([]float64, N1)
    sem := make(chan struct{}, N1)

    data = []float64{0.0, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9}

    for i, datum := range data {
        go func(i int, datum float64) {

            res[i] = doSomething1(i, datum)
            sem <- struct{}{}
        }(i, datum)
    }

    for i := 0; i < N1; i++ {
        <-sem
    }
    close(sem)
    fmt.Println(res)

}
