package logger

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/teyhouse/gowatch/filehandler"
)

// http.Get (http.DefaultClient) has zero Timeout, i.e. "wait forever" on an
// unresponsive upstream. Always use an explicit client with a timeout.
var httpClient = &http.Client{Timeout: 10 * time.Second}

// sendGET attaches the log message (base64-encoded) as the "logmessage"
// query parameter and issues the request. Query parameters are set via
// net/url so separators and escaping are always correct.
func sendGET(rawURL, message string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		if filehandler.CheckDebug() {
			fmt.Printf("invalid event URI %s: %s\n", rawURL, err)
		}
		return
	}

	q := u.Query()
	q.Set("logmessage", base64.StdEncoding.EncodeToString([]byte(message)))
	u.RawQuery = q.Encode()

	res, err := httpClient.Get(u.String())
	if err != nil {
		if filehandler.CheckDebug() {
			fmt.Printf("error making http request: %s\n", err)
		}
		return
	}
	defer res.Body.Close()

	if filehandler.CheckDebug() {
		fmt.Printf("HTTP response: %d\n URL: %s\n", res.StatusCode, u.String())
	}
}

func LogHTTP(message string) {
	if filehandler.FileExists("event.json") {
		uri, err := filehandler.GetEventURI()
		if err != nil {
			if filehandler.CheckDebug() {
				fmt.Printf("error reading event URI: %s\n", err)
			}
			return
		}
		sendGET(uri, message)
	}
}
