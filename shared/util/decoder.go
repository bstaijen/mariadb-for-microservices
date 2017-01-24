package util

import (
	"encoding/json"
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
		return err
	}
	return nil
}
