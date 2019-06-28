package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

// Request type
type Request struct {
	Action  string
	Command string
}

// Response type
type Response struct {
	Error  string
	Result []byte
}

func encodeRequest(w io.Writer, r Request) {
	encoder := gob.NewEncoder(w)
	encoder.Encode(r)
}

func decodeRequest(r io.Reader) Request {
	req := Request{}
	decoder := gob.NewDecoder(r)
	decoder.Decode(&req)
	return req
}

func encodeResponse(w io.Writer, r Response) {
	encoder := gob.NewEncoder(w)
	encoder.Encode(r)
}

func decodeResponse(r io.Reader) Response {
	resp := Response{}
	decoder := gob.NewDecoder(r)
	decoder.Decode(&resp)
	return resp
}

var (
	serverCommand = flag.NewFlagSet("serve", flag.ExitOnError)
	serverPort    = serverCommand.String("port", ":7891", "port to listen on")
	debug         = serverCommand.Bool("debug", false, "print debug information")
	clientCommand = flag.NewFlagSet("connect", flag.ExitOnError)
	remoteHost    = clientCommand.String("remoteHost", "localhost", "hostname/ip to connect to")
	remotePort    = clientCommand.String("remotePort", ":7891", "remote port to connect to")
	command       = clientCommand.String("command", "", "command to run at server")
)

func init() {
	if len(os.Args) == 1 {
		fmt.Println("usage: gomon <command> [<args>]")
		fmt.Println("The most common gmon commands are")
		fmt.Println("    serve    start the gomon server")
		fmt.Println("  connect    connect to the gomon server")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "serve":
		serverCommand.Parse(os.Args[2:])
	case "connect":
		clientCommand.Parse(os.Args[2:])
	}
}

func main() {
	if serverCommand.Parsed() {
		startServer(*serverPort)

	} else if clientCommand.Parsed() {
		if *command == "" {
			fmt.Println("-command requires an argument")
			os.Exit(2)
		}
		startClient(*remoteHost, *remotePort, *command)
	}
}

func startServer(port string) {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	log.Printf("server listening on %s\n", l.Addr().String())
	for {
		c, _ := l.Accept()
		go func(c net.Conn) {
			req := decodeRequest(c)
			switch req.Action {
			case "run":
				output, err := exec.Command("bash", "-c", req.Command).CombinedOutput()
				resp := Response{Result: output}
				if err != nil {
					resp.Error = err.Error()
				}
				if *debug {
					log.Printf("response object is :%p, resp size: %d\n", &resp, len(resp.Result))
				}
				encodeResponse(c, resp)
				log.Printf("%s [%s]\n", c.RemoteAddr().String(), req.Command)
				c.Close()
			default:
				c.Write([]byte("invalid action\n"))
				c.Close()
			}
		}(c)
	}
}

func startClient(host, port, command string) {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	c, err := net.Dial("tcp", host+port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	log.Printf("Connected to %s\n", c.RemoteAddr().String())
	req := Request{Action: "run", Command: command}
	encodeRequest(c, req)
	resp := decodeResponse(c)
	if resp.Error != "" {
		fmt.Println(resp.Error)
	}
	fmt.Print(string(resp.Result))
}
