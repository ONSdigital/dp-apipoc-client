package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ONSdigital/dp-apipoc-client/logging"
	"github.com/ONSdigital/dp-apipoc-client/model"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

//Shim for 2 param return values
func M(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}

//Shim for 3 param return values
func M3(a, b, c interface{}) []interface{} {
	return []interface{}{a, b, c}
}

func TestPingWhenAPIServerIsAvailable(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"HEAD",
		"https://api.develop.onsdigital.co.uk/ops/ping",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	client := NewApiClient()

	assert.Equal(t, M(200, nil), M(client.Ping()))
}

func TestPingWhenAPIServerIsNotAvailable(t *testing.T) {
	config.Host = "https://api.develop.onsdigital.co.uk"

	client := NewApiClient()

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	assert.Equal(t, M(0, hystrix.CircuitError{Message: "timeout"}), M(client.Ping()))
}

func TestStatus(t *testing.T) {
	statusJson := []byte(`{
		  "applicationName": "API POC Server",
		  "dependencies": {
		    "elasticsearch": {
		      "status": "RUNNING",
		      "statusCode": 200,
		      "pingResponse": {
		        "name": "golden_child",
		        "cluster_name": "elasticsearch",
		        "version": {
		          "number": "2.4.3",
		          "build_hash": "d38a34e7b75af4e17ead16f156feffa432b22be3",
		          "build_timestamp": "2016-12-07T16:28:56Z",
		          "build_snapshot": false,
		          "lucene_version": "5.5.2"
		        },
		        "tagline": "You Know, for Search"
		      }
		    },
		    "website":{"status":"RUNNING","statusCode":200}
		  }
		}`)

	var respBody interface{}
	e := json.Unmarshal(statusJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/ops/status",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	expectedJson := model.Status{
		ApplicationName: "API POC Server",
		Dependencies: &model.Dependency{
			Elasticsearch: &model.Elastic{
				Status: "RUNNING",
				Code:   200,
				PingResponse: &model.PingResponse{
					Name:        "golden_child",
					ClusterName: "elasticsearch",
					Version: &model.ElasticVersion{
						Number:         "2.4.3",
						BuildHash:      "d38a34e7b75af4e17ead16f156feffa432b22be3",
						BuildTimestamp: "2016-12-07T16:28:56Z",
						BuildSnapshot:  false,
						LuceneVersion:  "5.5.2",
					},
					Tagline: "You Know, for Search",
				},
			},
			Website: &model.Website{
				Status: "RUNNING",
				Code:   200,
			},
		},
	}

	client := NewApiClient()

	assert.Equal(t, M3(200, expectedJson, nil), M3(client.Status()))
}

