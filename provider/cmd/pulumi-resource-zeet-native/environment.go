package main

import (
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/config"
	"time"
)

type Environment struct{}

// Each resources has in input struct, defining what arguments it accepts.
type EnvironmentArgs struct {
	// Fields environmented into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but its generally a
	// good idea.
	ProjectID string `pulumi:"projectId"`
	Name      string `pulumi:"name"`
}

// Each resources has a state, describing the fields that exist on the created resources.
type EnvironmentState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	EnvironmentArgs
	// Here we define a required output called result.
	EnvironmentID string    `pulumi:"environmentId"`
	ProjectID     string    `pulumi:"projectId"`
	UpdatedAt     time.Time `pulumi:"updatedAt"`
}

// All resources must implement Create at a minumum.
func (Environment) Create(ctx provider.Context, name string, input EnvironmentArgs, preview bool) (string, EnvironmentState, error) {
	state := EnvironmentState{EnvironmentArgs: input}
	if preview {
		return name, state, nil
	}
	newEnvironment, err := config.ZeetClient.CreateEnvironment(ctx, input.ProjectID, input.Name)
	if err != nil {
		return name, state, err
	}
	state.EnvironmentID = newEnvironment.ID
	state.ProjectID = newEnvironment.ProjectID
	state.UpdatedAt = newEnvironment.UpdatedAt
	return name, state, nil
}
