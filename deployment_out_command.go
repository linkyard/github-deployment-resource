package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/shipt/go-github/v32/github"
)

type DeploymentOutCommand struct {
	github GitHub
	writer io.Writer
}

func NewDeploymentOutCommand(github GitHub, writer io.Writer) *DeploymentOutCommand {
	return &DeploymentOutCommand{
		github: github,
		writer: writer,
	}
}

func (c *DeploymentOutCommand) Run(sourceDir string, request OutRequest) (OutResponse, error) {
	if request.Params.Ref == nil {
		return OutResponse{}, errors.New("ref is a required parameter")
	}

	newDeployment := &github.DeploymentRequest{
		Ref:              request.Params.Ref,
		RequiredContexts: &[]string{},
	}

	concoursePayload := GetConcourseMetadata()

	if request.Params.Payload != nil && *request.Params.Payload != nil {
		payload := *request.Params.Payload
		payload["concourse_payload"] = concoursePayload
	} else {
		request.Params.Payload = &map[string]interface{}{
			"concourse_payload": concoursePayload,
		}
	}
	p, err := json.Marshal(request.Params.Payload)
	newDeployment.Payload = github.String(string(p))

	if request.Params.Task != nil {
		newDeployment.Task = request.Params.Task
	}
	if request.Params.Environment != nil {
		newDeployment.Environment = request.Params.Environment
	}
	if request.Params.Description != nil {
		newDeployment.Description = request.Params.Description
	}

	if request.Params.AutoMerge != nil {
		newDeployment.AutoMerge = request.Params.AutoMerge
	}

	fmt.Fprintln(c.writer, "creating deployment")
	deployment, err := c.github.CreateDeployment(newDeployment)
	if err != nil {
		return OutResponse{}, err
	}

	return OutResponse{
		Version:  Version{ID: strconv.FormatInt(*deployment.ID, 10)},
		Metadata: metadataFromDeployment(deployment, []*github.DeploymentStatus{}),
	}, nil
}

func (c *DeploymentOutCommand) fileContents(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
