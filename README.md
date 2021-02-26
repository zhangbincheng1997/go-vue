# go-vue

## 环境变量

- GOROOT=D:\Go\
- path + %GOROOT%\bin
- GO111MODULE=on
- GOPROXY=https://goproxy.cn

## vscode 配置

1. Go
2. CodeRunner
3. ctrl+shift+p > Go: Install/Update Tools

```Comment
1. gocode (Auto-completion, does not work with modules)
2. gopkgs (Auto-completion of unimported packages && Add Import feature)
3. go-outline (Go to symbol in file)
4. go-symbols (Go to symbol in workspace)
5. guru (Find all references and Go to implementation of symbols)
6. gorename (Rename symbols)
7. gotests (Generate unit tests)
8. gomodifytags (Modify tags on structs)
9. impl (Stubs for interfaces)
10. fillstruct (Fill structs with defaults)
11. goplay (The Go playground)
12. godocter (Extract to functions and variables)
13. dlv (Debugging)
14. gocode-gomod (Auto-completion, works with modules)
15. godef (Go to definition)
16. goreturns (Formatter)
17. golint (Linter)
```

## Web框架

```Shell
# go get -u 强制使用网络更新

# 热更新
$ go get -u github.com/cosmtrek/air

# Web框架
$ go get -u github.com/gin-gonic/gin
$ go get -u github.com/gin-contrib/cors
$ go get -u github.com/appleboy/gin-jwt/v2
# 核心
$ go get -u github.com/spf13/viper # Config
$ go get -u go.uber.org/zap # Log
$ go get -u gopkg.in/natefinch/lumberjack.v2
$ go get -u gorm.io/gorm # ORM
$ go get -u gorm.io/driver/mysql # MySQL
$ go get -u github.com/go-redis/redis/v8 # Redis
$ go get -u go.mongodb.org/mongo-driver/mongo # MongoDB
$ go get -u github.com/qiniu/go-sdk/v7 # OSS
# 文档
$ go get -u github.com/swaggo/swag/cmd/swag
$ go get -u github.com/swaggo/gin-swagger
$ go get -u github.com/swaggo/files
# 其他
$ go get -u github.com/jinzhu/copier
# ...
$ go get -u github.com/fsnotify/fsnotify
$ go get -u github.com/alecthomas/template
```
