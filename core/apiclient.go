package core

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/ONSdigital/dp-apipoc-client/logging"
	"github.com/ONSdigital/dp-apipoc-client/model"
)

type ApiClient interface {
	Ping() (int, error)
	Status() (int, model.Status, error)
	GetDataset(id string, start int, limit int) (int, model.Metadata, error)
	getMetadata(path string, start int, limit int) (int, model.Metadata, error)
	//	GetData(url string) (model.Response, error)
}

func NewApiClient() ApiClient {
	return &apiService{httpClient: NewHttpClient()}
}

type apiService struct {
	httpClient HttpClient
}

func (s *apiService) Ping() (int, error) {
	resp := s.httpClient.Head("/ops/ping")

	if resp.Failure != nil {
		logging.Error.Println(resp.Failure)

		return 0, resp.Failure
	}

	return resp.Success.StatusCode, nil
}

func (s *apiService) Status() (int, model.Status, error) {
	resp := s.httpClient.Get("/ops/status", nil)

	if resp.Failure != nil {
		logging.Error.Println(resp.Failure)

		return resp.Success.StatusCode, model.Status{}, resp.Failure
	}

	defer resp.Success.Body.Close()

	bodyBytes, e := ioutil.ReadAll(resp.Success.Body)
	if e != nil {
		logging.Error.Println(e)
		panic(e)
	}

	var body model.Status
	err := json.Unmarshal(bodyBytes, &body)
	if err != nil {
		logging.Error.Println(err)
		panic(err)
	}

	return resp.Success.StatusCode, body, nil
}

func (s *apiService) GetDataset(id string, start int, limit int) (int, model.Metadata, error) {
	return s.getMetadata("/dataset/"+id, start, limit)
}

func (s *apiService) getMetadata(path string, start int, limit int) (int, model.Metadata, error) {
	params := make(map[string]string)
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	resp := s.httpClient.Get(path, params)

	if resp.Failure != nil {
		logging.Error.Println(resp.Failure)

		return resp.Success.StatusCode, model.Metadata{}, resp.Failure
	}

	defer resp.Success.Body.Close()

	bodyBytes, e := ioutil.ReadAll(resp.Success.Body)
	if e != nil {
		logging.Error.Println(e)
		panic(e)
	}

	var body model.Metadata
	err := json.Unmarshal(bodyBytes, &body)
	if err != nil {
		logging.Error.Println(err)
		panic(err)
	}

	return resp.Success.StatusCode, body, nil
}