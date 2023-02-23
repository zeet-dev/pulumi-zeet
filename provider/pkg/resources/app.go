package resources

import (
	"fmt"
	"time"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/config"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/gql"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/model"
)

type App struct{}

var _ = (infer.CustomCheck[AppArgs])((*App)(nil))
var _ = (infer.CustomRead[AppArgs, AppState])((*App)(nil))
var _ = (infer.CustomUpdate[AppArgs, AppState])((*App)(nil))
var _ = (infer.CustomDelete[AppState])((*App)(nil))

func (a App) Check(ctx provider.Context, name string, oldInputs resource.PropertyMap, newInputs resource.PropertyMap) (AppArgs, []provider.CheckFailure, error) {
	// check that one of the required "AppArgs" params is set
	if err := checkAppArgs(newInputs); err != nil {
		return AppArgs{}, []provider.CheckFailure{}, err
	}
	parsedArgs, i, err := infer.DefaultCheck[AppArgs](newInputs)
	if err != nil {
		return parsedArgs, i, err
	}
	if parsedArgs.BuildInput.Type == "DOCKER" {
		if parsedArgs.BuildInput.DockerfilePath == nil {
			return parsedArgs, i, fmt.Errorf("must specify DockerfilePath for DOCKER build type")
		}
	}
	return parsedArgs, i, nil
}

func checkAppArgs(newInputs resource.PropertyMap) error {
	if _, ok := newInputs["github"]; ok {
		return nil
	} else if _, ok := newInputs["docker"]; ok {
		return nil
	} else {
		return fmt.Errorf("must specify one of: github, docker")
	}
}

// Each resources has in input struct, defining what arguments it accepts.
type AppArgs struct {
	// Fields projected into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but its generally a
	// good idea.
	UserID               string                                    `pulumi:"userId"`
	ProjectID            string                                    `pulumi:"projectId"`
	EnvironmentID        string                                    `pulumi:"environmentId"`
	Name                 string                                    `pulumi:"name"`
	Enabled              bool                                      `pulumi:"enabled"`
	ResourcesInput       model.CreateAppResourcesInput             `pulumi:"resources"`
	BuildInput           model.CreateAppBuildInput                 `pulumi:"build,optional"`
	DeployInput          model.CreateAppDeployInput                `pulumi:"deploy"`
	GithubInput          *model.CreateAppGithubInput               `pulumi:"github,optional"`
	DockerInput          *model.CreateAppDockerInput               `pulumi:"docker,optional"`
	EnvironmentVariables []model.CreateAppEnvironmentVariableInput `pulumi:"environmentVariables,optional"`
}

// Each resources has a state, describing the fields that exist on the created resources.
type AppState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	AppArgs
	AppID     string    `pulumi:"appId"`
	UpdatedAt time.Time `pulumi:"updatedAt"`
}

// All resources must implement Create at a minumum.
func (App) Create(ctx provider.Context, name string, input AppArgs, preview bool) (string, AppState, error) {
	state := AppState{AppArgs: input}
	if preview {
		return state.AppID, state, nil
	}

	if input.GithubInput != nil {
		ctx.Log(diag.Warning,
			"app resources configuration is not currently functional for github-based apps, app wil default to 'tiny' instance size")
	}

	args := model.CreateAppInput{
		UserID:               input.UserID,
		ProjectID:            input.ProjectID,
		EnvironmentID:        input.EnvironmentID,
		Name:                 input.Name,
		Enabled:              input.Enabled,
		Resources:            input.ResourcesInput,
		Build:                input.BuildInput,
		Deploy:               input.DeployInput,
		GithubInput:          input.GithubInput,
		DockerInput:          input.DockerInput,
		EnvironmentVariables: input.EnvironmentVariables,
	}
	newApp, err := config.ZeetClient.CreateApp(ctx, args)
	if err != nil {
		return name, state, err
	}
	state.AppID = newApp.ID
	state.UpdatedAt = newApp.UpdatedAt
	return state.AppID, state, nil
}

func (a App) Read(ctx provider.Context, id string, inputs AppArgs, state AppState) (canonicalID string, normalizedInputs AppArgs, normalizedState AppState, err error) {
	resp, err := config.ZeetClient.ReadApp(ctx, id)
	if err != nil {
		return
	}
	normalizedState = responseToAppState(resp)
	return resp.ID, inputs, normalizedState, nil
}

func responseToAppState(resp gql.CreateAppResponse) AppState {
	state := AppState{
		AppArgs: AppArgs{
			UserID:        resp.UserID,
			ProjectID:     resp.ProjectID,
			EnvironmentID: resp.EnvironmentID,
			Name:          resp.Name,
			Enabled:       resp.Enabled,
			ResourcesInput: model.CreateAppResourcesInput{
				Cpu:              resp.Resources.Cpu,
				Memory:           resp.Resources.Memory,
				EphemeralStorage: resp.Resources.EphemeralStorage,
				SpotInstance:     resp.Resources.SpotInstance,
			},
			BuildInput: model.CreateAppBuildInput{
				Type:           resp.Build.Type,
				DockerfilePath: resp.Build.DockerfilePath,
			},
			DeployInput: model.CreateAppDeployInput{
				DeployTarget: resp.Deploy.DeployTarget,
				ClusterID:    resp.Deploy.ClusterID,
			},

			EnvironmentVariables: resp.EnvironmentVariables,
		},
		AppID:     resp.ID,
		UpdatedAt: resp.UpdatedAt,
	}

	if resp.GithubInput != nil {
		state.GithubInput = &model.CreateAppGithubInput{
			Url:              resp.GithubInput.Url,
			ProductionBranch: resp.GithubInput.ProductionBranch,
		}
	}

	if resp.DockerInput != nil {
		state.DockerInput = &model.CreateAppDockerInput{
			DockerImage: resp.DockerInput.DockerImage,
		}
	}

	return state
}

func (a App) Update(ctx provider.Context, id string, olds AppState, news AppArgs, preview bool) (AppState, error) {
	if preview {
		newState := olds
		newState.AppArgs = news
		return newState, nil
	}
	input := model.CreateAppInput{
		// update not possible yet
		UserID:               olds.UserID,
		ProjectID:            olds.ProjectID,
		EnvironmentID:        olds.EnvironmentID,
		GithubInput:          olds.GithubInput,
		Deploy:               olds.DeployInput,
		EnvironmentVariables: olds.EnvironmentVariables,

		// can be updated
		Name:    news.Name,
		Enabled: news.Enabled,
		Build:   news.BuildInput,
		Resources: model.CreateAppResourcesInput{
			Cpu:    news.ResourcesInput.Cpu,
			Memory: news.ResourcesInput.Memory,
		},
		DockerInput: news.DockerInput,
	}
	resp, err := config.ZeetClient.UpdateApp(ctx, id, input)
	if err != nil {
		return AppState{}, err
	}

	return responseToAppState(resp), nil
}

func (a App) Delete(ctx provider.Context, id string, props AppState) error {
	return config.ZeetClient.DeleteApp(ctx, id)
}
