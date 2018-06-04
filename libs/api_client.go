package libs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/denperov/owm-task/libs/log"
	"github.com/gorilla/schema"
)

var queryEncoder = schema.NewEncoder()

type APIClient struct {
	Address url.URL
}

func NewAPIClient(address string) *APIClient {
	addressURL, err := url.Parse(address)
	if err != nil {
		log.Panic(err)
	}
	return &APIClient{
		Address: *addressURL,
	}
}

func (c *APIClient) Get(apiPath string, val interface{}, result interface{}) error {
	queryValues := make(url.Values)
	if err := queryEncoder.Encode(val, queryValues); err != nil {
		log.Panic(err)
	}

	u := c.Address
	u.Path = path.Join(u.Path, apiPath)
	u.RawQuery = queryValues.Encode()

	log.Infof("GET: %s", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	return c.readResponse(resp, result)
}

func (c *APIClient) Post(apiPath string, val interface{}, result interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		log.Panic(err)
	}

	u := c.Address
	u.Path = path.Join(u.Path, apiPath)

	log.Infof("POST: %s, %s", u.String(), string(data))

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	return c.readResponse(resp, result)
}

func (c *APIClient) readResponse(resp *http.Response, result interface{}) error {
	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			log.Panic(err)
		}
	case http.StatusBadRequest:
		var errBody struct {
			Message string
		}
		err := json.NewDecoder(resp.Body).Decode(&errBody)
		if err != nil {
			log.Panic(err)
		}
		return errors.New(errBody.Message)
	default:
		errBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicf("unexpected status code: %d", resp.Status)
		}
		log.Panicf("unexpected status code %d, body: %s", resp.Status, string(errBody))
	}
	return nil
}
