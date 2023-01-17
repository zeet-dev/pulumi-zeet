package gql

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/model"
	"net/http"
	"time"
)

type ZeetGraphqlClient struct {
	endpoint string
	apiToken string
	client   graphql.Client
}

func NewZeetGraphqlClient(endpoint string, apiToken string) ZeetGraphqlClient {
	return ZeetGraphqlClient{
		endpoint: endpoint,
		apiToken: apiToken,
		client: graphql.NewClient(endpoint,
			&http.Client{Transport: &authedTransport{apiToken: apiToken, wrapped: http.DefaultTransport}}),
	}
}

type authedTransport struct {
	apiToken string
	wrapped  http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.apiToken)
	return t.wrapped.RoundTrip(req)
}

func (c *ZeetGraphqlClient) GetCurrentUserID(ctx context.Context) (string, error) {
	resp, err := currentUserID(ctx, c.client)
	if err != nil {
		return "", err
	}
	return resp.CurrentUser.Id, nil
}

type CreateProjectResponse struct {
	ID        string
	Name      string
	UpdatedAt time.Time
}

func (c *ZeetGraphqlClient) CreateProject(ctx context.Context, userID string, name string) (CreateProjectResponse, error) {
	resp, err := createProject(ctx, c.client, userID, name)
	if err != nil {
		return CreateProjectResponse{}, err
	}
	project := resp.CreateProjectV2
	return CreateProjectResponse{
		ID:        project.Id,
		Name:      project.Name,
		UpdatedAt: project.UpdatedAt,
	}, nil
}

type CreateEnvironmentResponse struct {
	ID        string
	Name      string
	ProjectID string
	UpdatedAt time.Time
}

func (c *ZeetGraphqlClient) CreateEnvironment(ctx provider.Context, projectID string, name string) (CreateEnvironmentResponse, error) {
	resp, err := createEnvironment(ctx, c.client, projectID, name)
	if err != nil {
		return CreateEnvironmentResponse{}, err
	}
	environment := resp.CreateProjectEnvironment
	return CreateEnvironmentResponse{
		ID:        environment.Id,
		Name:      environment.Name,
		ProjectID: environment.Project.Id,
		UpdatedAt: environment.UpdatedAt,
	}, nil
}

type CreateAppResponse struct {
	ID        string
	UpdatedAt time.Time
}

func (c *ZeetGraphqlClient) CreateApp(ctx provider.Context, args model.CreateAppInput) (CreateAppResponse, error) {
	buildType := BuildType(args.Build.Type)
	deployTarget := DeployTarget(args.Deploy.DeployTarget)
	input := CreateProjectGitInput{
		UserID:        &args.UserID,
		ProjectID:     &args.ProjectID,
		EnvironmentID: &args.EnvironmentID,
		Name:          &args.Name,
		Build: &ProjectBuildInput{
			BuildType:      &buildType,
			DockerfilePath: &args.Build.DockerfilePath,
		},
		DeployTarget: &ProjectDeployInput{
			DeployTarget: deployTarget,
			ClusterID:    &args.Deploy.ClusterID,
		},
	}
	if args.GithubInput != nil {
		input.Url = args.GithubInput.Url
		input.ProductionBranch = &args.GithubInput.ProductionBranch
	} else {
		return CreateAppResponse{}, fmt.Errorf("must specify one app spec")
	}
	resp, err := createAppGit(ctx, c.client, &input)
	if err != nil {
		return CreateAppResponse{}, err
	}
	return CreateAppResponse{
		ID:        resp.CreateProjectGit.Id,
		UpdatedAt: resp.CreateProjectGit.UpdatedAt,
	}, nil
}
