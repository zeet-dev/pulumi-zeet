package gql

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/model"
)

type zeetGraphqlClient struct {
	apiToken string
	client   graphql.Client
}

func NewZeetGraphqlClient(endpoint string, apiToken string, version string) ZeetClient {
	userAgent := fmt.Sprintf("Pulumi/3.0 (https://www.pulumi.com) pulumi-zeet/%s", version)
	transport := authedTransport{apiToken: apiToken, wrapped: http.DefaultTransport, userAgent: userAgent}
	graphqlEndpoint := endpoint + "/graphql"
	return &zeetGraphqlClient{
		apiToken: apiToken,
		client:   graphql.NewClient(graphqlEndpoint, &http.Client{Transport: &transport}),
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
		// idempotent delete - return nil if already deleted
		if strings.Contains(err.Error(), "deleteProjectV2 record not found") {
			return nil
		}
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
		// idempotent delete - return nil if already deleted
		if strings.Contains(err.Error(), "deleteProjectEnvironment record not found") {
			return nil
		}
		return err
	}
	// server returned false
	if !resp.GetDeleteProjectEnvironment() {
		return fmt.Errorf("unable to delete environment '%s'", environmentID)
	}
	// server returned true
	return nil
}

func NewCreateAppResponse(args model.CreateAppInput, state AppStateFragment) CreateAppResponse {
	return CreateAppResponse{
		CreateAppInput: args,
		ID:             state.GetId(),
		UpdatedAt:      state.GetUpdatedAt(),
	}
}

type CreateAppResponse struct {
	model.CreateAppInput
	ID        string
	UpdatedAt time.Time
}

func (c *zeetGraphqlClient) CreateApp(ctx provider.Context, args model.CreateAppInput) (CreateAppResponse, error) {
	var response AppStateFragment
	if args.GithubInput != nil {
		input, err := newCreateProjectGitInput(args, *args.GithubInput)
		if err != nil {
			return CreateAppResponse{}, err
		}
		resp, err := createAppGit(ctx, c.client, &input)
		if err != nil {
			return CreateAppResponse{}, err
		}
		response = resp.GetCreateProjectGit().AppStateFragment
	} else if args.DockerInput != nil {
		input := newCreateProjectDockerInput(args, *args.DockerInput)
		resp, err := createAppDocker(ctx, c.client, &input)
		if err != nil {
			return CreateAppResponse{}, err
		}
		response = resp.GetCreateProjectDocker().AppStateFragment
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

func CreateAppResponseFromAppState(state AppStateFragment) (CreateAppResponse, error) {
	var err error
	var cpu float64
	if state.GetCpu() != nil {
		cpu, err = strconv.ParseFloat(*state.GetCpu(), 64)
		if err != nil {
			return CreateAppResponse{}, fmt.Errorf("unable to parse float for cpu value '%s'", *state.GetCpu())
		}
	} else {
		cpu = *new(float64)
	}
	var memory string
	if state.GetMemory() != nil {
		memory = *state.GetMemory()
	} else {
		memory = ""
	}

	var spot bool
	if state.GetDedicated() != nil {
		spot = !*state.GetDedicated()
	} else {
		spot = false
	}

	var ephemeralStorage float64
	if state.GetEphemeralStorage() != nil {
		ephemeralStorage = *state.GetEphemeralStorage()
	} else {
		ephemeralStorage = 0
	}

	args := model.CreateAppInput{
		UserID:        state.GetOwner().GetId(),
		ProjectID:     state.GetProject().GetId(),
		EnvironmentID: state.GetProjectEnvironment().GetId(),
		Name:          state.GetName(),
		GithubInput: &model.CreateAppGithubInput{
			// NB: anchor cannot resolve githubRepository for Repo's created with x-access-token:
			// ZEET-1480: https://linear.app/zeet/issue/ZEET-1480/anchor-or-createprojectgitgithubrepository-error-not-found
			ProductionBranch: state.GetProductionBranch(),
		},
		Enabled: state.GetEnabled(),
		Resources: model.CreateAppResourcesInput{
			Cpu:              cpu,
			Memory:           memory,
			SpotInstance:     &spot,
			EphemeralStorage: &ephemeralStorage,
		},
		Deploy:               model.CreateAppDeployInput{},
		EnvironmentVariables: environmentVariablesToModel(state.GetEnvs()),
	}

	if state.GetBuildMethod() != nil {
		args.Build = model.CreateAppBuildInput{
			Type:           string(state.GetBuildMethod().GetType()),
			DockerfilePath: state.GetBuildMethod().GetDockerfilePath(),
		}
	}

	if state.GetSource() != nil {
		switch state.GetSource().GetType() {
		case RepoSourceTypeDocker:
			args.DockerInput = &model.CreateAppDockerInput{
				DockerImage: state.GetSource().GetId(),
			}
		}
	}

	return NewCreateAppResponse(args, state), nil
}

func (c *zeetGraphqlClient) ReadApp(ctx provider.Context, appID string) (CreateAppResponse, error) {
	resp, err := getApp(ctx, c.client, appID)
	if err != nil {
		return CreateAppResponse{}, err
	}
	return CreateAppResponseFromAppState(resp.GetRepo().AppStateFragment)
}

func (c *zeetGraphqlClient) UpdateApp(ctx provider.Context, appID string, args model.CreateAppInput) (CreateAppResponse, error) {
	cpuString := args.GetCpuString()
	memoryString := args.GetMemoryString()
	input := UpdateProjectInput{
		Id:               appID,
		Name:             &args.Name,
		Cpu:              cpuString,
		Memory:           memoryString,
		EphemeralStorage: args.Resources.EphemeralStorage,
		// TODO: can 'Spot' bool be updated?
	}
	resp, err := updateApp(ctx, c.client, &input)
	if err != nil {
		return CreateAppResponse{}, err
	}
	var state AppStateFragment
	state = resp.GetUpdateProject().AppStateFragment
	if state.GetEnabled() != args.Enabled {
		if args.Enabled {
			resp, err := enableApp(ctx, c.client, appID)
			if err != nil {
				return CreateAppResponse{}, err
			}
			state = resp.GetEnableRepo().AppStateFragment
		} else {
			resp, err := disableApp(ctx, c.client, appID)
			if err != nil {
				return CreateAppResponse{}, err
			}
			state = resp.GetDisableRepo().AppStateFragment
		}
	}
	return NewCreateAppResponse(args, state), nil
}

func (c *zeetGraphqlClient) DeleteApp(ctx provider.Context, appID string) error {
	resp, err := deleteApp(ctx, c.client, appID)
	if err != nil {
		// idempotent delete - if app is already deleted, return nil
		if strings.Contains(err.Error(), "deleteRepo record not found") {
			return nil
		}
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
