package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-apipoc-client/logging"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestHeadWhenAPIServerIsAvailable(t *testing.T) {
	os.Setenv("API_SERVER_ROOT", "http://bah.com")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("HEAD", "http://bah.com/ops/ping",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	client := NewHttpClient()

	response := client.Head("/ops/ping")

	assert.Equal(t, nil, response.Failure)
	assert.Equal(t, 200, response.Success.StatusCode)

	os.Unsetenv("API_SERVER_ROOT")
}

func TestHeadWhenAPIServerIsNotAvailable(t *testing.T) {
	os.Setenv("API_SERVER_ROOT", "http://bah.com")

	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	expectedFailure := hystrix.CircuitError{Message: "timeout"}

	client := NewHttpClient()

	response := client.Head("/ops/ping")

	assert.Equal(t, expectedFailure, response.Failure)

	os.Unsetenv("API_SERVER_ROOT")
}

func TestGetWhenAPIServerIsAvailable(t *testing.T) {
	os.Setenv("API_SERVER_ROOT", "http://bah.com")

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

	httpmock.RegisterResponder(
		"GET",
		"http://bah.com/ops/test",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

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

	os.Unsetenv("API_SERVER_ROOT")
}

func TestGetWithQueryParamsWhenAPIServerIsAvailable(t *testing.T) {
	os.Setenv("API_SERVER_ROOT", "http://bah.com")

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

	httpmock.RegisterResponder(
		"GET",
		"http://bah.com/ops/test?a=1&b=24",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

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

	os.Unsetenv("API_SERVER_ROOT")
}

func TestGetWhenAPIServerIsNotAvailable(t *testing.T) {
	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	expectedFailure := hystrix.CircuitError{Message: "timeout"}

	client := NewHttpClient()

	response := client.Get("/ops/status", nil)

	assert.Equal(t, expectedFailure, response.Failure)

	os.Unsetenv("API_SERVER_ROOT")
}

func TestPingWhenEnvironmentVariableIsNotSet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"HEAD",
		"https://api.develop.onsdigital.co.uk/ops/ping",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	client := NewHttpClient()

	response := client.Head("/ops/ping")

	assert.Equal(t, nil, response.Failure)
	assert.Equal(t, 200, response.Success.StatusCode)
}
