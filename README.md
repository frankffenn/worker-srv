# worker-srv

## 安装 protoc

下载 protoc ,选择合适版本, [官网 Release 链接](https://github.com/protocolbuffers/protobuf/releases)        

**linux x86_64 版本**
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
    go get github.com/micro/micro/v2/cmd/protoc-gen-micro@master
```

重新编译      
```
    make
```

编译 go ，不编译 protobuf 文件      
```
    make build
```

## 启动服务 

启动 registry-srv  
```
./registry-srv --server_address=127.0.0.1:8080

```

启动 worker-srv
```
REGISTRY_ADDR=127.0.0.1:8080 ./worker-srv
```