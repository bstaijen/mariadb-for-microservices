package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"io/ioutil"

	simplejson "github.com/bitly/go-simplejson"
)

func RequestToJSON(req *http.Request, target interface{}) error {
	defer req.Body.Close()
	return toJSON(req.Body, target)
}

func ResponseToJSON(res *http.Response, target interface{}) error {
	defer res.Body.Close()
	return toJSON(res.Body, target)
}

func toJSON(r io.Reader, target interface{}) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		fmt.Printf("json decoder error occured: %v \n", err.Error())
		return err
	}
	return nil
}

func ToSimpleJson(req *http.Request) (*simplejson.Json, error) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	js, err := simplejson.NewJson(body)
	/*if err != nil {
		fmt.Println(err.Error())
	}*/
	return js, err
}

func ResponseToSimpleJSON(res *http.Response) (*simplejson.Json, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	js, err := simplejson.NewJson(body)
	return js, err
}

/*func JSONDecoder(req *http.Request, target interface{}) error {
	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(target)
	if err != nil {
		fmt.Printf("json decoder error occured: %v \n", err.Error())
		return err
	}
	return nil
}*/
