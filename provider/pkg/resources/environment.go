package resources

import (
	"fmt"
	"github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/config"
	"time"
)

type Environment struct{}

var _ = (infer.CustomRead[EnvironmentArgs, EnvironmentState])((*Environment)(nil))
var _ = (infer.CustomUpdate[EnvironmentArgs, EnvironmentState])((*Environment)(nil))
var _ = (infer.CustomDelete[EnvironmentState])((*Environment)(nil))

// Each resources has in input struct, defining what arguments it accepts.
type EnvironmentArgs struct {
	ProjectID string `pulumi:"projectId"`
	Name      string `pulumi:"name"`
}

// Each resources has a state, describing the fields that exist on the created resources.
type EnvironmentState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	EnvironmentArgs
	// Here we define a required output called result.
	EnvironmentID string    `pulumi:"environmentId"`
	UpdatedAt     time.Time `pulumi:"updatedAt"`
}

// All resources must implement Create at a minumum.
func (Environment) Create(ctx provider.Context, name string, input EnvironmentArgs, preview bool) (string, EnvironmentState, error) {
	state := EnvironmentState{EnvironmentArgs: input}
	if preview {
		return state.EnvironmentID, state, nil
	}
	newEnvironment, err := config.ZeetClient.CreateEnvironment(ctx, input.ProjectID, input.Name)
	if err != nil {
		return name, state, err
	}
	state.EnvironmentID = newEnvironment.ID
	state.ProjectID = newEnvironment.ProjectID
	state.UpdatedAt = newEnvironment.UpdatedAt
	return state.EnvironmentID, state, nil
}

func (e Environment) Read(ctx provider.Context, id string, inputs EnvironmentArgs, state EnvironmentState) (canonicalID string, normalizedInputs EnvironmentArgs, normalizedState EnvironmentState, err error) {
	if inputs.ProjectID == "" {
		err = fmt.Errorf("must specify projectId to read environment")
		return
	}
	environment, err := config.ZeetClient.ReadEnvironment(ctx, inputs.ProjectID, id)
	if err != nil {
		return "", EnvironmentArgs{}, EnvironmentState{}, err
	}
	normalizedInputs = inputs
	normalizedState.EnvironmentArgs = normalizedInputs
	normalizedState.Name = environment.Name
	normalizedState.EnvironmentID = environment.ID
	normalizedState.UpdatedAt = environment.UpdatedAt

	canonicalID = state.EnvironmentID
	return
}

func (e Environment) Update(ctx provider.Context, id string, olds EnvironmentState, news EnvironmentArgs, preview bool) (EnvironmentState, error) {
	newState := EnvironmentState{
		EnvironmentArgs: news,
		EnvironmentID:   id,
		UpdatedAt:       time.Time{},
	}
	if preview {
		return newState, nil
	}
	updated, err := config.ZeetClient.UpdateEnvironment(ctx, id, &news.Name)
	if err != nil {
		return newState, err
	}
	newState.UpdatedAt = updated.UpdatedAt
	return newState, nil
}

func (e Environment) Delete(ctx provider.Context, id string, props EnvironmentState) error {
	return config.ZeetClient.DeleteEnvironment(ctx, id)
}
