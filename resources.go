package resource

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterbourgon/mergemap"
	"github.com/shipt/go-github/v32/github"
)

type Source struct {
	User         string   `json:"user"`
	Repository   string   `json:"repository"`
	AccessToken  string   `json:"access_token"`
	GitHubAPIURL string   `json:"github_api_url"`
	Environments []string `json:"environments"`
}

type Version struct {
	ID       string `json:"id"`
	Statuses string `json:"status"`
	ETag     string `json:"etag"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

type InResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

type OutResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

type OutParams struct {
	Type           *string `json:"type"`
	ID             *string
	Ref            *string
	Environment    *string
	Task           *string
	State          *string
	Description    *string
	AutoMerge      *bool
	Payload        *map[string]interface{}
	PayloadPath    *string `json:"payload_path"`
	LogURL         *string
	EnvironmentURL *string

	RawID             json.RawMessage `json:"id"`
	RawState          json.RawMessage `json:"state"`
	RawRef            json.RawMessage `json:"ref"`
	RawTask           json.RawMessage `json:"task"`
	RawEnvironment    json.RawMessage `json:"environment"`
	RawEnvironmentURL json.RawMessage `json:"environment_url"`
	RawDescription    json.RawMessage `json:"description"`
	RawAutoMerge      json.RawMessage `json:"auto_merge"`
	RawPayload        json.RawMessage `json:"payload"`
	RawLogURL         json.RawMessage `json:"log_url"`
}

// Used to avoid recursion in UnmarshalJSON below.
type outParams OutParams

func (p *OutParams) UnmarshalJSON(b []byte) (err error) {
	j := outParams{
		Type: github.String("status"),
	}

	if err = json.Unmarshal(b, &j); err == nil {
		*p = OutParams(j)
		if p.RawID != nil {
			p.ID = github.String(getStringOrStringFromFile(p.RawID))
		}

		if p.RawState != nil {
			p.State = github.String(getStringOrStringFromFile(p.RawState))
		}

		if p.RawRef != nil {
			p.Ref = github.String(getStringOrStringFromFile(p.RawRef))
		}

		if p.RawTask != nil {
			p.Task = github.String(getStringOrStringFromFile(p.RawTask))
		}

		if p.RawEnvironment != nil {
			p.Environment = github.String(getStringOrStringFromFile(p.RawEnvironment))
		}

		if p.RawEnvironmentURL != nil {
			envUrl := os.ExpandEnv(getStringOrStringFromFile(p.RawEnvironmentURL)) // Interpolate ENV variables
			p.EnvironmentURL = github.String(envUrl)
		}

		if p.RawDescription != nil {
			p.Description = github.String(getStringOrStringFromFile(p.RawDescription))
		}

		if p.RawLogURL != nil {
			logURL := os.ExpandEnv(getStringOrStringFromFile(p.RawLogURL)) // Interpolate ENV variables
			p.LogURL = github.String(logURL)
		}

		if p.RawAutoMerge != nil {
			p.AutoMerge = github.Bool(getBoolOrDefault(p.RawAutoMerge, false))
		}

		var payload map[string]interface{}
		json.Unmarshal(p.RawPayload, &payload)

		if p.PayloadPath != nil && *p.PayloadPath != "" {
			stringFromFile := fileContents(*p.PayloadPath)
			var payloadFromFile map[string]interface{}
			json.Unmarshal([]byte(stringFromFile), &payloadFromFile)

			payload = mergemap.Merge(payloadFromFile, payload)
		}

		p.Payload = &payload

		return
	}
	return
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewCheckRequest() CheckRequest {
	return CheckRequest{}
}

func NewInRequest() InRequest {
	return InRequest{}
}

func NewOutRequest() OutRequest {
	return OutRequest{}
}

func getStringOrStringFromFile(field json.RawMessage) string {
	var rawValue interface{}
	if err := json.Unmarshal(field, &rawValue); err == nil {
		switch rawValue := rawValue.(type) {
		case string:
			return rawValue
		case map[string]interface{}:
			return fileContents(rawValue["file"].(string))
		default:
			panic("Could not read string out of Params field")
		}
	}
	return ""
}

func getBoolOrDefault(field json.RawMessage, defaultValue bool) bool {
	var rawValue interface{}
	err := json.Unmarshal(field, &rawValue)
	if err != nil {
		log.Printf("Error unmarshalling json: %s.\nDefaulting to given default value %t", err, defaultValue)
		return defaultValue

	}
	switch rawValue := rawValue.(type) {
	case bool:
		return rawValue
	default:
		log.Printf("Could not parse bool from json; defaulting to false")
		return defaultValue
	}
}

func fileContents(path string) string {
	sourceDir := os.Args[1]
	contents, err := ioutil.ReadFile(filepath.Join(sourceDir, path))
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(contents))
}
