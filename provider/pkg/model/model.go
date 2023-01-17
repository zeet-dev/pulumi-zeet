package model

type CreateAppResourcesInput struct {
	Cpu    string `pulumi:"cpu"`
	Memory string `pulumi:"memory"`
}

type CreateAppBuildInput struct {
	Type           string `pulumi:"type"` // DOCKER
	DockerfilePath string `pulumi:"dockerfilePath,optional"`
}

type CreateAppDeployInput struct {
	DeployTarget string `pulumi:"deployTarget"`
	//DeployType    DeployType    `pulumi:"deployType"`
	//DeployRuntime DeployRuntime `pulumi:"deployRuntime"`
	ClusterID string `pulumi:"clusterId,optional"`
}

type CreateAppGithubInput struct {
	Url              string `pulumi:"url"`
	ProductionBranch string `pulumi:"productionBranch"`
}

type CreateAppInput struct {
	UserID        string
	ProjectID     string
	EnvironmentID string
	Name          string
	GithubInput   *CreateAppGithubInput
	Enabled       bool
	Build         CreateAppBuildInput
	Resources     CreateAppResourcesInput
	Deploy        CreateAppDeployInput
}
