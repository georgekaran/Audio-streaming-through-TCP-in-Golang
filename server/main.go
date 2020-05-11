package main

import (
	"encoding/binary"
	"github.com/gordonklaus/portaudio"
	"log"
	"net"
)

const sampleRate = 44100
const seconds = 0.1

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	buffer := make([]float32, sampleRate * seconds)
	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(buffer), func(in []float32) {
		for i := range buffer {
			buffer[i] = in[i]
		}
	})
	must(err)
	must(stream.Start())
	defer stream.Close()

	listen, errNet := net.Listen("tcp", ":8080")
	if errNet != nil {
		log.Fatal(errNet)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn, buffer)
	}
}

func handle(con net.Conn, buffer []float32) {
	defer con.Close()
	binary.Write(con, binary.BigEndian, &buffer)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
