package logger

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/teyhouse/gowatch/filehandler"
)

func sendGET(url, params string) {
	url += fmt.Sprintf("&logmessage=%s", params)
	res, err := http.Get(url)
	if err != nil {
		if filehandler.CheckDebug() {
			fmt.Printf("error making http request: %s\n", err)
		}
	}

	if filehandler.CheckDebug() {
		fmt.Printf("HTTP response: %d\n URL: %s\n", res.StatusCode, url)

	}
}

func LogHTTP(message string) {
	if filehandler.FileExists("event.json") {
		//Attach Log-Message base64-encoded as last param to GET-Request
		sendGET(filehandler.GetEventURI(), base64.StdEncoding.EncodeToString([]byte(message)))
	}
}
