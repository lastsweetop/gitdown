# gitdown

gitlab的下载库
- 断点续传    下载大项目时，不需要因为各种情况的中断而重新下载
- 目录支持    支持单个目录下载，可以只下载你关心的目录，不需要整个工程都下载
- 最小下载    只下载某一版本某一目录的最新版本，不会去下载各种提交过程中的文件

# Install
## go的方式
```
go get -u  github.com/lastsweetop/gitdown
```
[GO install](https://github.com/golang/go)

## bin方式

在bin目录下有各种系统的可执行包，下载即可执行
```
├── darwin_386
│   └── gitdown
├── linux_386
│   └── gitdown
├── linux_amd64
│   └── gitdown
├── windows_386
│   └── gitdown.exe
└── windows_amd64
    └── gitdown.exe
```


# Getting Started

``` bash
gitdown -g https://gitlab.com/TeeFirefly/FireNow-Nougat.git -p firefly-rk3399 -d FFTools

-g  git的url，默认为https://gitlab.com/TeeFirefly/FireNow-Nougat.git
-p  分支，默认为master
-d  目录，默认为空
-t  线程数

```

# Command

```
Usage:
  gitdown [flags]

Flags:
  -d, --directory string   directories path
  -g, --git string         git url (default "https://gitlab.com/TeeFirefly/FireNow-Nougat/")
  -h, --help               help for gitdown
  -p, --patch string       patch (default "master")
  -t, --thread int         thread num (default 100)
```

# vertify

检查每个文件是否有文件损坏

```
vertify the repo

Usage:
  gitdown vertify [flags]

Flags:
  -g, --git string     git url (default "https://gitlab.com/TeeFirefly/FireNow-Nougat")
  -h, --help           help for vertify
  -p, --patch string   patch (default "firefly-rk3399")
```

# License
This repository is Copyright (c) 2017 lastsweetop, Inc. All rights reserved. It is licensed under the MIT license. Please see the LICENSE file for applicable license terms.

# Authors
The primary author is [lastsweetop](http://www.lastsweetop.com), with some documentation and other minor contributions by others at lastsweetop.
