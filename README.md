# SAP PDI Util Tools

[![CircleCI](https://circleci.com/gh/Soontao/pdi-util.svg?style=shield)](https://circleci.com/gh/Soontao/pdi-util)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/Soontao/pdi-util.svg)](https://github.com/Soontao/pdi-util/releases)
[![](https://godoc.org/github.com/Soontao/pdi-util?status.svg)](http://godoc.org/github.com/Soontao/pdi-util)

[![Docker Automated build](https://img.shields.io/docker/cloud/automated/theosun/pdiutil.svg)](https://cloud.docker.com/repository/docker/theosun/pdiutil)
[![Size](https://shields.beevelop.com/docker/image/image-size/theosun/pdiutil/latest.svg?style=flat-square)](https://cloud.docker.com/repository/docker/theosun/pdiutil)

Cli for SAP PDI.

## Latest Build

Just download latest binary files from the [Github Release](https://github.com/Soontao/pdi-util/releases) page

## Features

* [x] [list all solutions](https://github.com/Soontao/pdi-util/wiki/How-to-list-all-solutions)
* [x] [static check & export results to excel](https://github.com/Soontao/pdi-util/wiki/How-to-do-static-check)
* [x] download current sources in solution
* [x] download assembled package (in history)
* [x] view single file history
* [x] [activate, assemble and download assembled package](https://github.com/Soontao/pdi-util/wiki/How-to-assemble-solution)
* [x] [deploy solution](https://github.com/Soontao/pdi-util/wiki/How-to-deploy-solution)
* [x] static text spell check
* [x] CI/CD [Jenkinsfile](./jenkins) provided
* [x] [solution source statistics](https://github.com/Soontao/pdi-util/wiki/How-to-statistics-solution-scale)
* [ ] [code ast parser](https://github.com/Soontao/grammar-pdi) (in progress)

## Help

run with `--help` to show command help

Almost all options can be configured in the system environment variables.

## To Do

* [ ] Documents

## [CHANGELOG](./CHANGELOG.md)

## [LICENSE](./LICENSE)
