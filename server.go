package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/thulio63/thulchat/internal/database"
)

type Server struct {
	host string
	port string
	context context.Context
	serv *database.Server
	clients *map[uuid.UUID]Client
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

	cliList := make(map[uuid.UUID]Client)
	serv := Server{
		host: hostStr,
		port: portStr,
		context: ctx,
		serv: &server,
		clients: &cliList,
	}
	//handle error for invalid server address
	cfg.servers_active = append(cfg.servers_active, &serv)

	ch := make(chan bool)
	go serv.Run(ch)
	
	cfg.colorCon.success.Println("Started the server")
	<- ch

}

// need to return channel AND send same channel to the Respond function to bridge to client and the server via channel *******************************
func (server *Server) Run(ch chan bool) {
	mu := sync.Mutex{}
	
	//connect to database
	dbURL := "postgres://andrewthul:@localhost:5432/thulchat?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	defer db.Close()
	dbQ := database.New(db)

	hp := fmt.Sprintf("%s:%s", server.host, server.port)
	listener, err := net.Listen("tcp", hp)
	if err != nil {
		fmt.Println("error listening:", err)
	}
	fmt.Println("Listening on", hp)
	fmt.Println("")
	ch <- true
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection:", err)
		}
		//newClient := Client{Client: conn}
		//edit 

		go server.Respond(conn, *dbQ, &mu)
	}
}

func (server *Server) Respond(c net.Conn, db database.Queries, mu *sync.Mutex) {
	for {
		//channel here to block until data received
		//server.comm <-
		//actually i think the listener.accept does this already

		var buf bytes.Buffer
		_, err := buf.ReadFrom(c) //will this require io.readall?
		if err != nil {
			fmt.Println("server error reading from connection:", err)
		}

		if buf.String() == "death" {
			//fmt.Println("exiting server")
			break
		}

		var incDat Message
		err = json.Unmarshal(buf.Bytes(), &incDat)
		if err != nil {
			fmt.Println("server error reading sent data:", err)
			break
		}



		mu.Lock()
		mess := database.SendMessageParams{
			SenderID: incDat.Sender,
			Body: incDat.Body,
			ServerID: server.serv.ServerID,
		}
		db.SendMessage(server.context, mess)


		allMessages, err := db.RetrieveMessages(server.context, server.serv.ServerID)
		if err != nil {
			mu.Unlock()
			fmt.Println("server error retrieving messages from databse:", err)
			break
		}
		mu.Unlock()

		var chatData []Message //get sender nickname with inner join ******************************
		for _, log := range allMessages {
			chatData = append(chatData, Message{Body: log.Body, TimeSent: log.SentAt, Sender: log.SenderID})
		}

		jsonData, err := json.Marshal(chatData)
		if err != nil {
			fmt.Println("server error marshalling data:", err)
		}

		_, err = c.Write(jsonData)
		if err != nil {
			fmt.Println("server error writing data:", err)
		}

	}
	c.Close()
}

//func RegisterClient()

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        fmt.Println("error finding IP:", err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}