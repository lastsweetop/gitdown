# gitdown

支持断点续传的gitlab库下载工具

# 示例

``` bash
gitdown -g https://gitlab.com/TeeFirefly/FireNow-Nougat.git -p firefly-rk3399 -d FFTools

-g  git的url，默认为https://gitlab.com/TeeFirefly/FireNow-Nougat.git 
-p  分支，默认为master
-d  目录，默认为空

```

# 命令行

```
gitdown -h
Usage:
  gitdown [flags]

Flags:
  -d, --directory string   directories path
  -g, --git string         git url (default "https://gitlab.com/TeeFirefly/FireNow-Nougat/")
  -h, --help               help for gitdown
  -p, --patch string       patch (default "master")
```
