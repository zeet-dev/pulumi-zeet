// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Zeet.Resources
{
    [ZeetResourceType("zeet:resources:App")]
    public partial class App : global::Pulumi.CustomResource
    {
        [Output("appId")]
        public Output<string> AppId { get; private set; } = null!;

        [Output("build")]
        public Output<Pulumi.Zeet.Model.Outputs.CreateAppBuildInput?> Build { get; private set; } = null!;

        [Output("deploy")]
        public Output<Pulumi.Zeet.Model.Outputs.CreateAppDeployInput> Deploy { get; private set; } = null!;

        [Output("docker")]
        public Output<Pulumi.Zeet.Model.Outputs.CreateAppDockerInput?> Docker { get; private set; } = null!;

        [Output("enabled")]
        public Output<bool> Enabled { get; private set; } = null!;

        [Output("environmentId")]
        public Output<string> EnvironmentId { get; private set; } = null!;

        [Output("environmentVariables")]
        public Output<ImmutableArray<Pulumi.Zeet.Model.Outputs.CreateAppEnvironmentVariableInput>> EnvironmentVariables { get; private set; } = null!;

        [Output("github")]
        public Output<Pulumi.Zeet.Model.Outputs.CreateAppGithubInput?> Github { get; private set; } = null!;

        [Output("name")]
        public Output<string> Name { get; private set; } = null!;

        [Output("projectId")]
        public Output<string> ProjectId { get; private set; } = null!;

        [Output("resources")]
        public Output<Pulumi.Zeet.Model.Outputs.CreateAppResourcesInput> Resources { get; private set; } = null!;

        [Output("updatedAt")]
        public Output<Pulumi.Zeet.Time.Outputs.Time> UpdatedAt { get; private set; } = null!;

        [Output("userId")]
        public Output<string> UserId { get; private set; } = null!;


        /// <summary>
        /// Create a App resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public App(string name, AppArgs args, CustomResourceOptions? options = null)
            : base("zeet:resources:App", name, args ?? new AppArgs(), MakeResourceOptions(options, ""))
        {
        }

        private App(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("zeet:resources:App", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing App resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static App Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new App(name, id, options);
        }
    }

    public sealed class AppArgs : global::Pulumi.ResourceArgs
    {
        [Input("build")]
        public Input<Pulumi.Zeet.Model.Inputs.CreateAppBuildInputArgs>? Build { get; set; }

        [Input("deploy", required: true)]
        public Input<Pulumi.Zeet.Model.Inputs.CreateAppDeployInputArgs> Deploy { get; set; } = null!;

        [Input("docker")]
        public Input<Pulumi.Zeet.Model.Inputs.CreateAppDockerInputArgs>? Docker { get; set; }

        [Input("enabled", required: true)]
        public Input<bool> Enabled { get; set; } = null!;

        [Input("environmentId", required: true)]
        public Input<string> EnvironmentId { get; set; } = null!;

        [Input("environmentVariables")]
        private InputList<Pulumi.Zeet.Model.Inputs.CreateAppEnvironmentVariableInputArgs>? _environmentVariables;
        public InputList<Pulumi.Zeet.Model.Inputs.CreateAppEnvironmentVariableInputArgs> EnvironmentVariables
        {
            get => _environmentVariables ?? (_environmentVariables = new InputList<Pulumi.Zeet.Model.Inputs.CreateAppEnvironmentVariableInputArgs>());
            set => _environmentVariables = value;
        }

        [Input("github")]
        public Input<Pulumi.Zeet.Model.Inputs.CreateAppGithubInputArgs>? Github { get; set; }

        [Input("name", required: true)]
        public Input<string> Name { get; set; } = null!;

        [Input("projectId", required: true)]
        public Input<string> ProjectId { get; set; } = null!;

        [Input("resources", required: true)]
        public Input<Pulumi.Zeet.Model.Inputs.CreateAppResourcesInputArgs> Resources { get; set; } = null!;

        [Input("userId", required: true)]
        public Input<string> UserId { get; set; } = null!;

        public AppArgs()
        {
        }
        public static new AppArgs Empty => new AppArgs();
    }
}
