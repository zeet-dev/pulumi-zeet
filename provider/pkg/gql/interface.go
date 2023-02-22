package gql

import (
	"context"

	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/model"
)

type ZeetClient interface {
	GetCurrentUserID(ctx context.Context) (string, error)
	CreateProject(ctx context.Context, userID string, name string) (CreateProjectResponse, error)
	ReadProject(ctx context.Context, projectID string) (CreateProjectResponse, error)
	UpdateProject(ctx context.Context, projectID string, name *string) (CreateProjectResponse, error)
	DeleteProject(ctx provider.Context, projectID string) error
	CreateEnvironment(ctx provider.Context, projectID string, name string) (CreateEnvironmentResponse, error)
	ReadEnvironment(ctx context.Context, projectID string, environmentID string) (CreateEnvironmentResponse, error)
	UpdateEnvironment(ctx context.Context, environmentID string, name *string) (CreateEnvironmentResponse, error)
	DeleteEnvironment(ctx provider.Context, environmentID string) error
	CreateApp(ctx provider.Context, args model.CreateAppInput) (CreateAppResponse, error)
	ReadApp(ctx provider.Context, appID string) (CreateAppResponse, error)
	UpdateApp(ctx provider.Context, appID string, args model.CreateAppInput) (CreateAppResponse, error)
	DeleteApp(ctx provider.Context, appID string) error
}
