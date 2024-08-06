package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Buffer struct {
	buf []byte
	len int
}

func (b *Buffer) Reader() io.Reader {
	return b
}

func (b *Buffer) Read(p []byte) (int, error) {
	n := copy(p, b.buf)
	return n, nil
}

func handler(conn net.Conn) {
	data := make([]byte, 0)
	recv := make([]byte, 4096)
	length := 0

	// for {
	// read to the tmp var
	n, err := conn.Read(recv)
	if err != nil {
		// log if not normal error
		if err != io.EOF {
			fmt.Printf("Read error - %s\n", err)
		}
		return
	}

	fmt.Printf("Read - %d\n", n)

	// append read data to full data
	data = append(data, recv[:n]...)

	// update total read var
	length += n
	// }

	buffer := Buffer{
		buf: data,
		len: length,
	}

	fmt.Printf("data :\n%s\n", data)
	buf := bufio.NewReaderSize(buffer.Reader(), length)

	req, err := http.ReadRequest(buf)
	if err != nil {
		fmt.Printf("err : %s\n", err)
		return
	}

	fmt.Printf("req : %+v\n", req)
}

func main() {
	l, err := net.Listen("tcp", ":33114")
	if err != nil {
		fmt.Println("Failed to Listen : ", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to Accept : ", err)
			continue
		}

		go handler(conn)
	}
}
