package resources

import (
	"fmt"
	"github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/config"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/model"
	"time"
)

type App struct{}

var _ = (infer.CustomCheck[AppArgs])((*App)(nil))

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
		if parsedArgs.BuildInput.DockerfilePath == "" {
			return parsedArgs, i, fmt.Errorf("must specify DockerfilePath for DOCKER build type")
		}
	}
	return parsedArgs, i, nil
}

func checkAppArgs(newInputs resource.PropertyMap) error {
	if _, ok := newInputs["github"]; ok {
		return nil
	} else {
		return fmt.Errorf("must specify one of: github")
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
	BuildInput           model.CreateAppBuildInput                 `pulumi:"build"`
	DeployInput          model.CreateAppDeployInput                `pulumi:"deploy"`
	GithubInput          model.CreateAppGithubInput                `pulumi:"github,optional"`
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
		return name, state, nil
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
		GithubInput:          &input.GithubInput,
		EnvironmentVariables: input.EnvironmentVariables,
	}
	newApp, err := config.ZeetClient.CreateApp(ctx, args)
	if err != nil {
		return name, state, err
	}
	state.AppID = newApp.ID
	state.UpdatedAt = newApp.UpdatedAt
	return name, state, nil
}