func TestGetDatasets(t *testing.T) {
	datasetJson := []byte(`{
		  "startIndex": 0,
		  "itemsPerPage": 1,
		  "totalItems": 35,
		  "items": [
		    {
		      "uri": "/businessindustryandtrade/business/businessinnovation/datasets/scienceandtechnologyclassification/current",
		      "type": "dataset",
		      "description": {
		        "title": "Science and Technology Classification",
		        "summary": "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
		        "keywords": [],
		        "metaDescription": "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
		        "nationalStatistic": false,
		        "contact": {
		          "email": "james.p.harris@ons.gsi.gov.uk",
		          "name": "James P Harris",
		          "telephone": "+44 (0) 1329 444 656"
		        },
		        "releaseDate": "2015-02-13T00:00:00Z",
		        "nextRelease": "",
		        "edition": "Current",
		        "datasetId": "",
		        "unit": "",
		        "preUnit": "",
		        "source": ""
		      },
		      "searchBoost": []
		    }
		  ]
		}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/dataset?limit=1&start=0",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	dt, er := time.Parse(time.RFC3339, "2015-02-13T00:00:00Z")
	if er != nil {
		panic(er)
	}

	expectedJson := model.Metadata{
		StartIndex:   0,
		ItemsPerPage: 1,
		TotalItems:   35,
		Items: &[]model.Record{
			{
				RecordUri:  "/businessindustryandtrade/business/businessinnovation/datasets/scienceandtechnologyclassification/current",
				RecordType: "dataset",
				Description: &model.Description{
					Title:             "Science and Technology Classification",
					Summary:           "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
					Keywords:          []string{},
					MetaDescription:   "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
					NationalStatistic: false,
					Contact: &model.Contact{
						Email:     "james.p.harris@ons.gsi.gov.uk",
						Name:      "James P Harris",
						Telephone: "+44 (0) 1329 444 656",
					},
					ReleaseDate: dt,
					NextRelease: "",
					Edition:     "Current",
					DatasetId:   "",
					DataUnit:    "",
					PreUnit:     "",
					Source:      "",
				},
				SearchBoost: []string{},
			},
		},
	}

	client := NewApiClient()

	assert.Equal(t, M3(200, expectedJson, nil), M3(client.GetDatasets(0, 1)))
}

func TestGetDatasetsForId(t *testing.T) {
	datasetJson := []byte(`{
		  "startIndex": 0,
		  "itemsPerPage": 1,
		  "totalItems": 35,
		  "items": [
		    {
		      "uri": "/businessindustryandtrade/business/businessinnovation/datasets/scienceandtechnologyclassification/current",
		      "type": "dataset",
		      "description": {
		        "title": "Science and Technology Classification",
		        "summary": "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
		        "keywords": [],
		        "metaDescription": "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
		        "nationalStatistic": false,
		        "contact": {
		          "email": "james.p.harris@ons.gsi.gov.uk",
		          "name": "James P Harris",
		          "telephone": "+44 (0) 1329 444 656"
		        },
		        "releaseDate": "2015-02-13T00:00:00Z",
		        "nextRelease": "",
		        "edition": "Current",
		        "datasetId": "",
		        "unit": "",
		        "preUnit": "",
		        "source": ""
		      },
		      "searchBoost": []
		    }
		  ]
		}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/dataset/ukea?limit=1&start=0",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	dt, er := time.Parse(time.RFC3339, "2015-02-13T00:00:00Z")
	if er != nil {
		panic(er)
	}

	expectedJson := model.Metadata{
		StartIndex:   0,
		ItemsPerPage: 1,
		TotalItems:   35,
		Items: &[]model.Record{
			{
				RecordUri:  "/businessindustryandtrade/business/businessinnovation/datasets/scienceandtechnologyclassification/current",
				RecordType: "dataset",
				Description: &model.Description{
					Title:             "Science and Technology Classification",
					Summary:           "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
					Keywords:          []string{},
					MetaDescription:   "The full list of UK Standard Industrial Classification of Economic Activities 2007 (SIC07) codes and comparison between different classifications, with identifiers to categorise all 5-digit SIC07 codes according to the Science and Technology classification.",
					NationalStatistic: false,
					Contact: &model.Contact{
						Email:     "james.p.harris@ons.gsi.gov.uk",
						Name:      "James P Harris",
						Telephone: "+44 (0) 1329 444 656",
					},
					ReleaseDate: dt,
					NextRelease: "",
					Edition:     "Current",
					DatasetId:   "",
					DataUnit:    "",
					PreUnit:     "",
					Source:      "",
				},
				SearchBoost: []string{},
			},
		},
	}

	client := NewApiClient()

	assert.Equal(t, M3(200, expectedJson, nil), M3(client.GetDatasetsForId("ukea", 0, 1)))
}

func TestGetDatasetsForTimeseries(t *testing.T) {
	datasetJson := []byte(`{}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/timeseries/nmcu/dataset?limit=1&start=0",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	client := NewApiClient()

	assert.Equal(t, M3(200, model.Metadata{}, nil), M3(client.GetDatasetsForTimeseries("nmcu", 0, 1)))
}

func TestGetTimeseries(t *testing.T) {
	datasetJson := []byte(`{}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/timeseries?limit=1&start=0",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	client := NewApiClient()

	assert.Equal(t, M3(200, model.Metadata{}, nil), M3(client.GetTimeseries(0, 1)))
}

func TestGetTimeseriesForId(t *testing.T) {
	datasetJson := []byte(`{}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/timeseries/nmcu?limit=1&start=0",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	client := NewApiClient()

	assert.Equal(t, M3(200, model.Metadata{}, nil), M3(client.GetTimeseriesForId("nmcu", 0, 1)))
}

func TestGetTimeseriesForDataset(t *testing.T) {
	datasetJson := []byte(`{}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/dataset/UKEA/timeseries?limit=1&start=0",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	client := NewApiClient()

	assert.Equal(t, M3(200, model.Metadata{}, nil), M3(client.GetTimeseriesForDataset("UKEA", 0, 1)))
}

func TestGetDataset(t *testing.T) {
	datasetJson := []byte(`{}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/dataset/UKEA/timeseries/nmcu",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	client := NewApiClient()

	assert.Equal(t, M3(200, model.Record{}, nil), M3(client.GetDataset("UKEA", "nmcu")))
}

func TestSearchMetadata(t *testing.T) {
	datasetJson := []byte(`{}`)

	var respBody interface{}
	e := json.Unmarshal(datasetJson, &respBody)

	if e != nil {
		panic(e)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder(
		"GET",
		"https://api.develop.onsdigital.co.uk/search?limit=7&q=travel&start=3",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, respBody)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	logging.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	client := NewApiClient()

	assert.Equal(t, M3(200, model.Metadata{}, nil), M3(client.Search("travel", 3, 7)))
}
