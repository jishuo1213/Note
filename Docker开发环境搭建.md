## Docker开发环境搭建

### 环境

* OS:Arch Linux

### 依赖

* git
* make
* docker 版本高点比较好，我用pacman装的最新版 20.10.11

### 步骤

####  安装依赖

```shell
sudo pacman -Sy docker docker-compose git make
```

#### 克隆源码

找一个空文件夹

```shell
git clone https://github.com/moby/moby.git
```

docker的源码现在叫moby

#### 编译开发镜像

```shell
cd moby
```

需要在DockerFile里面加入

```yaml
ENV GOPROXY="https://goproxy.io,direct"
ENV HTTP_PROXY="http://10.9.11.26:18888"
ENV HTTPS_PROXY="http://10.9.11.26:18888"
```

然后执行

```shell
make BIND_DIR=. shell
```

等待完成

