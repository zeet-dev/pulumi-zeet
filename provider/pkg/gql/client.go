package gql

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/model"
	"net/http"
	"strconv"
	"time"
)

type zeetGraphqlClient struct {
	endpoint  string
	apiToken  string
	client    graphql.Client
	userAgent string
}

func NewZeetGraphqlClient(endpoint string, apiToken string, version string) ZeetClient {
	userAgent := fmt.Sprintf("Pulumi/3.0 (https://www.pulumi.com) pulumi-zeet-native/%s", version)
	transport := authedTransport{apiToken: apiToken, wrapped: http.DefaultTransport, userAgent: userAgent}
	println("version", version, "ua", userAgent)
	return &zeetGraphqlClient{
		endpoint: endpoint,
		apiToken: apiToken,
		client:   graphql.NewClient(endpoint, &http.Client{Transport: &transport}),
	}
}

type authedTransport struct {
	apiToken  string
	wrapped   http.RoundTripper
	userAgent string
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.apiToken)
	req.Header.Set("User-Agent", t.userAgent)
	return t.wrapped.RoundTrip(req)
}

func (c *zeetGraphqlClient) GetCurrentUserID(ctx context.Context) (string, error) {
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

func (c *zeetGraphqlClient) CreateProject(ctx context.Context, userID string, name string) (CreateProjectResponse, error) {
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

func (c *zeetGraphqlClient) ReadProject(ctx context.Context, projectID string) (CreateProjectResponse, error) {
	resp, err := getProjectByID(ctx, c.client, projectID)
	if err != nil {
		return CreateProjectResponse{}, err
	}
	return CreateProjectResponse{
		ID:        resp.Project.Id,
		Name:      resp.Project.Name,
		UpdatedAt: resp.Project.UpdatedAt,
	}, nil
}

func (c *zeetGraphqlClient) UpdateProject(ctx context.Context, projectID string, name *string) (CreateProjectResponse, error) {
	resp, err := updateProject(ctx, c.client, projectID, name)
	if err != nil {
		return CreateProjectResponse{}, err
	}
	return CreateProjectResponse{
		ID:        resp.UpdateProjectV2.Id,
		Name:      resp.UpdateProjectV2.Name,
		UpdatedAt: resp.UpdateProjectV2.UpdatedAt,
	}, nil
}

func (c *zeetGraphqlClient) DeleteProject(ctx provider.Context, projectID string) error {
	resp, err := deleteProject(ctx, c.client, projectID)
	// graphql error
	if err != nil {
		return err
	}
	// server returned false
	if !resp.GetDeleteProjectV2() {
		return fmt.Errorf("unable to delete project '%s'", projectID)
	}
	// server returned true
	return nil
}

type CreateEnvironmentResponse struct {
	ID        string
	Name      string
	ProjectID string
	UpdatedAt time.Time
}

func (c *zeetGraphqlClient) CreateEnvironment(ctx provider.Context, projectID string, name string) (CreateEnvironmentResponse, error) {
	resp, err := createEnvironment(ctx, c.client, projectID, name)
	if err != nil {
		return CreateEnvironmentResponse{}, err
	}
	environment := resp.CreateProjectEnvironment
	project := environment.GetProject()
	if project == nil {
		return CreateEnvironmentResponse{}, fmt.Errorf("expected to find project for environment '%s'", environment.GetId())
	}
	return CreateEnvironmentResponse{
		ID:        environment.GetId(),
		Name:      environment.GetName(),
		ProjectID: project.GetId(),
		UpdatedAt: environment.GetUpdatedAt(),
	}, nil
}

func (c *zeetGraphqlClient) ReadEnvironment(ctx context.Context, projectID string, environmentID string) (CreateEnvironmentResponse, error) {
	resp, err := getProjectEnvironments(ctx, c.client, projectID)
	if err != nil {
		return CreateEnvironmentResponse{}, err
	}
	proj := resp.GetProject()
	if proj == nil {
		return CreateEnvironmentResponse{}, fmt.Errorf("could not fetch project by id '%s'", projectID)
	}
	// find environment by id in array
	for _, environment := range proj.GetEnvironments() {
		if environment == nil {
			continue
		}
		if environment.GetId() == environmentID {
			return CreateEnvironmentResponse{
				ID:        environment.GetId(),
				Name:      environment.GetName(),
				ProjectID: proj.GetId(),
				UpdatedAt: environment.GetUpdatedAt(),
			}, nil
		}
	}
	// environment was not found in project
	return CreateEnvironmentResponse{}, fmt.Errorf("no environment with id '%s' found in project '%s'", environmentID, projectID)
}

func (c *zeetGraphqlClient) UpdateEnvironment(ctx context.Context, environmentID string, name *string) (CreateEnvironmentResponse, error) {
	resp, err := updateEnvironment(ctx, c.client, environmentID, name)
	if err != nil {
		return CreateEnvironmentResponse{}, err
	}
	updatedEnv := resp.GetUpdateProjectEnvironment()
	return CreateEnvironmentResponse{
		ID:        updatedEnv.GetId(),
		Name:      updatedEnv.GetName(),
		UpdatedAt: updatedEnv.GetUpdatedAt(),
	}, nil
}

func (c *zeetGraphqlClient) DeleteEnvironment(ctx provider.Context, environmentID string) error {
	resp, err := deleteEnvironment(ctx, c.client, environmentID)
	// graphql error
	if err != nil {
		return err
	}
	// server returned false
	if !resp.GetDeleteProjectEnvironment() {
		return fmt.Errorf("unable to delete environment '%s'", environmentID)
	}
	// server returned true
	return nil
}

func NewCreateAppResponse(args model.CreateAppInput, state appStateResponse) CreateAppResponse {
	return CreateAppResponse{
		CreateAppInput: args,
		ID:             state.GetId(),
		UpdatedAt:      state.GetUpdatedAt(),
	}
}

type appStateResponse interface {
	GetId() string
	GetUpdatedAt() time.Time
}

type CreateAppResponse struct {
	model.CreateAppInput
	ID        string
	UpdatedAt time.Time
}

func (c *zeetGraphqlClient) CreateApp(ctx provider.Context, args model.CreateAppInput) (CreateAppResponse, error) {
	var response appStateResponse
	if args.GithubInput != nil {
		input, err := newCreateProjectGitInput(args, *args.GithubInput)
		if err != nil {
			return CreateAppResponse{}, err
		}
		resp, err := createAppGit(ctx, c.client, &input)
		if err != nil {
			return CreateAppResponse{}, err
		}
		response = resp.GetCreateProjectGit()
	} else if args.DockerInput != nil {
		input := newCreateProjectDockerInput(args, *args.DockerInput)
		resp, err := createAppDocker(ctx, c.client, &input)
		if err != nil {
			return CreateAppResponse{}, err
		}
		response = resp.GetCreateProjectDocker()
	} else {
		return CreateAppResponse{}, fmt.Errorf("must specify one app spec")
	}
	return NewCreateAppResponse(args, response), nil
}

func newCreateProjectDockerInput(args model.CreateAppInput, dockerInput model.CreateAppDockerInput) CreateProjectDockerInput {
	input := CreateProjectDockerInput{
		Enabled:       &args.Enabled,
		UserID:        &args.UserID,
		ProjectID:     &args.ProjectID,
		EnvironmentID: &args.EnvironmentID,
		Name:          &args.Name,
		DeployTarget: &ProjectDeployInput{
			DeployTarget: getDeployTarget(args.Deploy),
			ClusterID:    &args.Deploy.ClusterID,
		},
		Envs:        environmentVariablesToRequestInput(args.EnvironmentVariables),
		DockerImage: dockerInput.DockerImage,
		Cpu:         args.GetCpuString(),
		Memory:      args.GetMemoryString(),
		Dedicated:   nil,
		Gpu:         nil,
		TeamID:      nil,
	}
	return input
}

func newCreateProjectGitInput(args model.CreateAppInput, githubInput model.CreateAppGithubInput) (CreateProjectGitInput, error) {
	buildType := getBuildType(args.Build)
	deployTarget := getDeployTarget(args.Deploy)
	memory, err := args.Resources.GetMemoryFloat()
	if err != nil {
		return CreateProjectGitInput{}, err
	}
	input := CreateProjectGitInput{
		UserID:        &args.UserID,
		ProjectID:     &args.ProjectID,
		EnvironmentID: &args.EnvironmentID,
		Name:          &args.Name,
		Build: &ProjectBuildInput{
			BuildType:      &buildType,
			DockerfilePath: args.Build.DockerfilePath,
		},
		DeployTarget: &ProjectDeployInput{
			DeployTarget: deployTarget,
			ClusterID:    &args.Deploy.ClusterID,
		},
		Envs: environmentVariablesToRequestInput(args.EnvironmentVariables),
		Resources: &ContainerResourcesSpecInput{
			Cpu: args.Resources.Cpu,
			// TODO: what kind of float is CreateProjectGitInput expecting? bytes?
			Memory:           memory,
			EphemeralStorage: args.Resources.EphemeralStorage,
			Spot:             args.Resources.SpotInstance,
		},
		Url:              githubInput.Url,
		ProductionBranch: githubInput.ProductionBranch,
	}
	return input, nil
}

func getDeployTarget(deployInput model.CreateAppDeployInput) DeployTarget {
	return DeployTarget(deployInput.DeployTarget)
}

func getBuildType(buildInput model.CreateAppBuildInput) BuildType {
	return BuildType(buildInput.Type)
}

func (c *zeetGraphqlClient) ReadApp(ctx provider.Context, appID string) (CreateAppResponse, error) {
	resp, err := getApp(ctx, c.client, appID)
	if err != nil {
		return CreateAppResponse{}, err
	}
	repo := resp.GetRepo()
	var cpu float64
	if repo.GetCpu() != nil {
		cpu, err = strconv.ParseFloat(*repo.GetCpu(), 64)
		if err != nil {
			return CreateAppResponse{}, fmt.Errorf("unable to parse float for cpu value '%s'", *repo.GetCpu())
		}
	} else {
		cpu = *new(float64)
	}
	var memory string
	if repo.GetMemory() != nil {
		memory = *repo.GetMemory()
	} else {
		memory = ""
	}
	args := model.CreateAppInput{
		UserID:        repo.GetOwner().GetId(),
		ProjectID:     repo.GetProject().GetId(),
		EnvironmentID: repo.GetProjectEnvironment().GetId(),
		Name:          repo.GetName(),
		GithubInput: &model.CreateAppGithubInput{
			// NB: anchor cannot resolve githubRepository for Repo's created with x-access-token:
			// ZEET-1480: https://linear.app/zeet/issue/ZEET-1480/anchor-or-createprojectgitgithubrepository-error-not-found
			ProductionBranch: repo.GetProductionBranch(),
		},
		Enabled: repo.GetEnabled(),
		Build: model.CreateAppBuildInput{
			Type:           string(repo.GetBuildMethod().GetType()),
			DockerfilePath: repo.GetBuildMethod().GetDockerfilePath(),
		},
		Resources: model.CreateAppResourcesInput{
			Cpu:    cpu,
			Memory: memory,
		},
		Deploy:               model.CreateAppDeployInput{},
		EnvironmentVariables: environmentVariablesToModel(repo.GetEnvs()),
	}
	return NewCreateAppResponse(args, resp.GetRepo()), nil
}

func (c *zeetGraphqlClient) UpdateApp(ctx provider.Context, appID string, args model.CreateAppInput) (CreateAppResponse, error) {
	cpuString := args.GetCpuString()
	memoryString := args.GetMemoryString()
	input := UpdateProjectInput{
		Id:               appID,
		Name:             &args.Name,
		DockerfilePath:   args.Build.DockerfilePath,
		Cpu:              cpuString,
		Memory:           memoryString,
		EphemeralStorage: args.Resources.EphemeralStorage,
		// TODO: can 'Spot' bool be updated?
	}
	resp, err := updateApp(ctx, c.client, &input)
	if err != nil {
		return CreateAppResponse{}, err
	}
	return NewCreateAppResponse(args, resp.GetUpdateProject()), nil
}

func (c *zeetGraphqlClient) DeleteApp(ctx provider.Context, appID string) error {
	resp, err := deleteApp(ctx, c.client, appID)
	if err != nil {
		return err
	}
	if !resp.GetDeleteRepo() {
		return fmt.Errorf("unable to delete app '%s'", appID)
	}
	return nil
}

func environmentVariablesToRequestInput(variables []model.CreateAppEnvironmentVariableInput) []*EnvVarInput {
	out := []*EnvVarInput{}
	for _, variable := range variables {
		input := &EnvVarInput{
			Name:   variable.Name,
			Value:  variable.Value,
			Sealed: variable.Sealed,
		}
		out = append(out, input)
	}
	return out
}

func environmentVariablesToModel(variables []*AppStateFragmentEnvsEnvVar) []model.CreateAppEnvironmentVariableInput {
	out := []model.CreateAppEnvironmentVariableInput{}
	for _, variable := range variables {
		sealed := variable.GetSealed()
		envVar := model.CreateAppEnvironmentVariableInput{
			Name:   variable.GetName(),
			Value:  variable.GetValue(),
			Sealed: &sealed,
		}
		out = append(out, envVar)
	}
	return out
}
