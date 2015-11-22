package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
	"utils"
)

type Client struct {
	id          uint64
	name        string
	conn        net.Conn
	connectedAt int64

	reader *bufio.Reader
	writer *bufio.Writer

	incoming chan Message
	outgoing chan Message
	quiting  chan *Client
}

func InitClient(conn net.Conn, id uint64) (client *Client) {
	client = &Client{
		id:          id,
		name:        fmt.Sprintf("%v的%v", utils.RandomPrefix(), utils.RandomName()),
		conn:        conn,
		connectedAt: time.Now().UTC().Unix(),
		incoming:    make(chan Message),
		outgoing:    make(chan Message),
		quiting:     make(chan *Client),
		reader:      bufio.NewReader(conn),
		writer:      bufio.NewWriter(conn),
	}

	client.listen()
	return
}

func (client *Client) listen() {
	go client.read()
	go client.write()
}

func (client *Client) read() {
	for {
		if line, _, err := client.reader.ReadLine(); err == nil {
			input := strings.TrimSpace(string(line))
			if input != "" {
				msg := Message{
					client:      client,
					messageType: MESSAGE_TYPE_PLAYER,
					message:     string(line),
					sentAt:      time.Now().UTC().Unix(),
				}

				client.incoming <- msg
			}

		} else {
			utils.Log().Error("接收「%v」的消息失败: %v", client.name, err)
			client.quit()
			return

		}
	}
}

func (client *Client) write() {
	for data := range client.outgoing {
		if _, err := client.writer.WriteString(data.Stringer()); err != nil {
			utils.Log().Error("为「%v」发送消息失败1: %v", client.name, err)
			client.quit()
			return
		}

		if err := client.writer.Flush(); err != nil {
			utils.Log().Error("为「%v」发送消息失败2: %v", client.name, err)
			client.quit()
			return
		}
	}

}

func (client *Client) quit() {
	client.quiting <- client

}

func (client *Client) close() {
	client.conn.Close()
}
