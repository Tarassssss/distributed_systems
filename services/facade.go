package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"basement"
)

func main() {
	panic(http.ListenAndServe(basement.FacadeAddress, &FacadeResponse{}))
}

type FacadeResponse struct{}

func (m *FacadeResponse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responseBody string
	if request.Method == "GET" {
		responseBody = basement.GetData(basement.LoggingServiceAddress) + "\n" + basement.GetData(basement.MessagesServiceAddress)
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
		info := basement.RequestInfo{uuid.New().String(), text.Msg}
		messageToLogging, _ := json.Marshal(info)
		_, err = http.Post(
			basement.Localhost+basement.LoggingServiceAddress,
			"application/json",
			bytes.NewReader(messageToLogging))
		if err != nil {
			responseBody = "Error occured!"
		} else {
			responseBody = "Message sent successfully!"
		}
	}
	_, _ = writer.Write([]byte(responseBody))
}