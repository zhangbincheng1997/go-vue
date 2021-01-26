# GO 配置

## vscode 插件

1. Go
2. CodeRunner

## Go 开发包

ctrl+shift+p > Go: Install/Update Tools

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

```Shell
gopkgs 安装失败：
> go get -v github.com/uudashr/gopkgs/cmd/gopkgs
```

## 环境变量

- GOROOT=D:\Go\
- path + %GOROOT%\bin
- GO111MODULE=on
- GOPROXY=https://goproxy.cn

## 报错处理

```Shell
go: cannot find main module; see 'go help modules'
> go mod init main
```
