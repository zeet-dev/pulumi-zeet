import * as zeet from "@pulumi/zeet-native";
import * as pulumi from "@pulumi/pulumi";


let config = new pulumi.Config();

let teamIdConfig = config.get("team-id");
let clusterId = config.get("cluster-id");
let githubRepoUrl = config.get("github-repo-url");

const project = new zeet.resources.Project("my-project", {
    userId: teamIdConfig,
    name: "pulumi-test-project-01",
})

const newProjectWithId = new zeet.resources.Project("second-project", {
    userId: teamIdConfig,
    name: "pulumi-test-project-02"
})

let existingProjectIdConfig = config.get("existing-project-id");

const existingProject = zeet.resources.Project.get("existing-project", existingProjectIdConfig)
const existingProject2 = zeet.resources.Project.get("existing-project2", existingProjectIdConfig)

const newProjectPreview = new zeet.resources.Project("new-project-for-preview", {
    userId: teamIdConfig,
    name: "pulumi-test-preview"
})

const environment = new zeet.resources.Environment("my-environment", {
    projectId: project.projectId,
    name: "pulumi-test-environment-01",
})

const app = new zeet.resources.App("github-app", {
    userId: teamIdConfig,
    projectId: project.projectId,
    environmentId: environment.environmentId,
    name: "pulumi-github-test-01",
    enabled: true,
    build: {
        type: "DOCKER",
        dockerfilePath: "Dockerfile",
    },
    resources: {
        cpu: "1",
        memory: "1"
    },
    deploy: {
        deployTarget: "KUBERNETES",
        // deployType: "KUBERNETES",
        // deployRuntime: "KUBERNETES",
        clusterId: clusterId
    },
    github: {
        url: githubRepoUrl,
        productionBranch: "main"
    },
    environmentVariables: [
        {
            name: "TEST_ENV",
            value: "1"
        },
        {
            name: "TEST_ENV_SEALED",
            value: "xyz",
            sealed: true
        }
    ]
})

export const output = "<obsolete>";

export const projectId = project.projectId;
export const projectPulumiId = project.id;
export const newProjectPulumiId = newProjectWithId.id;
export const existingProjectName = existingProject.name;
export const existingProjectId = existingProject2.id;
export const newProjectPreviewId = newProjectPreview.id;
export const environmentId = environment.environmentId;
