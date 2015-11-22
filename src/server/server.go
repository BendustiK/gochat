package server

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"utils"
)

type MesasgeType int

const (
	MESSAGE_TYPE_SYSTEM MesasgeType = 0
	MESSAGE_TYPE_PLAYER MesasgeType = 1
)

type Message struct {
	client      *Client
	messageType MesasgeType
	message     string
	sentAt      int64
}

func (msg *Message) Stringer() string {
	t := time.Unix(msg.sentAt, 0)
	nt := t.Format("2006-01-02 15:04:05")

	switch msg.messageType {
	case MESSAGE_TYPE_PLAYER:
		return fmt.Sprintf("(%v)「%v」说「%v」\r\n", nt, msg.client.name, strings.TrimSpace(msg.message))
	case MESSAGE_TYPE_SYSTEM:
		return fmt.Sprintf("*********** (%v) - %v\r\n", msg.sentAt, msg.message)
	}

	return msg.message
}

type Server struct {
	startedAt       int64
	currentClientId uint64
	clients         map[uint64]*Client
	listener        net.Listener

	connecting chan net.Conn
	incoming   chan Message
	quiting    chan *Client
}

func SystemMessage(msg string) (message Message) {
	message = Message{
		messageType: MESSAGE_TYPE_SYSTEM,
		message:     msg,
		sentAt:      time.Now().UTC().Unix(),
	}
	return
}

func InitServer() (server *Server) {
	server = &Server{
		startedAt:  time.Now().UTC().Unix(),
		clients:    make(map[uint64]*Client, 100),
		connecting: make(chan net.Conn),
		incoming:   make(chan Message),
		quiting:    make(chan *Client),
	}

	server.listen()

	return
}

func (server *Server) listen() {
	go func() {
		for {
			select {
			case conn := <-server.connecting:
				server.join(conn)
			case msg := <-server.incoming:
				server.broadcast(msg)
			case client := <-server.quiting:
				server.quit(client)
			default:
				// utils.Log().Debug("Idle...")
			}
		}
	}()
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		utils.Log().Error("服务器启动失败: %v", err)
		os.Exit(-1)
	}
	server.listener = listener

	utils.Log().Notice("服务器已经启动, 正在监听8080端口")
	for {
		if conn, err := listener.Accept(); err != nil {
			utils.Log().Error("接收信息失败: %v", err)
		} else {
			server.connecting <- conn
		}

	}
}

func (server *Server) join(conn net.Conn) {
	server.currentClientId = server.currentClientId + 1
	client := InitClient(conn, server.currentClientId)
	server.clients[client.id] = client

	msg := SystemMessage(fmt.Sprintf("「%v」加入聊天:「%v」", client.name, conn.RemoteAddr().String()))
	server.broadcast(msg)
	utils.Log().Debug(msg.Stringer())

	// Message
	go func() {
		for {
			msg := <-client.incoming

			server.incoming <- msg
		}
	}()

	// Quit
	go func() {
		for {
			c := <-client.quiting
			server.quiting <- c
		}
	}()
}

func (server *Server) quit(client *Client) {
	if _, ok := server.clients[client.id]; ok && client.conn != nil {
		msg := SystemMessage(fmt.Sprintf("「%v」退出聊天:「%v」", client.name, client.conn.RemoteAddr().String()))
		client.close()
		delete(server.clients, client.id)

		utils.Log().Debug(msg.Stringer())
		server.broadcast(msg)

	}

}

func (server *Server) broadcast(msg Message) {
	utils.Log().Debug(msg.Stringer())
	for _, client := range server.clients {
		client.outgoing <- msg
	}
}
