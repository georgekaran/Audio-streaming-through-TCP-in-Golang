package main

import (
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"log"
	"net"
	"os"
	"os/exec"
)

const sampleRate = 88200
const seconds = 0.06

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

	clearTerminal()
	fmt.Println("Running Server TCP on port :8080")

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

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
