package model

import (
	"net/http"
	"time"
)

type Response struct {
	Success *http.Response
	Failure error
}

type Status struct {
	ApplicationName string      `json:"applicationName"`
	Dependencies    *Dependency `json:"dependencies"`
}

type Dependency struct {
	Elasticsearch *Elastic `json:"elasticsearch"`
	Website       *Website `json:"website"`
}

type Elastic struct {
	Status       string        `json:"status"`
	Code         int           `json:"statusCode"`
	PingResponse *PingResponse `json:"pingResponse"`
}

type PingResponse struct {
	Name        string          `json:"name"`
	ClusterName string          `json:"cluster_name"`
	Version     *ElasticVersion `json:"version"`
	Tagline     string          `json:"tagline"`
}

type ElasticVersion struct {
	Number         string `json:"number"`
	BuildHash      string `json:"build_hash"`
	BuildTimestamp string `json:"build_timestamp"`
	BuildSnapshot  bool   `json:"build_snapshot"`
	LuceneVersion  string `json:"lucene_version"`
}

type Website struct {
	Status string `json:"status"`
	Code   int    `json:"statusCode"`
}

type Metadata struct {
	StartIndex   int       `json:"startIndex"`
	ItemsPerPage int       `json:"itemsPerPage"`
	TotalItems   int       `json:"totalItems"`
	Items        *[]Record `json:"items"`
}

type Record struct {
	RecordUri   string       `json:"uri"`
	RecordType  string       `json:"type"`
	Description *Description `json:"description"`
	SearchBoost []string     `json:"searchBoost"`
}

type Description struct {
	Title             string    `json:"title"`
	Summary           string    `json:"summary"`
	Keywords          []string  `json:"keywords"`
	MetaDescription   string    `json:"metaDescription"`
	NationalStatistic bool      `json:"nationalStatistic"`
	Contact           *Contact  `json:"contact"`
	ReleaseDate       time.Time `json:"releaseDate,string"`
	NextRelease       string    `json:"nextRelease"`
	Edition           string    `json:"edition"`
	DatasetId         string    `json:"datasetId"`
	DatasetUri        string    `json:"datasetUri"`
	CDID              string    `json:"cdid,omitempty"`
	DataUnit          string    `json:"unit"`
	PreUnit           string    `json:"preUnit"`
	Source            string    `json:"source"`
	DataDate          string    `json:"date,omitempty"`
	DataNumber        string    `json:"number,omitempty"`
	KeyNote           string    `json:"keyNote,omitempty"`
	SampleSize        string    `json:"sampleSize,omitempty"`
	VersionLabel      string    `json:"versionLabel,omitempty"`
}

type Contact struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}
