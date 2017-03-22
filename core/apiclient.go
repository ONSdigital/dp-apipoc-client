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

	GetDatasets(start int, limit int) (int, model.Metadata, error)
	GetDatasetsForId(datasetId string, start int, limit int) (int, model.Metadata, error)
	GetDatasetsForTimeseries(timeseriesId string, start int, limit int) (int, model.Metadata, error)
	GetTimeseries(start int, limit int) (int, model.Metadata, error)
	GetTimeseriesForId(timeseriesId string, start int, limit int) (int, model.Metadata, error)
	GetTimeseriesForDataset(datasetId string, start int, limit int) (int, model.Metadata, error)
	GetDataset(datasetId string, timeseriesId string) (int, model.Record, error)
	Search(term string, start int, limit int) (int, model.Metadata, error)
	getMetadata(path string, params map[string]string) (int, model.Metadata, error)

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

func (s *apiService) GetDatasets(start int, limit int) (int, model.Metadata, error) {
	params := make(map[string]string)
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	return s.getMetadata("/dataset", params)
}

func (s *apiService) GetDatasetsForId(datasetId string, start int, limit int) (int, model.Metadata, error) {
	path := buildPath([]string{"/dataset/", datasetId})
	params := make(map[string]string)
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	return s.getMetadata(path, params)
}

func (s *apiService) GetDatasetsForTimeseries(timeseriesId string, start int, limit int) (int, model.Metadata, error) {
	path := buildPath([]string{"/timeseries/", timeseriesId, "/dataset"})
	params := make(map[string]string)
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	return s.getMetadata(path, params)
}

func (s *apiService) GetTimeseries(start int, limit int) (int, model.Metadata, error) {
	params := make(map[string]string)
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	return s.getMetadata("/timeseries", params)
}

func (s *apiService) GetTimeseriesForId(timeseriesId string, start int, limit int) (int, model.Metadata, error) {
	path := buildPath([]string{"/timeseries/", timeseriesId})
	params := make(map[string]string)
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	return s.getMetadata(path, params)
}

func (s *apiService) GetTimeseriesForDataset(datasetId string, start int, limit int) (int, model.Metadata, error) {
	path := buildPath([]string{"/dataset/", datasetId, "/timeseries"})
	params := make(map[string]string)
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	return s.getMetadata(path, params)
}

func (s *apiService) GetDataset(datasetId string, timeseriesId string) (int, model.Record, error) {
	path := buildPath([]string{"/dataset/", datasetId, "/timeseries/", timeseriesId})

	resp := s.httpClient.Get(path, nil)

	if resp.Failure != nil {
		logging.Error.Println(resp.Failure)

		return resp.Success.StatusCode, model.Record{}, resp.Failure
	}

	defer resp.Success.Body.Close()

	bodyBytes, e := ioutil.ReadAll(resp.Success.Body)
	if e != nil {
		logging.Error.Println(e)
		panic(e)
	}

	var body model.Record
	err := json.Unmarshal(bodyBytes, &body)
	if err != nil {
		logging.Error.Println(err)
		panic(err)
	}

	return resp.Success.StatusCode, body, nil
}

func (s *apiService) Search(term string, start int, limit int) (int, model.Metadata, error) {
	params := make(map[string]string)
	params["q"] = term
	params["start"] = strconv.Itoa(start)
	params["limit"] = strconv.Itoa(limit)

	return s.getMetadata("/search", params)
}

func (s *apiService) getMetadata(path string, params map[string]string) (int, model.Metadata, error) {
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
