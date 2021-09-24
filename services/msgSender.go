   
package main

import (
	"bytes"
	"net/http"
	"basement"
	"strconv"
)

func main() {
	for i := 0; i < 10; i++ {
		message := `{"msg": "msg` + strconv.Itoa(i) + `"}`
		_, _ = http.Post(
			basement.Localhost+basement.FacadeAddress,
			"application/json",
			bytes.NewReader([]byte(message)))
	}
}