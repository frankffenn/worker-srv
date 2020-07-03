# worker-srv

## Install

下载 protoc ,选择合适版本, [官网 Release 链接](https://github.com/protocolbuffers/protobuf/releases)        

```shell
 cd ~
 wget https://github.com/protocolbuffers/protobuf/releases/download/v3.12.3/protoc-3.12.3-linux-x86_64.zip
 unzip protoc-3.12.3-linux-x86_64.zip
 cp protoc-3.12.3-linux-x86_64/protoc /usr/local/bin

```


安装 protoc-gen-go      
```shell
    go get -u github.com/golang/protobuf/protoc-gen-go
```

安装 micro 插件     
```shell
    go get github.com/micro/protobuf/proto github.com/micro/protobuf/protoc-gen-go
    ##
```

编译      
```
    make all
```