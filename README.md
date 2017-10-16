# gitdown

支持断点续传的gitlab库下载工具


# Install

```
go get -u  github.com/lastsweetop/gitdown
```


# Getting Started

``` bash
gitdown -g https://gitlab.com/TeeFirefly/FireNow-Nougat.git -p firefly-rk3399 -d FFTools

-g  git的url，默认为https://gitlab.com/TeeFirefly/FireNow-Nougat.git
-p  分支，默认为master
-d  目录，默认为空

```

# Command

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

# License
This repository is Copyright (c) 2017 lastsweetop, Inc. All rights reserved. It is licensed under the MIT license. Please see the LICENSE file for applicable license terms.

# Authors
The primary author is [lastsweetop](http://www.lastsweetop.com), with some documentation and other minor contributions by others at lastsweetop.
