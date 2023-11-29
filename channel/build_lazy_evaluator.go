package channel

import "fmt"

type Any interface{}

type EvalFunc func(Any) (Any, Any)

func BuildLazyEvaluator(evalFunc EvalFunc, initState Any) func() Any {
    retValChan := make(chan Any)
    // 惰性生成器之所以惰性生成，也是因为后台起了一个协程来不停的计算下一个参数的返回值和状态参数
    loopFunc := func() {
        var actState = initState
        var retVal Any
        for {
            retVal, actState = evalFunc(actState)
            retValChan <- retVal
        }
    }
    go loopFunc()
    retFunc := func() Any {
        return <-retValChan
    }

    return retFunc
}

func BuildLazyIntEvaluator(evalFunc EvalFunc, initState Any) func() int {
    ef := BuildLazyEvaluator(evalFunc, initState)
    return func() int {
        return ef().(int)
    }
}

func TestBuildLazyEvaluator() {
    // 工厂函数三个原则：空接口，闭包，高阶函数
    // 工厂函数(BuildLazyIntEvaluator)需要一个函数和一个初始状态作为输入参数，返回一个无参、返回值是生成序列的函数
    // 传入的函数(evalFunc)需要计算出下一个返回值以及下一个状态参数
    evalFunc := func(state Any) (Any, Any) {
        os := state.(int)
        ns := os + 2
        return os, ns
    }
    ef := BuildLazyIntEvaluator(evalFunc, 0)
    for i := 0; i < 10; i++ {
        fmt.Println("even:", ef())
    }

}
