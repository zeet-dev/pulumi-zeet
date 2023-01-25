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
	ClusterID    string `pulumi:"clusterId,optional"`
	//DeployType    DeployType    `pulumi:"deployType"`
	//DeployRuntime DeployRuntime `pulumi:"deployRuntime"`
}

type CreateAppGithubInput struct {
	Url              string `pulumi:"url"`
	ProductionBranch string `pulumi:"productionBranch"`
}

type CreateAppEnvironmentVariableInput struct {
	Name   string `pulumi:"name"`
	Value  string `pulumi:"value"`
	Sealed *bool  `pulumi:"sealed,optional"`
}

type CreateAppInput struct {
	UserID               string
	ProjectID            string
	EnvironmentID        string
	Name                 string
	GithubInput          *CreateAppGithubInput
	Enabled              bool
	Build                CreateAppBuildInput
	Resources            CreateAppResourcesInput
	Deploy               CreateAppDeployInput
	EnvironmentVariables []CreateAppEnvironmentVariableInput
}
