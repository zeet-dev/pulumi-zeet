# pulumi-zeet-native

A pulumi provider for Zeet.


> :warning: **This project is in _alpha_**: This is an experimental project, and is not supported at this time!

## Overview

#### Supported Languages:
- Go 1.18+

### Getting Started

1. Install the pulumi provider plugin
   ```
    pulumi plugin install resource zeet-native [VERSION] --server github://api.github.com/zeet-dev
   ```
1. Install the SDK
   ```
   go get github.com/zeet-dev/pulumi-zeet-native
   ```
1. Configure your pulumi stack
   ```
   pulumi config set zeet-native:endpoint ${ZEET_API_ENDPOINT} 
   pulumi config set --secret zeet-native:api-token ${ZEET_API_TOKEN}
   ```
