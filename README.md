# SAP PDI Util Tools

[![CircleCI](https://circleci.com/gh/Soontao/pdi-util.svg?style=shield)](https://circleci.com/gh/Soontao/pdi-util)

Cli for SAP PDI

## Latest Build

Just download latest binary files from [Github Release](https://github.com/Soontao/pdi-util/releases) Page

## Help

```bash
bash> pdi-util 
NAME:
   PDI Util - A cli util for SAP PDI

USAGE:
   pdi-util [global options] command [command options] [arguments...]

VERSION:
   SNAPSHOT

COMMANDS:
     check     code static check
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

Almost all options can be configured in the system environment variables.

## list all solutions

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com solution list 
+----------+----------------+----------------+------------+----------+-------+
|   NAME   |   DESCRIPTION  |     STATUS     |  CUSTOMER  |  CONTACT | EMAIL |
+----------+----------------+----------------+------------+----------+-------+
| Removed  | Removed        | Deployed       |            | Removed  |       |
| Removed  | Removed        | In Development |            | Removed  |       |
+----------+----------------+----------------+------------+----------+-------+
```

## download source from repo

```bash
bash> pdi-util source download -h
NAME:
   PDI Util source download - download all files in a solution

USAGE:
   PDI Util source download [command options] [arguments...]

OPTIONS:
   --solution value, -s value    The PDI Solution Name [$SOLUTION_NAME]
   --output value, -o value      Output directory (default: "output") [$OUTPUT]
   --concurrent value, -c value  concurrent goroutine number (default: 35) [$DOWNLOAD_CONCURRENT]
   
```

Extremely fast, download with `35` goroutines defaultly. 

(PDI download project files one by one in a single thread).

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com source download -s YQABCDEFG_ 
2018/11/24 21:20:52 Will download 1226 files to /Users/theosun/go/src/github.com/Soontao/pdi-util/output
1226 / 1226  100.00% 34s
2018/11/24 21:21:26 Done
```

## check copy right header

```bash
bash> pdi-util check header -h
NAME:
   PDI Util check header - check copyright header

USAGE:
   
make sure all absl & bo have copyright header with following format:

/*
  Function: make sure all absl & bo have copyright header
  Author: Theo Sun
  Copyright: ?
*/

OPTIONS:
   --solution value, -s value    The PDI Solution Name [$SOLUTION_NAME]
   --concurrent value, -c value  concurrent goroutine number when download from remote (default: 35) [$DOWNLOAD_CONCURRENT]
```

example

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com check header -s YQABCDEFG_
2018/11/25 16:13:08 Will check 532 ABSL/BO/XBO Defination
 532 / 532  100.00% 5s
2018/11/25 16:13:14 Not found copyright header in: /API/XXXXXXX.absl
... 
...
2018/11/25 16:16:34 Totally 247 file (of 532) lost copyright header
```

## [CHNAGELOG](./CHANGELOG.md)

## [LICENSE](./LICENSE)