import * as pulumi from "@pulumi/pulumi";
import * as zeet from "@zeet-dev/pulumi-zeet";

let config = new pulumi.Config();

let userId = config.require("user-id");
let clusterId = config.require("cluster-id");

const image = config.get("image") || "nginx:latest";

const project = new zeet.resources.Project("simple-apps-project", {
  userId: userId,
  name: "pulumi-simple-apps",
});

const environment = new zeet.resources.Environment("simple-apps-environment", {
  projectId: project.projectId,
  name: "staging",
});

for (let i = 0; i < 5; i++) {
  const dockerApp = new zeet.resources.App(`simple-app-${i}`, {
    userId: userId,
    projectId: project.projectId,
    environmentId: environment.environmentId,
    name: `simple-app-${i}`,
    enabled: false,
    docker: {
      dockerImage: image,
    },
    deploy: {
      deployTarget: "KUBERNETES",
      clusterId: clusterId,
    },
    resources: {
      cpu: 0.1,
      memory: "0.1Gi",
      ephemeralStorage: 0,
      spotInstance: true,
    },
    environmentVariables: [
      {
        name: "simple_APP_ID",
        value: `${i}`,
      },
    ],
  });
}

// Export the id of the project
export const projectId = project.id;
