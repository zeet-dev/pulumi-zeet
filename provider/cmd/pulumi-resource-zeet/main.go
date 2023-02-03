// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/config"
	"github.com/zeet-dev/pulumi-zeet/provider/pkg/resources"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

func main() {
	// We tell the config what resources it needs to support.
	// NOTE: resources must be implemented in this main package! See README
	providerOptions := infer.Options{
		Resources: []infer.InferredResource{
			infer.Resource[resources.Project, resources.ProjectArgs, resources.ProjectState](),
			infer.Resource[resources.Environment, resources.EnvironmentArgs, resources.EnvironmentState](),
			infer.Resource[resources.App, resources.AppArgs, resources.AppState](),
		},
		Metadata: schema.Metadata{
			LanguageMap: map[string]any{
				"go": map[string]any{
					"importBasePath": "github.com/zeet-dev/pulumi-zeet/sdk/go/zeet",
				},
			},
		},
	}
	provider := infer.Provider(providerOptions)
	provider.Configure = config.MakeConfigureFunc(Version)
	p.RunProvider("zeet", Version, provider)
}
