package config

import (
	"fmt"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/gql"
)

var ZeetClient gql.ZeetClient

func MakeConfigureFunc(version string) func(ctx p.Context, request p.ConfigureRequest) error {
	println("configure func", version)
	return func(ctx p.Context, request p.ConfigureRequest) error {
		endpoint, err := getConfigValue(request, "endpoint")
		if err != nil {
			return err
		}
		apiToken, err := getConfigValue(request, "api-token")
		if err != nil {
			return err
		}
		ZeetClient = gql.NewZeetGraphqlClient(endpoint, apiToken, version)
		return nil
	}
}

func getConfigValue(request p.ConfigureRequest, key string) (string, error) {
	fullKey := configKey(key)
	value, ok := request.Variables[fullKey]
	if !ok {
		return "", fmt.Errorf("missing config value: %s", fullKey)
	}
	return value, nil
}

func configKey(key string) string {
	return fmt.Sprintf("zeet-native:config:%s", key)
}
