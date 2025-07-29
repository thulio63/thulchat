package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/chzyer/readline"
)

// Server ...
type Server struct {
	host string
	port string
}

type Client struct {
	Client net.Conn
}

// Config ...
type ServerConfig struct {
	Host string
	Port string
}

// modify host and port selection via repl, start go run at end
func (cfg *config)New() {
	fmt.Println("Please enter the desired server host:")

	//readline for repl commands
	rl, err := readline.New("- ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	
	
	hostStr, err := rl.Readline()
	if err!= nil {
		fmt.Println("error reading input:", err)
	}
	
	fmt.Println("Select which port to open for connection:")
	//add logic to avoid crucial port overlap
	
	portStr, err := rl.Readline()
	if err!= nil {
		fmt.Println("error reading input:", err)
	}
	serv := Server{
		host: hostStr,
		port: portStr,
	}
	cfg.servers_active = append(cfg.servers_active, &serv)
	*cfg.servers_count += 1

	ch := make(chan bool)
	go serv.Run(ch)
	
	fmt.Println("Started the server")
	<- ch
}

func (server *Server) Run(ch chan bool) {
	hp := fmt.Sprintf("%s:%s", server.host, server.port)
	listener, err := net.Listen("tcp", hp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on", hp)
	ch <- true
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go Respond(conn)
	}
}

func Respond(c net.Conn) {
	io.Copy(c, c)
	c.Close()
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}