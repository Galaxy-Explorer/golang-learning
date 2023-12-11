package testhelper

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "net"
    "strconv"
    "sync/atomic"
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
    Err     error
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
        var sreq ServerReq
        err := json.Unmarshal(req, &sreq)
        if err != nil {
            errMsg = fmt.Sprintf("Server: Req Unmarshal Error: %s", err)
        } else {
            sresp.ID = sreq.ID
            sresp.Result = op(sreq.Operands, sreq.Operator)
            sresp.Formula = genFormula(sreq.Operands, sreq.Operator, sresp.Result, true)
        }
    }

    if errMsg != "" {
        sresp.Err = errors.New(errMsg)
    }

    bytes, err := json.Marshal(sresp)

    if err != nil {
        logger.Errorf("Server: Resp Marshal Error: %s", err)
    }

    _, err = write(conn, bytes, DELIM)
    if err != nil {
        logger.Errorf("Server: Resp Write error: %s", err)
    }
}

type TCPServer struct {
    listener net.Listener
    active   uint32
}

func NewTCPServer() *TCPServer {
    return &TCPServer{}
}

func (server *TCPServer) init(addr string) error {
    if !atomic.CompareAndSwapUint32(&server.active, 0, 1) {
        return nil
    }
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        atomic.StoreUint32(&server.active, 0)
    }
    server.listener = ln

    return nil
}

func (server *TCPServer) Listen(addr string) error {
    err := server.init(addr)
    if err != nil {
        return err
    }

    go func() {
        for {
            if atomic.LoadUint32(&server.active) != 1 {
                break
            }
            conn, err := server.listener.Accept()
            if err != nil {
                if atomic.LoadUint32(&server.active) == 1 {
                    logger.Errorf("Server: Request Acception Error: %s\n", err)
                } else {
                    logger.Warnf("Server: Broken acception because of closed network connection.")
                }
                continue
            }

            go reqHandler(conn)
        }
    }()
    return nil
}

func (server *TCPServer) Close() bool {
    if !atomic.CompareAndSwapUint32(&server.active, 1, 0) {
        return false
    }
    server.listener.Close()
    return true
}
