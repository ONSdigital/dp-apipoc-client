package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestHeadWhenAPIServerIsAvailable(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("HEAD", "https://api.develop.onsdigital.co.uk/ops/ping",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	config.Host = "https://api.develop.onsdigital.co.uk"

	client := NewHttpClient()

	response := client.Head("/ops/ping")

	assert.Equal(t, nil, response.Failure)
	assert.Equal(t, 200, response.Success.StatusCode)
}

func TestHeadWhenAPIServerIsNotAvailable(t *testing.T) {
	config.Host = "https://api.develop.onsdigital.co.uk"

	expectedFailure := hystrix.CircuitError{Message: "timeout"}

	client := NewHttpClient()

	response := client.Head("/ops/ping")

	assert.Equal(t, expectedFailure, response.Failure)
}

func TestGetWhenAPIServerIsAvailable(t *testing.T) {
	byteJson := []byte(`{
		  "applicationName": "API POC Server",
		  "dependencies": {
			"dependency1":{},
			"dependency2":{}
		  }
		}`)

	var respBody interface{}
	e := json.Unmarshal(byteJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/ops/test",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	client := NewHttpClient()

	response := client.Get("/ops/test", nil)

	defer response.Success.Body.Close()

	bodyBytes, e2 := ioutil.ReadAll(response.Success.Body)
	if e2 != nil {
		panic(e2)
	}

	var bodyJson interface{}
	err := json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, nil, response.Failure)
	assert.Equal(t, 200, response.Success.StatusCode)
	assert.Equal(t, respBody, bodyJson)
}

func TestGetWithQueryParamsWhenAPIServerIsAvailable(t *testing.T) {
	byteJson := []byte(`{
		  "applicationName": "API POC Server",
		  "dependencies": {
			"dependency1":{},
			"dependency2":{}
		  }
		}`)

	var respBody interface{}
	e := json.Unmarshal(byteJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/ops/test?a=1&b=24",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	client := NewHttpClient()

	params := make(map[string]string)
	params["a"] = strconv.Itoa(1)
	params["b"] = strconv.Itoa(24)

	response := client.Get("/ops/test", params)

	defer response.Success.Body.Close()

	bodyBytes, e2 := ioutil.ReadAll(response.Success.Body)
	if e2 != nil {
		panic(e2)
	}

	var bodyJson interface{}
	err := json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, nil, response.Failure)
	assert.Equal(t, 200, response.Success.StatusCode)
	assert.Equal(t, respBody, bodyJson)
}

func TestGetWhenAPIServerIsNotAvailable(t *testing.T) {
	config.Host = "https://api.develop.onsdigital.co.uk"

	expectedFailure := hystrix.CircuitError{Message: "timeout"}

	client := NewHttpClient()

	response := client.Get("/ops/status", nil)

	assert.Equal(t, expectedFailure, response.Failure)
}
