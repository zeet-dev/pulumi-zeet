package main

import (
	"github.com/pulumi/pulumi-go-provider"
	"math/rand"
	"time"
)

// Random a simple illustrative example resource
// Each resources has a controlling struct.
// Resource behavior is determined by implementing methods on the controlling struct.
// The `Create` method is mandatory, but other methods are optional.
// - Check: Remap inputs before they are typed.
// - Diff: Change how instances of a resources are compared.
// - Update: Mutate a resources in place.
// - Read: Get the state of a resources from the backing config.
// - Delete: Custom logic when the resources is deleted.
// - Annotate: Describe fields and set defaults for a resources.
// - WireDependencies: Control how outputs and secrets flows through values.
type Random struct{}

// Each resources has in input struct, defining what arguments it accepts.
type RandomArgs struct {
	// Fields projected into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but its generally a
	// good idea.
	Length int `pulumi:"length"`
}

// Each resources has a state, describing the fields that exist on the created resources.
type RandomState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	RandomArgs
	// Here we define a required output called result.
	Result string `pulumi:"result"`
}

// All resources must implement Create at a minumum.
func (Random) Create(ctx provider.Context, name string, input RandomArgs, preview bool) (string, RandomState, error) {
	state := RandomState{RandomArgs: input}
	if preview {
		return name, state, nil
	}
	state.Result = makeRandom(input.Length)
	return name, state, nil
}

func makeRandom(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
