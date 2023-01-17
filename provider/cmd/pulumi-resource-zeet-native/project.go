package main

import (
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/config"
	"time"
)

type Project struct{}

// Each resources has in input struct, defining what arguments it accepts.
type ProjectArgs struct {
	// Fields projected into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but its generally a
	// good idea.
	UserID string `pulumi:"userId"`
	Name   string `pulumi:"name"`
}

// Each resources has a state, describing the fields that exist on the created resources.
type ProjectState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	ProjectArgs
	// Here we define a required output called result.
	ProjectID string    `pulumi:"projectId"`
	UpdatedAt time.Time `pulumi:"updatedAt"`
}

// All resources must implement Create at a minumum.
func (Project) Create(ctx provider.Context, name string, input ProjectArgs, preview bool) (string, ProjectState, error) {
	state := ProjectState{ProjectArgs: input}
	if preview {
		return name, state, nil
	}
	newProject, err := config.ZeetClient.CreateProject(ctx, input.UserID, input.Name)
	if err != nil {
		return name, state, err
	}
	state.ProjectID = newProject.ID
	state.UpdatedAt = newProject.UpdatedAt
	return name, state, nil
}
