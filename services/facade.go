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
)

func main() {
	panic(http.ListenAndServe(basement.FacadeAddress, &FacadeResponse{}))
}

type FacadeResponse struct{}

func (m *FacadeResponse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responseBody string
	if request.Method == "GET" {
		responseBody = listLogs() + "\n" + listMessages()
	} else if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
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
		err = sendToLogger(messageToLogging)
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

func listMessages() string {
	messages, _ := basement.GetData(basement.MessagesServiceAddress)
	return messages
}

func sendToLogger(message []byte) error {
	for tryCount := 0; tryCount < 15; tryCount++ {
		i := rand.Int() % 3
		_, err := http.Post(
			basement.Localhost+basement.LoggingServiceAddress[i],
			"application/json",
			bytes.NewReader(message))
		if err == nil {
			return nil
		} else {
			tryCount++
			fmt.Println("Cannot send message to logger " + basement.LoggingServiceAddress[i])
		}
	}
	return errors.New(basement.MsgServicesNotResponding)
}