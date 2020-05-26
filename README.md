# SAP PDI Util Tools

[![CircleCI](https://circleci.com/gh/Soontao/pdi-util.svg?style=shield)](https://circleci.com/gh/Soontao/pdi-util)
[![Snapshot Build](https://github.com/Soontao/pdi-util/workflows/Snapshot%20Build/badge.svg)](https://github.com/Soontao/pdi-util/actions?query=workflow%3A%22Snapshot+Build%22)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/Soontao/pdi-util.svg)](https://github.com/Soontao/pdi-util/releases)
[![](https://godoc.org/github.com/Soontao/pdi-util?status.svg)](http://godoc.org/github.com/Soontao/pdi-util)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Soontao_pdi-util&metric=alert_status)](https://sonarcloud.io/dashboard?id=Soontao_pdi-util)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/theosun/pdiutil)](https://hub.docker.com/repository/docker/theosun/pdiutil)
[![Docker Image Size (tag)](https://img.shields.io/docker/image-size/theosun/pdiutil/latest)](https://hub.docker.com/repository/docker/theosun/pdiutil)

Cli for SAP PDI, **an UN-OFFCIAL command line tool.**

## Latest Build

Just download latest binary files from the [Github Release](https://github.com/Soontao/pdi-util/releases) page

## Features

* [x] [list all solutions](https://github.com/Soontao/pdi-util/wiki/How-to-list-all-solutions)
* [x] [static check & export results to excel](https://github.com/Soontao/pdi-util/wiki/How-to-do-static-check)
* [x] download current sources in solution
* [x] download assembled package (in history)
* [x] view single file history (PDI from release 2002 has support this feature)
* [x] [activate, assemble and download assembled package](https://github.com/Soontao/pdi-util/wiki/How-to-assemble-solution)
* [x] [deploy solution](https://github.com/Soontao/pdi-util/wiki/How-to-deploy-solution)
* [x] static text spell check
* [x] CI/CD [Jenkinsfile](./jenkins) pipeline provided
* [x] [solution source statistics](https://github.com/Soontao/pdi-util/wiki/How-to-statistics-solution-scale)
* [x] [code ast parser](https://github.com/Soontao/grammar-pdi)

## Help

run with `--help` to show command help

Almost all options can be configured in the system environment variables.

## [CHANGELOG](./CHANGELOG.md)

## [LICENSE](./LICENSE)
