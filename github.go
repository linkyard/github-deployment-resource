package resource

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"

	"github.com/shipt/go-github/v32/github"
)

//go:generate counterfeiter -o fakes/fake_git_hub.go . GitHub

type GitHub interface {
	ListDeployments(etag string) ([]*github.Deployment, string, error)
	ListDeploymentStatuses(ID int64) ([]*github.DeploymentStatus, error)
	GetDeployment(ID int64) (*github.Deployment, error)
	CreateDeployment(request *github.DeploymentRequest) (*github.Deployment, error)
	CreateDeploymentStatus(ID int64, request *github.DeploymentStatusRequest) (*github.DeploymentStatus, error)
}

type GitHubClient struct {
	client *github.Client

	user       string
	repository string
}

func NewGitHubClient(source Source) (*GitHubClient, error) {
	var client *github.Client

	client, err := oauthClient(source)
	if err != nil {
		return nil, err
	}

	if source.GitHubAPIURL != "" {
		client.BaseURL, err = url.Parse(source.GitHubAPIURL)
		if err != nil {
			return nil, err
		}
	}

	return &GitHubClient{
		client:     client,
		user:       source.User,
		repository: source.Repository,
	}, nil
}

func (g *GitHubClient) ListDeployments(etag string) ([]*github.Deployment, string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var opts *github.DeploymentsListOptions
	if etag != "" {
		opts = &github.DeploymentsListOptions{
			ETag: etag,
		}
	}
	deployments, res, err := g.client.Repositories.ListDeployments(ctx, g.user, g.repository, opts)

	if err != nil && res.StatusCode != http.StatusNotModified { // An error + a Not Modified status is actually a successful operation
		return []*github.Deployment{}, "", err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, "", err
	}

	return deployments, res.Header.Get("ETag"), nil
}

func (g *GitHubClient) GetDeployment(ID int64) (*github.Deployment, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	deployment, res, err := g.client.Repositories.GetDeployment(ctx, g.user, g.repository, ID)
	if err != nil {
		return &github.Deployment{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (g *GitHubClient) CreateDeployment(request *github.DeploymentRequest) (*github.Deployment, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	deployment, res, err := g.client.Repositories.CreateDeployment(ctx, g.user, g.repository, request)
	if err != nil {
		return &github.Deployment{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (g *GitHubClient) ListDeploymentStatuses(ID int64) ([]*github.DeploymentStatus, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	statuses, res, err := g.client.Repositories.ListDeploymentStatuses(ctx, g.user, g.repository, ID, nil)
	if err != nil {
		return []*github.DeploymentStatus{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (g *GitHubClient) CreateDeploymentStatus(ID int64, request *github.DeploymentStatusRequest) (*github.DeploymentStatus, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, res, err := g.client.Repositories.CreateDeploymentStatus(ctx, g.user, g.repository, ID, request)
	if err != nil {
		return &github.DeploymentStatus{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return status, nil
}

func oauthClient(source Source) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: source.AccessToken,
	})

	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)

	githubHTTPClient := &http.Client{
		Transport: oauthClient.Transport,
	}

	return github.NewClient(githubHTTPClient), nil
}
