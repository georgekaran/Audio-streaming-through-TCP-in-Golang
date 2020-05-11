package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"io/ioutil"
	"net"
	"time"
)

const sampleRate = 44100
const seconds = 0.1

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	buffer := make([]float32, sampleRate * seconds)

	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, len(buffer), func(out []float32) {
		go readFromServer(out, buffer)
	})
	must(err)
	must(stream.Start())

	for {
		time.Sleep(time.Millisecond)
	}
}

func readFromServer(out []float32, buffer []float32) {
	conn := dialServer()
	defer conn.Close()

	bs, _ := ioutil.ReadAll(conn)
	bytesReader := bytes.NewReader(bs)
	binary.Read(bytesReader, binary.BigEndian, &buffer)
	for i := range out {
		out[i] = buffer[i]
	}
}

func dialServer() net.Conn {
	conn, errConn := net.Dial("tcp", "localhost:8080")
	for errConn != nil {
		conn, errConn = net.Dial("tcp", "localhost:8080")
		fmt.Println("Trying to reconnect...")
		time.Sleep(time.Second)
	}
	return conn
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}