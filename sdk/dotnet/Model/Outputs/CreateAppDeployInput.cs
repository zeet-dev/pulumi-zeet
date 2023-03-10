// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Zeet.Model.Outputs
{

    [OutputType]
    public sealed class CreateAppDeployInput
    {
        public readonly string? ClusterId;
        public readonly string DeployTarget;

        [OutputConstructor]
        private CreateAppDeployInput(
            string? clusterId,

            string deployTarget)
        {
            ClusterId = clusterId;
            DeployTarget = deployTarget;
        }
    }
}
