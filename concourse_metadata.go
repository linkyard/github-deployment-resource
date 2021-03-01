package resource

import (
	"fmt"
	"os"
)

type ConcourseMetadata struct {
	ATCExternalURL    string `json:"atc_external_url"`
	BuildID           string `json:"build_id"`
	BuildJobName      string `json:"build_job_name"`
	BuildName         string `json:"build_name"`
	BuildPipelineName string `json:"build_pipeline_name"`
	BuildTeamName     string `json:"build_team_name"`
	BuildURL          string `json:"build_url"`
}

func GetConcourseMetadata() ConcourseMetadata {
	buildURL := fmt.Sprintf("%v/teams/%v/pipelines/%v/jobs/%v/builds/%v",
		os.Getenv("ATC_EXTERNAL_URL"), os.Getenv("BUILD_TEAM_NAME"), os.Getenv("BUILD_PIPELINE_NAME"), os.Getenv("BUILD_JOB_NAME"), os.Getenv("BUILD_NAME"))
	return ConcourseMetadata{
		ATCExternalURL:    os.Getenv("ATC_EXTERNAL_URL"),
		BuildID:           os.Getenv("BUILD_ID"),
		BuildJobName:      os.Getenv("BUILD_JOB_NAME"),
		BuildName:         os.Getenv("BUILD_NAME"),
		BuildPipelineName: os.Getenv("BUILD_PIPELINE_NAME"),
		BuildTeamName:     os.Getenv("BUILD_TEAM_NAME"),
		BuildURL:          buildURL,
	}
}
