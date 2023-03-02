## 安装Cobra-cli脚手架工具
```bash
$ go install github.com/spf13/cobra-cli@latest
```

## 在项目中下载Cobra依赖
```bash
$ go get -u github.com/spf13/cobra@latest
```

## 体验命令

```bash
# 常用方式
$ go run go-cobra/main.go hi --name zs --age 100 --like Coding,Running --address ShangHai
```

![image-20230302164212990](http://img.dingyuming.top/202303021642014.png)

```bash
# 未传参数时，显示帮助
$ go run go-cobra/main.go hi
```

![image-20230302163952456](http://img.dingyuming.top/202303021640809.png)