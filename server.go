package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/chzyer/readline"
	"github.com/thulio63/thulchat/internal/database"
)

type Server struct {
	host string
	port string
	context context.Context
	serv *database.Server
}

type Client struct {
	Client net.Conn
}

// modify host and port selection via repl, start go run at end
func (cfg *config)New() {
	//readline for repl commands
	rl, err := readline.New("- ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	cfg.colorCon.prompt.Println("Would you like to host the server on your machine? (Y/n)")
	resp, err := rl.Readline()
	if err != nil {
		cfg.colorCon.err.Println("error reading input:", err)
	}
	resp = Clean(resp)
	hostStr := cfg.MyIP.String()
	//add check for y/n, give prompt again, for loop

	if resp != "y" {
		cfg.colorCon.prompt.Println("Please enter the desired server host:")
		hostStr, err = rl.Readline()
		if err != nil {
			cfg.colorCon.err.Println("error reading input:", err)
		}
		hostStr = Clean(hostStr)
	}	
	
	cfg.colorCon.prompt.Println("Select which port to open for connection:")
	//add logic to avoid crucial port overlap
	
	portStr, err := rl.Readline()
	if err!= nil {
		cfg.colorCon.err.Println("error reading input:", err)
	}
	portStr = Clean(portStr)
	ctx := context.Background()
	
	server, err := cfg.db.CreateServer(cfg.ctx, database.CreateServerParams{
		CreatorID: cfg.User.UserID,
		Hostname: hostStr,
		Port: portStr,
	})
	if err != nil {
		cfg.colorCon.err.Println("error creating server:", err)
	}
	serv := Server{
		host: hostStr,
		port: portStr,
		context: ctx,
		serv: &server,
	}
	//handle error for invalid server address
	cfg.servers_active = append(cfg.servers_active, &serv)

	ch := make(chan bool)
	go serv.Run(ch)
	
	cfg.colorCon.success.Println("Started the server")
	<- ch
}

func (server *Server) Run(ch chan bool) {
	hp := fmt.Sprintf("%s:%s", server.host, server.port)
	listener, err := net.Listen("tcp", hp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on", hp)
	fmt.Println("")
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