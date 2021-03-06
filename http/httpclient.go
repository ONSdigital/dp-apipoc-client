package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ONSdigital/dp-apipoc-client/logging"
	"github.com/ONSdigital/dp-apipoc-client/model"
	"github.com/afex/hystrix-go/hystrix"
)

type HttpClient interface {
	Head(path string) model.Response
	Get(path string, params map[string]string) model.Response
}

func NewHttpClient() HttpClient {
	var serverRoot string
	serverRoot = os.Getenv("API_SERVER_ROOT")
	if len(serverRoot) <= 0 {
		serverRoot = "https://api.develop.onsdigital.co.uk"
	}

	hystrix.ConfigureCommand("default_config", hystrix.CommandConfig{
		Timeout: 20,
	})

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return &httpService{
		apiServerUrl: serverRoot,
		httpClient:   netClient,
	}
}

type httpService struct {
	apiServerUrl string
	httpClient   *http.Client
}

func (s *httpService) Head(path string) model.Response {
	url := fmt.Sprintf("%s%s", s.apiServerUrl, path)

	responseChannel := make(chan model.Response)

	hystrix.Go("default_config", func() error {
		resp, err := s.httpClient.Head(url)

		responseChannel <- model.Response{Success: resp, Failure: err}
		return nil
	}, func(err error) error {
		responseChannel <- model.Response{Success: &http.Response{}, Failure: err}
		return nil
	})

	return <-responseChannel
}

func (s *httpService) Get(path string, params map[string]string) model.Response {
	url := fmt.Sprintf("%s%s", s.apiServerUrl, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logging.Error.Println(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}

	req.Header.Add("Accept", "application/json")
	req.URL.RawQuery = q.Encode()

	logging.Info.Println(req)

	responseChannel := make(chan model.Response)

	hystrix.Go("default_config", func() error {
		resp, err := s.httpClient.Do(req)

		responseChannel <- model.Response{Success: resp, Failure: err}
		return nil
	}, func(err error) error {
		responseChannel <- model.Response{Success: &http.Response{}, Failure: err}
		return nil
	})

	return <-responseChannel
}
