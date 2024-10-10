
# 测试命令备注
常用测试命令笔记

## 单元测试(TestXXX)
1、执行包及其子包下所有单元测试用例
```shell
go test -v
```
2、执行包下特定单元测试用例
```shell
go test -v -run="xxx"
```

## 基准/压力测试(BenchmarkXXX)
1、执行包及其子包下所有基准测试用例
```shell
go test -bench=.
```
2、执行包下特定基准测试用例
```shell
go test -bench="xxx"
```
3、执行包下基准测试用例并显示内存信息(-benchmem)
```shell
go test -bench=. -benchmem
go test -bench="xxx" -benchmem
```
4、执行包下基准测试用例并指定执行时间(-benchtime=xs)
```shell
go test -bench=. -benchtime=5s
go test -bench="xxx" -benchtime=5s
```

## 参考文档
https://blog.csdn.net/qq_52678569/article/details/141786689