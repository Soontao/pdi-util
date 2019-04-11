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
   cli [global options] command [command options] [arguments...]

VERSION:
   SNAPSHOT

COMMANDS:
     check     static code check
     package   package related commands
     solution  solution related operations
     source    source code related operations
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value, -u value  The PDI Development User [$PDI_USER]
   --password value, -p value  The PDI Development User Password [$PDI_PASSWORD]
   --hostname value, -t value  The PDI Tenant host [$PDI_TENANT_HOST]
   --release value, -r value   The tenant release version, e.g. 1902 (default: "1902") [$TENANT_RELEASE]
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

### source version access

access different source version in history.

```bash
# try to use BO name search version
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com source version -s Test_Solution -f BO_Notification
2019/03/27 10:31:24 More than one files matched name: BO_Notification
2019/03/27 10:31:24 /YQCS6QLDY_MAIN/SRC/Theo/BO_Notification-Root-Event-BeforeSave.absl
2019/03/27 10:31:24 /YQCS6QLDY_MAIN/SRC/Theo/BO_Notification-Root.node
2019/03/27 10:31:24 /YQCS6QLDY_MAIN/SRC/Theo/BO_Notification.bo

# use the bo file name search version
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com source version -s Test_Solution -f BO_Notification-Root-Event-BeforeSave.absl
+-------------------------------+--------+--------+------------------------+
|           DATE TIME           | ACTION |  USER  |       VERSIONID        |
+-------------------------------+--------+--------+------------------------+
| 2018-12-20 08:26:31 +0000 UTC | ADD    | YS1004 | 20181220082631.8250240 |
| 2018-12-20 08:29:47 +0000 UTC | EDIT   | YS1004 | 20181220082947.8340020 |
| 2018-12-20 08:33:53 +0000 UTC | EDIT   | YS1004 | 20181220083353.8092850 |
| 2018-12-20 08:35:28 +0000 UTC | EDIT   | YS1004 | 20181220083528.7241340 |
| 2018-12-20 08:50:22 +0000 UTC | EDIT   | YS1004 | 20181220085022.6248990 |
| 2018-12-20 09:18:55 +0000 UTC | EDIT   | YS1004 | 20181220091855.8923780 |
+-------------------------------+--------+--------+------------------------+

# use version id access source version content
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com source version -s Test_Solution -f BO_Notification-Root-Event-BeforeSave.absl -v 20181220082631.8250240

# content ignore

# diff two version change
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com source version -s Test_Solution -f BO_Notification-Root-Event-BeforeSave.absl --from 20181220085022.6248990 --to 20181220091855.8923780
# result have color, markdown can not present that so ignore

```

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

## [CHNAGELOG](./CHANGELOG.md)

## [LICENSE](./LICENSE)