# PDI Solution CI/CD Repository

This is the PDI Solution CI/CD script host repository.

## Setup Variables

Setup following variables in the `jenkinsfile` for build and deployment

```groovy
def cliVersion = "v1.9.10"
def utilDownloadURI = ""
def cliName = "pdiutil-${cliVersion}-darwin-amd64"
def utilZipFileName = "pdiutil.zip"
// please use the readable 'Solution Description' instead of the technical 'Solution ID'
// so that pdi-util tool can auto detect the correct solution id after patch solution created
def solution = ""
// source tenant
def devTenant = ""
// target tenant
def targetTenant = ""
// Jenkins user/password credential ID
def devUserCredentialId = ''
// Jenkins user/password credential ID
def targetUserCredentialId = ''
// specific the current ByD/C4C release version
def releaseVersion = ""
```

## Build Only

If you only want to build/assemble, just comment the `Deploy` stage.

## Jenkins Guide

Please use this `BuildDeploy.Jenkinsfile` pipeline after the `first-time deployment`, because the `first-time deployment` is not tested by the [pdi-util](https://github.com/Soontao/pdi-util) project by now

Please setup this jenkinsfile as normal `Pipeline` (instead of `Multi Branch Pipeline`), so that you could control build/deployment flexible.

