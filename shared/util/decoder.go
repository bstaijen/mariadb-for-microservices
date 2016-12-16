package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RequestToJSON(req *http.Request, target interface{}) error {
	defer req.Body.Close()
	return toJSON(req.Body, target)
}

func toJSON(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		fmt.Printf("json decoder error occured: %v \n", err.Error())
		return err
	}
	return nil
}
