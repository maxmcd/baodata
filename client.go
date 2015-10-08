package baodata

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	secret   string
	baoId    string
	mock     bool
	mockData map[string][]Data
	endpoint string
)

func init() {
	endpoint = "https://gitbao.com/ds/"

	secret = os.Getenv("_B_secret")
	baoId = os.Getenv("_B_baoId")

	if secret == "" || baoId == "" {
		mock = true
		mockData = make(map[string][]Data)
	}
}

func Get(location string) (data []Data, err error) {
	br, err := request("get", location, Data{})
	if err != nil {
		return
	}
	data = br
	return
}

func Put(location string, putData Data) (data Data, err error) {
	br, err := request("put", location, putData)
	if len(br) == 1 {
		data = br[0]
	}
	return
}

func Delete(location string) (err error) {
	_, err = request("delete", location, Data{})
	return
}

func parseLocation(loc string) (data PathData, err error) {
	// add support for query strings, sorting, relations, etc....

	var locParts []string
	lp := strings.Split(loc, "/")
	for _, val := range lp {
		if val != "" {
			locParts = append(locParts, val)
		}
	}
	if len(locParts) > 2 || len(locParts) < 1 {
		err = errors.New(
			`path does not match "/resource_name/id" or "/resource_name" format`,
		)
		return
	}

	var id int
	if len(locParts) == 2 {
		idString := path.Base(locParts[1])
		id, err = strconv.Atoi(idString)
		if err != nil {
			return
		}
	}

	resourceName := locParts[0]

	data.ResourceName = resourceName
	data.Id = id
	return
}

func request(method, loc string, data Data) (
	br []Data, err error) {
	br, err = requestData(method, loc, data)

	if (err == nil) &&
		(method == "get" || method == "put") &&
		(len(br) == 0) {
		//
		err = errors.New("not found")
	}
	return
}

func requestData(method, loc string, data Data) (
	br []Data, err error) {

	pd, err := parseLocation(loc)
	if err != nil {
		return
	}

	if mock == false {
		var b bytes.Buffer

		breq := BaoRequest{
			Method: method,
			Secret: secret,
			BaoId:  baoId,
			Path:   pd,
			Data:   Data(data),
		}

		jsonBytes, err := json.Marshal(breq)
		if err != nil {
			return br, err
		}

		b.Write(jsonBytes)

		resp, err := http.Post(endpoint, "text/json", &b)
		if err != nil {
			return br, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return br, err
		}
		var baoResponse BaoResponse
		err = json.Unmarshal(body, &baoResponse)
		if baoResponse.Error != "" {
			err = errors.New(baoResponse.Error)
			return br, err
		}
		for _, d := range baoResponse.Data {
			br = append(br, Data(d))
		}
		return br, err

	} else {
		// mock mode
		if method == "get" {
			br = mockGet(pd)
			return
		} else if method == "put" {
			br = mockPut(pd, data)
			return
		} else if method == "delete" {
			mockDelete(pd)
			return
		}
	}
	return
}
