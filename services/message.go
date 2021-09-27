package main

import (
	"net/http"
	"basement"
	"fmt"
	"github.com/streadway/amqp"
	"net"
	"strings"
)

var ch *amqp.Channel
var messageQueue amqp.Queue
var messages []string

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ = conn.Channel()
	messageQueue, _ = ch.QueueDeclare("messaging_service", true, false, false, false, nil)
	msgs, _ := ch.Consume(messageQueue.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			msg := string(d.Body)
			messages = append(messages, msg)
			fmt.Println("Received a message:", msg)
		}
	}()

	panic(http.ListenAndServe(findAvailablePort(), &MessagingResponse{}))
	//panic(http.ListenAndServe(basement.MessagesServiceAddress, &MessagingResponse{}))
}

type MessagingResponse struct{}

func (m *MessagingResponse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		text := strings.Join(messages, "\n")
		_, _ = writer.Write([]byte(text))
	}

}

func findAvailablePort() string {
	for i := 0; i < len(basement.MessagesServiceAddress); i++ {
		listener, err := net.Listen("tcp", basement.MessagesServiceAddress[i])
		if err != nil {
			fmt.Println(basement.MessagesServiceAddress[i], " is already taken, try next one")
		} else {
			_ = listener.Close()
			return basement.MessagesServiceAddress[i]
		}
		//_, _ = writer.Write([]byte("Message-service is not implemented yet"))
	}
	panic("All ports are unavailable")
}