package testhelper

import (
    "bufio"
    "net"
)

func read(conn net.Conn, delim byte) ([]byte, error) {
    reader := bufio.NewReader(conn)
    content, err := reader.ReadBytes(delim)
    if err != nil {
        return nil, err
    }

    return content, nil
}

func write(conn net.Conn, content []byte, delim byte) (int, error) {
    writer := bufio.NewWriter(conn)
    n, err := writer.Write(content)
    if err == nil {
        writer.WriteByte(delim)
    }
    if err == nil {
        err = writer.Flush()
    }
    return n, nil
}
