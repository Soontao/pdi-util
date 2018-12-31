# SAP PDI Util Tools

[![CircleCI](https://circleci.com/gh/Soontao/pdi-util.svg?style=shield)](https://circleci.com/gh/Soontao/pdi-util)

Cli for SAP PDI, expose PDI apis to cli environment. For research.

## Latest Build

Just download latest binary files from [Github Release](https://github.com/Soontao/pdi-util/releases) Page

## Help

```bash
bash> pdi-util
NAME:
   PDI Util - A Command Line Tool for SAP Partner Development IDE

USAGE:
   pdi-util [global options] command [command options] [arguments...]

VERSION:
   SNAPSHOT

COMMANDS:
     check     static check
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

## Solution Operations

### list all solutions

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com solution list 
+----------+----------------+----------------+------------+----------+-------+
|   NAME   |   DESCRIPTION  |     STATUS     |  CUSTOMER  |  CONTACT | EMAIL |
+----------+----------------+----------------+------------+----------+-------+
| Removed  | Removed        | Deployed       |            | Removed  |       |
| Removed  | Removed        | In Development |            | Removed  |       |
+----------+----------------+----------------+------------+----------+-------+
```

## Source Operations

### download source from repo

Extremely fast, download with `35` goroutines defaultly.

(PDI download project files one by one in a single thread).

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com source download -s YQABCDEFG_ 
2018/11/24 21:20:52 Will download 1226 files to /Users/theosun/go/src/github.com/Soontao/pdi-util/output
1226 / 1226  100.00% 34s
2018/11/24 21:21:26 Done
```

## Static Check

Normally, just use `check all` sub command, it will do all available check & save result to excel file.

```bash
NAME:
   PDI Util check - static code check

USAGE:
   PDI Util check command [command options] [arguments...]

COMMANDS:
     all          do all check to file
     header       check copyright header
     backend      do backend check
     translation  do translation check
     name         check name convension

OPTIONS:
   --help, -h  show help
```

### Check All

example:

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com check all -s YQABCDEFG_
2018/12/31 09:52:59 Starting Backend Check...
2018/12/31 09:53:06 Backend Check Finished
2018/12/31 09:53:06 StartingTranslation Check...
2018/12/31 09:53:07 Translation Check Finished
2018/12/31 09:53:08 Copyright Header Check Finished
2018/12/31 09:53:08 Name Convention Check Finished
2018/12/31 09:53:08 In-Active File Check Finished
2018/12/31 09:53:08 Start Generating Excel File...
2018/12/31 09:53:08 Save Check Result File to check_all.xlsx
```

## [CHNAGELOG](./CHANGELOG.md)

## [LICENSE](./LICENSE)