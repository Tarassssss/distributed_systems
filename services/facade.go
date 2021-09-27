package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"basement"
	"errors"
	"fmt"
	"math/rand"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel
var messageQueue amqp.Queue

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ = conn.Channel()
	messageQueue, _ = ch.QueueDeclare("messaging_service", true, false, false, false, nil)
	panic(http.ListenAndServe(basement.FacadeAddress, &FacadeResponse{}))
}

type FacadeResponse struct{}

func (m *FacadeResponse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responseBody string
	if request.Method == "GET" {
		responseBody = listLogs() + "\n" + listMessages()
	} else if request.Method == "POST" {
		/*body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		var text basement.RequestData
		err = json.Unmarshal(body, &text)
		if err != nil {
			panic(err)
		}
		info := basement.RequestInfo{Id: uuid.New().String(), Msg: text.Msg}
		messageToLogging, _ := json.Marshal(info)
		err = sendToLogger(messageToLogging)*/
		reqParams := parseRequest(request)
		err := sendToLogger(reqParams)
		sendToMessenger(reqParams)
		if err != nil {
			responseBody = err.Error()
		} else {
			responseBody = "Message sent successfully!"
		}
	}
	_, _ = writer.Write([]byte(responseBody))
}

func listLogs() string {
	for i := 0; i < basement.LoggerPortsSize; i++ {
		logs, err := basement.GetData(basement.LoggingServiceAddress[i])
		if err != nil {
			fmt.Println("Logger on port " + basement.LoggingServiceAddress[i] + " is not responding")
		} else {
			return logs
		}
	}
	return basement.MsgServicesNotResponding
}

func parseRequest(request *http.Request) basement.RequestData {
	body, _ := ioutil.ReadAll(request.Body)
	var params basement.RequestData
	_ = json.Unmarshal(body, &params)
	return params
}

func listMessages() string {
	i := rand.Int() % 3
	messageServiceAddress := basement.MessagesServiceAddress[i]
	messages, _ := basement.GetData(messageServiceAddress)
	return messages
}

//func sendToLogger(message []byte) error {
func sendToLogger(reqParams basement.RequestData) error {
	info := basement.RequestInfo{Id: uuid.New().String(), Msg: reqParams.Msg}
	logRequestMessage, _ := json.Marshal(info)	
	for tryCount := 0; tryCount < 15; tryCount++ {
		i := rand.Int() % 3
		_, err := http.Post(
			basement.Localhost+basement.LoggingServiceAddress[i],
			"application/json",
			//bytes.NewReader(message))
			bytes.NewReader(logRequestMessage))
		if err == nil {
			return nil
		} else {
			tryCount++
			fmt.Println("Cannot send message to logger " + basement.LoggingServiceAddress[i])
		}
	}
	return errors.New(basement.MsgServicesNotResponding)
}

func sendToMessenger(reqParams basement.RequestData) {
	_ = ch.Publish(
		"", messageQueue.Name, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(reqParams.Msg),
		})
}