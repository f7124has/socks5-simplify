package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var (
	localAddr  = flag.String("l", "", "host:port to listen on (required).")
	remoteAddr = flag.String("r", "", "host:port to forward to (required).")
	silentRun  = flag.Bool("s", false, "run silently, do not stdout anything.")
)

func forward(conn net.Conn) {
	client, err := net.Dial("tcp", *remoteAddr)
	if err != nil {
		log.Printf("Dial failed: %v", err)
		defer conn.Close()
		return
	}
	log.Printf("Forwarding from %v to %v\n", conn.LocalAddr(), client.RemoteAddr())

	// write data to server (auth)
	_, err = client.Write([]byte{5, 2, 0, 1})
	if err != nil {
		defer client.Close()
		defer conn.Close()
		return
	}

	// Read real socks5 server out
	tmp := make([]byte, 8)
	_, err = client.Read(tmp)
	if err != nil {
		defer client.Close()
		defer conn.Close()
		return
	}

	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func main() {
	flag.Usage = func() {
		fmt.Printf("tcpforward - This thing ate all my forwarded packets!\n\n")
		fmt.Printf("Usage: \n\ttcpforward [-s] -l local_address:port -r remote_address:port\n")
		fmt.Printf("\t\t-l host:port\tto listen on (required).\n")
		fmt.Printf("\t\t-r host:port\tto forward to (required).\n")
		fmt.Printf("\t\t-s \t\trun silently, do not stdout anything.\n\n")
	}

	flag.Parse()

	if *localAddr == "" || *remoteAddr == "" {
		flag.Usage()
		os.Exit(255)
	}

	log.SetPrefix("tcpforward: ")

	if *silentRun == true {
		log.SetOutput(ioutil.Discard)
	}

	listener, err := net.Listen("tcp", *localAddr)
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr().String())
		go forward(conn)
	}
}
