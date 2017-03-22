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
	Summary           string    `json:"summary,omitempty"`
	Keywords          []string  `json:"keywords,omitempty"`
	MetaDescription   string    `json:"metaDescription,omitempty"`
	NationalStatistic bool      `json:"nationalStatistic,omitempty"`
	Contact           *Contact  `json:"contact"`
	ReleaseDate       time.Time `json:"releaseDate,string"`
	NextRelease       string    `json:"nextRelease"`
	Edition           string    `json:"edition,omitempty"`
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

type Data struct {
	Years            *[]Period    `json:"years"`
	Quarters         *[]Period    `json:"quarters"`
	Months           *[]Period    `json:"months"`
	RelatedDatasets  *[]Relation  `json:"relatedDatasets"`
	RelatedDocuments *[]Relation  `json:"relatedDocuments"`
	RelatedData      *[]Relation  `json:"relatedData,omitempty"`
	Versions         *[]Version   `json:"versions"`
	DataType         string       `json:"type"`
	DataUri          string       `json:"uri"`
	Description      *Description `json:"description"`
}

type Period struct {
	PeriodDate    string    `json:"date"`
	Value         string    `json:"value"`
	Label         string    `json:"label"`
	PeriodYear    string    `json:"year"`
	PeriodMonth   string    `json:"month"`
	Quarter       string    `json:"quarter"`
	SourceDataset string    `json:"sourceDataset"`
	UpdateDate    time.Time `json:"updateDate"`
}

type Relation struct {
	RelationUri string `json:"uri"`
}

type Version struct {
	VersionUri       string    `json:"uri"`
	UpdateDate       time.Time `json:"updateDate"`
	CorrectionNotice string    `json:"correctionNotice"`
	Label            string    `json:"label"`
}
