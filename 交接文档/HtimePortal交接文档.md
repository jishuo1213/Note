

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

### 3	用户中心

#### 3.1 概述

用户中心是云上协同App
