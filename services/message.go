package main

import (
	"net/http"
	"basement"
)

func main() {
	panic(http.ListenAndServe(basement.MessagesServiceAddress, &MessagingResponse{}))
}

type MessagingResponse struct{}

func (m *MessagingResponse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		_, _ = writer.Write([]byte("Message-service is not implemented yet"))
	}
}