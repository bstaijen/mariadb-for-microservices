package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
)

// RequestToJSON converts a Request containing a json object to a targeted interface
func RequestToJSON(req *http.Request, target interface{}) error {
	if req.Body == nil {
		return errors.New("Bad json")
	}
	defer req.Body.Close()
	return toJSON(req.Body, target)
}

// ResponseJSONToObject converts a Response containing a json object to a targeted interface
func ResponseJSONToObject(res *http.Response, target interface{}) error {
	if res.Body == nil {
		return errors.New("Bad json")
	}
	defer res.Body.Close()
	return toJSON(res.Body, target)
}

// ResponseRecorderJSONToObject converts a ResponseRecorder containing a json to a targeted interface. This is only used in test cases when a httptest server is used.
func ResponseRecorderJSONToObject(res *httptest.ResponseRecorder, target interface{}) error {
	if res.Body == nil {
		return errors.New("Bad json")
	}
	return toJSON(res.Body, target)
}

func toJSON(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		return err
	}
	return nil
}
