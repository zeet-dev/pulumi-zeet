// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";

export namespace model {
    export interface CreateAppBuildInput {
        dockerfilePath?: string;
        type: string;
    }

    export interface CreateAppDeployInput {
        clusterId?: string;
        deployTarget: string;
    }

    export interface CreateAppDockerInput {
        dockerImage: string;
    }

    export interface CreateAppEnvironmentVariableInput {
        name: string;
        sealed?: boolean;
        value: string;
    }

    export interface CreateAppGithubInput {
        productionBranch: string;
        url: string;
    }

    export interface CreateAppResourcesInput {
        cpu: number;
        ephemeralStorage?: number;
        memory: string;
        spotInstance?: boolean;
    }

}

export namespace time {
    export interface Time {
    }

}
