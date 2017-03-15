package core

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestSubtract(t *testing.T) {
	//	goatJson := []byte(`{
	//                "result": [{
	//			"id": 11111,
	//			"lat": "38.898648N",
	//			"lon": "77.037692W"
	//		},
	//		{
	//			"id": 44444,
	//			"lat": "6.55555555N",
	//			"lon": "89.9999999W"
	//		}
	//	]}`)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our database of articles
	//	articles := make([]map[string]interface{}, 0)

	config.Host = "https://api.develop.onsdigital.co.uk"

	httpmock.RegisterResponder("HEAD", "https://api.develop.onsdigital.co.uk/ops/ping",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	response := PingResponse{StatusCode: 200, Error: nil}

	client := NewHttpClient()

	assert.Equal(t, response, client.Ping())
}
