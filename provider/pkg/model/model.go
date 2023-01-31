package model

import (
	"regexp"
	"strconv"
)

type CreateAppResourcesInput struct {
	Cpu              float64  `pulumi:"cpu"`
	Memory           string   `pulumi:"memory"`
	EphemeralStorage *float64 `pulumi:"ephemeralStorage"`
	SpotInstance     *bool    `pulumi:"spotInstance"`
}

var nonNumericRegex = regexp.MustCompile(`[^0-9\\.]+`)

func (r *CreateAppResourcesInput) GetMemoryFloat() (float64, error) {
	numeric := nonNumericRegex.ReplaceAllString(r.Memory, "")
	float, err := strconv.ParseFloat(numeric, 64)
	if err != nil {
		return 0.0, err
	}
	return float, nil
}

type CreateAppBuildInput struct {
	Type           string  `pulumi:"type"` // DOCKER
	DockerfilePath *string `pulumi:"dockerfilePath,optional"`
}

type CreateAppDeployInput struct {
	DeployTarget string `pulumi:"deployTarget"`
	ClusterID    string `pulumi:"clusterId,optional"`
	//DeployType    DeployType    `pulumi:"deployType"`
	//DeployRuntime DeployRuntime `pulumi:"deployRuntime"`
}

type CreateAppGithubInput struct {
	Url              string  `pulumi:"url"`
	ProductionBranch *string `pulumi:"productionBranch"`
}

type CreateAppEnvironmentVariableInput struct {
	Name   string `pulumi:"name"`
	Value  string `pulumi:"value"`
	Sealed *bool  `pulumi:"sealed,optional"`
}

type CreateAppDockerInput struct {
	DockerImage string `pulumi:"dockerImage"`
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
	DockerInput          *CreateAppDockerInput
}

func (i CreateAppInput) GetCpuString() *string {
	cpuString := strconv.FormatFloat(i.Resources.Cpu, 'f', -1, 64)
	return &cpuString
}

func (i CreateAppInput) GetMemoryString() *string {
	return &i.Resources.Memory
}
