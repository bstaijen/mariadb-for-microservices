package util

import (
	"bytes"
	"log"
	"net/http"
)

// Request is a helper which executes a request over the network and returns an error or a response
func Request(method, url string, body []byte, cb func(*http.Response)) error {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request: " + err.Error())
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error executing request: " + err.Error())
		return err
	}

	// callback
	log.Printf("[info] %v %v %v", method, url, resp.StatusCode)
	cb(resp)

	//return
	return nil
}
