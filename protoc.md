
# protobuf 命令备注
常用命令笔记

1.生成标准的 protobuf go 代码
```shell
protoc --go_out=. example.proto
```
2.生成更高效 protobuf go代码
```shell
protoc --gofast_out=. example.proto
```
3.生成更高效 protobuf go代码
```shell
protoc --go-grpc_out=. example.proto
```
4.生成更高效 protobuf go代码
```shell
protoc --go_out=. --go-grpc_out=. example.proto
```
5.生成更高效 protobuf go代码
```shell
protoc --gofast_out=. --go-grpc_out=. example.proto
```