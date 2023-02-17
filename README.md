# pulumi-zeet

A pulumi provider for Zeet.


> :warning: **This project is in _alpha_**: This is an experimental project, and is not supported at this time!

## Overview

#### Supported Languages:
- Go 1.18+

### Getting Started

1. Install the pulumi provider plugin
   ```
   pulumi plugin install resource zeet ${VERSION} --server github://api.github.com/zeet-dev
   ```

1. Configure your pulumi stack
   ```
   pulumi config set zeet:endpoint ${ZEET_API_URL} 
   pulumi config set --secret zeet:api-token ${ZEET_API_TOKEN}
   ```

## Golang Stack
1. Install the SDK
   ```
   go get github.com/zeet-dev/pulumi-zeet
   ```


# Typescript Stack
1. Install the SDK
   ```
   npm i @zeet-dev/pulumi-zeet
   ```