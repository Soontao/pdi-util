# pdi-util

[![CircleCI](https://circleci.com/gh/Soontao/pdi-util.svg?style=shield)](https://circleci.com/gh/Soontao/pdi-util)

Cli for SAP PDI

## Help

```bash
bash> pdi-util 
NAME:
   PDI Util - A cli util for SAP PDI

USAGE:
   pdi-util [global options] command [command options] [arguments...]

VERSION:
   v1-alpha

COMMANDS:
     session   session related operations
     solution  solution related operations
     source    source code related operations
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value, -u value  The PDI Development User [$PDI_USER]
   --password value, -p value  The PDI Development User Password [$PDI_PASSWORD]
   --hostname value, -t value  The PDI Tenant host [$PDI_TENANT_HOST]
   --help, -h                  show help
   --version, -v               print the version
```

## list all solutions

```bash
bash> pdi-util -u DEV_USER -p DEV_PASSWORD -t myxxxxx.c4c.saphybriscloud.com solution list 
+----------+----------------+----------------+------------+----------+-------+
|    ID    |   DESCRIPTION  |     STATUS     |  CUSTOMER  |  CONTACT | EMAIL |
+----------+----------------+----------------+------------+----------+-------+
| Removed  | Removed        | Deployed       |            | Removed  |       |
| Removed  | Removed        | In Development |            | Removed  |       |
+----------+----------------+----------------+------------+----------+-------+
```

## download source from repo

```bash
bash> pdi-util -u DEV_USER -p DEV_PASSWORD -t myxxxxx.c4c.saphybriscloud.com source download -s SOLUTION_NAME 
2018/11/24 12:59:36 Will download 1226 files to /Users/theosun/go/src/github.com/Soontao/pdi-util/output
2018/11/24 13:00:05 Done
```