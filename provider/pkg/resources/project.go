package resources

import (
	"time"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/config"
)

type Project struct{}

var _ = (infer.CustomRead[ProjectArgs, ProjectState])((*Project)(nil))
var _ = (infer.CustomUpdate[ProjectArgs, ProjectState])((*Project)(nil))
var _ = (infer.CustomDelete[ProjectState])((*Project)(nil))

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
		return state.ProjectID, state, nil
	}
	newProject, err := config.ZeetClient.CreateProject(ctx, input.UserID, input.Name)
	if err != nil {
		return name, state, err
	}
	state.ProjectID = newProject.ID
	state.UpdatedAt = newProject.UpdatedAt
	return state.ProjectID, state, nil
}

func (p Project) Read(ctx provider.Context, id string, inputs ProjectArgs, state ProjectState) (canonicalID string, normalizedInputs ProjectArgs, normalizedState ProjectState, err error) {
	remoteState, err := config.ZeetClient.ReadProject(ctx, id)
	if err != nil {
		return "", ProjectArgs{}, ProjectState{}, err
	}
	normalizedState.UserID = inputs.UserID

	normalizedState.ProjectID = remoteState.ID
	normalizedState.Name = remoteState.Name
	normalizedState.UpdatedAt = remoteState.UpdatedAt

	return normalizedState.ProjectID, inputs, normalizedState, nil
}

func (p Project) Update(ctx provider.Context, id string, olds ProjectState, news ProjectArgs, preview bool) (ProjectState, error) {
	newState := ProjectState{
		ProjectArgs: news,
		ProjectID:   olds.ProjectID,
		UpdatedAt:   time.Time{},
	}
	if preview {
		return newState, nil
	}
	updated, err := config.ZeetClient.UpdateProject(ctx, id, &news.Name)
	if err != nil {
		return newState, err
	}
	newState.UpdatedAt = updated.UpdatedAt
	return newState, nil
}

func (p Project) Delete(ctx provider.Context, id string, props ProjectState) error {
	return config.ZeetClient.DeleteProject(ctx, id)
}
