package testhelper

import (
    "bytes"
    "fmt"
    "net"
    "strconv"
)

type ServerReq struct {
    ID       int64
    Operands []int
    Operator string
}

type ServerResp struct {
    ID      int64
    Formula string
    Result  int
    errMsg  error
}

func op(operands []int, operator string) int {
    var result int
    switch {
    case operator == "+":
        for _, operand := range operands {
            if result == 0 {
                result = operand
            } else {
                result += operand
            }
        }
    case operator == "-":
        for _, operand := range operands {
            if result == 0 {
                result = operand
            } else {
                result -= operand
            }
        }
    case operator == "*":
        for _, operand := range operands {
            if result == 0 {
                result = operand
            } else {
                result *= operand
            }
        }
    case operator == "/":
        for _, operand := range operands {
            if result == 0 {
                result = operand
            } else {
                result /= operand
            }
        }
    }
    return result
}

func genFormula(operands []int, operator string, result int, equal bool) string {
    var buff bytes.Buffer
    n := len(operands)
    for i := 0; i < n; i++ {
        if i > 0 {
            buff.WriteString(" ")
            buff.WriteString(operator)
            buff.WriteString(" ")
        }
        buff.WriteString(strconv.Itoa(operands[i]))
    }
    if equal {
        buff.WriteString("=")
    } else {
        buff.WriteString("!=")
    }
    buff.WriteString(strconv.Itoa(result))
    return buff.String()
}

func reqHandler(conn net.Conn) {
    var errMsg string
    var sresp ServerResp
    req, err := read(conn, DELIM)
    if err != nil {
        errMsg = fmt.Sprintf("Server: Req Read Error: %s", err)
    } else {

    }
}
