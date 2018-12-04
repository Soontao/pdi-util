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

Extremely fast, download with `35` goroutines defaultly.

(PDI download project files one by one in a single thread).

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com source download -s YQABCDEFG_ 
2018/11/24 21:20:52 Will download 1226 files to /Users/theosun/go/src/github.com/Soontao/pdi-util/output
1226 / 1226  100.00% 34s
2018/11/24 21:21:26 Done
```

## copyright header check

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com check header -s YQABCDEFG_
2018/11/25 16:13:08 Will check 532 ABSL/BO/XBO Defination
 532 / 532  100.00% 5s
2018/11/25 16:13:14 Not found copyright header in: /API/XXXXXXX.absl
...
...
2018/11/25 16:13:14 Totally 247 files (of 532) lost copyright header
```

## name convension check

check name convension of source code filename

```bash
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com check name -s YQABCDEFG_
2018/12/03 20:13:50 Name Convension BPM\CSD_BPMSystem.csd: filename should be CS_BPMSystem.csd
2018/12/03 20:13:50 Name Convension HCM\EWSI_CH_USER_ID.csd: filename should be CS_CH_USER_ID.csd
2018/12/03 20:13:50 Name Convension _Common\DT_MDRInputData.bo: filename should be BO_MDRInputData.bo
2018/12/03 20:13:50 Name Convension _EXT\EBO_ServiceRequest.xbo: filename should be BOE_ServiceRequest.xbo
2018/12/03 20:13:50 finished
```

## backend check

execute runtime check on backend

support follow files now

```json
{
  	".absl": "ABSL",
	".bo":   "BUSINESS_OBJECT",
	".qry":  "QUERYDEF",
	".xbo":  "EXTENSION_ENTITY",
	".bco":  "BCO",
	".bcc":  "BCSET"
}
```

extremely fast !


```
bash> pdi-util -u USER -p PASS -t myxxxxx.c4c.saphybriscloud.com check backend -s YQABCDEFG_
 133 / 133  100.00% 1s
2018/12/04 22:27:08 [W] CustomBO.bo(8 ,26 ): Use of data type 'Description' is not supported in queries
2018/12/04 22:27:08 [W] CustomBO.bo(8 ,26 ): Do not store external document data in unrestricted data type 'Description'. Recommendation is to use Attachment Folder (refer SDK Help Documentation Section 7.2.2.12). Please refer blog "Text Types Usage" in Community Forum for more on text data types.
...
...
2018/12/04 22:28:32 Finished
```

## [CHNAGELOG](./CHANGELOG.md)

## [LICENSE](./LICENSE)