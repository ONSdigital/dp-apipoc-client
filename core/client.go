package core

import (
	"fmt"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	//	"github.com/ONSdigital/go-ns/log"
)

type HttpClient interface {
	Ping() PingResponse
	//	GetMetadata(url string) (model.Response, error)
	//	GetData(url string) (model.Response, error)
}

func NewHttpClient() HttpClient {
	hystrix.ConfigureCommand("default_config", hystrix.CommandConfig{
		Timeout: 20,
	})

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return &httpService{httpClient: netClient}
}

type httpService struct {
	httpClient *http.Client
}

func (s *httpService) Ping() PingResponse {
	url := fmt.Sprintf("%s/ops/ping", config.Host)

	responseChannel := make(chan PingResponse)

	hystrix.Go("default_config", func() error {
		resp, err := s.httpClient.Head(url)

		responseChannel <- PingResponse{StatusCode: resp.StatusCode, Error: err}
		return nil
	}, func(err error) error {
		responseChannel <- PingResponse{Error: err}
		return nil
	})

	return <-responseChannel
}

//func retrieveOpenIssues(userConfig UserDefinedConfig) JiraSearchResponse {
//	url := fmt.Sprintf("%s/rest/api/2/search", userConfig.Host)
//	query := fmt.Sprintf("project = %s AND resolution = Unresolved ORDER BY priority DESC", userConfig.Project)
//	fields := []string{"summary", "key", "status", "assignee", "description", "issuetype", "created"}

//	jsonStr, _ := json.Marshal(JiraSearchRequest{query, 0, fields})

//	responseChannel := make(chan JiraSearchResponse)

//	hystrix.ConfigureCommand("Get all Issues", hystrix.CommandConfig{Timeout: 2000})

//	hystrix.Go("Get all Issues", func() error {
//		body := httpRequest(url, "POST", jsonStr, userConfig)
//		var searchRes JiraSearchResponse
//		json.Unmarshal(body, &searchRes)
//		responseChannel <- searchRes
//		return nil
//	}, func(err error) error {
//		var searchRes JiraSearchResponse
//		responseChannel <- searchRes
//		return nil
//	})

//	return <-responseChannel
//}
