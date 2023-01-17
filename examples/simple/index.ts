import * as zeet from "@pulumi/zeet-native";
import * as pulumi from "@pulumi/pulumi";


const random = new zeet.Random("my-random", { length: 24 });

let config = new pulumi.Config();

let teamIdConfig = config.get("team-id");
let clusterId = config.get("cluster-id");
let githubRepoUrl = config.get("github-repo-url");

const project = new zeet.Project("my-project", {
    userId: teamIdConfig,
    name: "pulumi-test-project-01",
})

const environment = new zeet.Environment("my-environment", {
    projectId: project.projectId,
    name: "pulumi-test-environment-01",
})

const app = new zeet.App("github-app", {
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
    }
})

export const output = random.result;

export const projectId = project.projectId;
export const environmentId = environment.environmentId;
