

## 交接文档

### 1	总体概述

此文档是云上协同部分go工程交接文档，包括用户中心、开放平台、新消息引擎等服务

### 2	开发环境配置

所有的工程都使用gomod进行依赖管理，在进行开发之前，需要搭建go开发环境、git环境然后进行gomod的配置。

#### 2.1 go开发环境搭建

具体请参照网上文档，进行环境变量的配置

#### 2.2 gomod配置

首先联系gitlab管理员加上相关开发人员工程源码权限，之后前往gitlab个人中心，生成accesstoken，如下图：

![](./resource/2022-02-14_17-53.png)

生成之后，保存形如

```shell
vJtDFca4wNz8pAyM8KMG
```

的access token。

之后使用如下命令配置gomod依赖git私服替换地址:

```shell
git config --global url."http://{{git_user_name}}:{{access_token}}@git.icity.cn:{{port}}".insteadOf "https://git.icity.cn"
```

然后配置gomod环境变量:

```shell
go env -w GOPROXY=https://goproxy.io,direct
go env -w GOPRIVATE=git.icity.cn
```

最后 在host文件添加映射:

```
{{gitlab_ip}} git.icity.cn
```

其中:

```
git_user_name:用户的git用户名
access_token:上一步生成得access_token
port:gitlab私服的端口，如果是80端口，可以不写
gitlab_ip:私服gitlab ip地址
```

完成配置后，可以从gitlab clone工程源码，在go.mod文件所在目录，使用

```
go mod tidy
```

下载依赖，每次go.mod文件出现改动之后，需要重新运行此命令。

### 3    基础框架

#### 3.1 概述

所有go工程都有一些公用的依赖，包含常用工具方法、加解密、web框架、服务注册、服务发现、配置管理等。主要有以下四个工程:

|   工程名    | git地址                                        | 简介                                                         |
| :---------: | :--------------------------------------------- | :----------------------------------------------------------- |
|   Logger    | http://10.110.9.195:8888/ccwork/go/Logger.git  | 日志库，包含基本的日志分级、分文件等功能                     |
| DBOperation | http://10.110.9.195:8888/ccwork/go/DBOperation | 数据库操作库，包含mysql、es、redis、sqlite、etcd、hbase、国产数据库等基础操作封装 |
|    Utils    | http://10.110.9.195:8888/ccwork/go/utils.git   | 包含常用加解密、http请求、服务注册发现、动态配置管理、web框架封装等功能 |
|    proto    | http://10.110.9.195:8888/ccwork/go/proto.git   | 各个服务的grpc描述文件                                       |

#### 3.2 Logger

logger是基于 [lumberjack Logger]: https://github.com/natefinch/lumberjack 日志记录封装库，一般用法如下：

```go
logger.InitLogConfig(logger.DEBUG, isDebug)//初始化配置，包含log等级和是否debug模式，debug模式会把所有日志打到标准输出,注意，必须进行初始化，不初始化直接使用会panic
logger.Debugf("upload path not existed:%v", err)
logger.Errorf("upload path not existed:%v", err)
logger.Fatalf("upload path not existed:%v", err)
```

非debug模式，会把不同等级日志打到服务包所在目录logfile文件夹下不同文件里面，具体的日志文件保存和删除策略，可以参考源码。

一些需要改进的点:

* 同一目录不同服务日志分开
* 日志格式美化

#### 3.3 DBOperation

DBOperation是数据库操作封装库，包含一些常用的诸如mysql、mongo、redis等操作。

目录结构如下:

```
├── es
│   ├── es.go
│   └── es_test.go
├── etcd
│   ├── etcd.go
│   └── etcd_test.go
├── hbase
│   ├── hbase.go
│   └── hbase_test.go
├── ldap
│   ├── ldap.go
│   └── ldap_test.go
├── main.go
├── mongo
│   └── mgo_operation.go
├── mysql---------------------------------------------------------------------------es
│   ├── dmsql_drive.go
│   ├── mysql_drive.go
│   ├── mysql.go
│   ├── mysql_test.go
│   └── oracle_drive.go
├── nsq
│   ├── consumer.go
│   ├── consumer_test.go
│   ├── producer.go
│   └── producer_test.go
├── orm
│   ├── orm.go
│   └── orm_test.go
├── redis
│   ├── redis.go
│   └── redis_test.go
├── sqlite
│   ├── orm.go
│   ├── sqlite.go
│   └── sqlite_test.go
└── zk
    ├── lock.go
    └── operation.go
```

每个包代表一类数据库操作封装，具体使用方式可以参照test文件或者具体项目写法，其中， mysql包中包含达蒙和oracle

### 3	用户中心

#### 3.1 概述

用户中心是云上协同App和管理后台负责登录、认证、用户管理等功能的一系列服务集合。git地址为:

```shel
http://10.110.9.195:8888/ccwork/go/HtimePortal.git
```

目录结构如下:

```
.
├── config---------------------------------------------------------------------------配置读取
├── eventsync------------------------------------------------------------------------事件同步服务
├── htimefile------------------------------------------------------------------------文件上传代理
├── htimesearch----------------------------------------------------------------------搜索服务
├── kafka
├── lock
├── orgmanager-----------------------------------------------------------------------互联网组织绑定服务
├── patch
├── proof
├── rbac-----------------------------------------------------------------------------角色服务
├── router
├── rpc
├── server
├── service
│   ├── common
│   ├── email
│   ├── htime
│   ├── inner
│   ├── logic
│   ├── portal
│   ├── repo
│   └── sms
├── servicenumber---------------------------------------------------------------------服务号管理服务
├── sync
├── thirdauth-------------------------------------------------------------------------第三方认证服务
├── uploads
└── version
```

