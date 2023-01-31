import * as zeet from "@pulumi/zeet-native";
import * as pulumi from "@pulumi/pulumi";


let config = new pulumi.Config();

let teamIdConfig = config.require("team-id");
let clusterId = config.require("cluster-id");
let githubRepoUrl = config.require("github-repo-url");

const project = new zeet.resources.Project("my-project", {
    userId: teamIdConfig,
    name: "pulumi-test-project-01",
})

// deleted
// const newProjectWithId = new zeet.resources.Project("second-project", {
//     userId: teamIdConfig,
//     name: "pulumi-test-project-02-renamed"
// })


let existingProjectIdConfig = config.get("existing-project-id");

// const existingProject = zeet.resources.Project.get("existing-project", existingProjectIdConfig)
// const existingProject2 = zeet.resources.Project.get("existing-project2", existingProjectIdConfig)

const newProjectPreview = new zeet.resources.Project("new-project-for-preview", {
    userId: teamIdConfig,
    name: "pulumi-test-preview"
})

const environment = new zeet.resources.Environment("my-environment", {
    projectId: project.projectId,
    name: "pulumi-test-environment-02-renamed",
})

// Deleted
// const environmentToDelete = new zeet.resources.Environment("env-to-delete", {
//     projectId: project.projectId,
//     name: "pulumi-env-to-delete"
// })

let existingEnvIdConfig = config.get("existing-env-id");

let existingEnvProjectIdConfig = config.get("existing-env-project-id");

// NB: .get() doesn't work when the resource id is insufficient to fetch the resource
// const existingEnvironment = zeet.resources.Environment.get("existing-env", existingEnvIdConfig, {
//     projectId: existingEnvProjectIdConfig
// } as any)

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
        cpu: 2.0,
        memory: 2.0,
        ephemeralStorage: 10.0,
        spotInstance: true,
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
export const newProjectPulumiId = "<deleted>";
export const newProjectPreviewId = newProjectPreview.id;
export const envDeleteId = "<deleted>";
export const environmentId = environment.environmentId;
export const appId = app.id;