package config

import (
	"fmt"
	"github.com/pulumi/pulumi-go-provider"
	"github.com/zeet-dev/pulumi-zeet-native/provider/pkg/gql"
)

var ZeetClient gql.ZeetGraphqlClient

func ConfigureFunc(ctx provider.Context, request provider.ConfigureRequest) error {
	endpoint, err := getConfigValue(request, "endpoint")
	if err != nil {
		return err
	}
	apiToken, err := getConfigValue(request, "api-token")
	if err != nil {
		return err
	}
	ZeetClient = gql.NewZeetGraphqlClient(endpoint, apiToken)
	return nil
}

func configKey(key string) string {
	return fmt.Sprintf("zeet-native:config:%s", key)
}

func getConfigValue(request provider.ConfigureRequest, key string) (string, error) {
	fullKey := configKey(key)
	value, ok := request.Variables[fullKey]
	if !ok {
		return "", fmt.Errorf("missing config value: %s", fullKey)
	}
	return value, nil
}
