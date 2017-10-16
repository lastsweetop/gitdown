#!/bin/bash

set -e

. build/envsetup.sh >/dev/null && setpaths

ANDROID_ROOT="$(get_abs_build_var)"
TARGET_PRODUCT=`get_build_var TARGET_PRODUCT`
IMG_ROOT="$ANDROID_ROOT/rockdev/Image-$TARGET_PRODUCT"
IMG_LIST="parameter.txt MiniLoaderAll.bin trust.img uboot.img resource.img kernel.img boot.img recovery.img misc.img system.img"
SRC_RAR_FILE="$(dirname $0)/AndroidTool.rar"
DST_RAR_PATH="$ANDROID_ROOT/rockdev/Image-$TARGET_PRODUCT/"

if [ -d $IMG_ROOT ];then
    if [ ! -e $SRC_RAR_FILE ];then
        echo "Make sure you have file \"$SRC_RAR_FILE\"!"
        exit 2
    fi        
    for img in $IMG_LIST
    do
        if [ ! -e $IMG_ROOT/$img ];then
            echo "Make sure you have file \"$img\"!"
            exit 3
        fi
    done
else
    echo "Make sure you have directory \"$IMG_ROOT\"!"
    exit 1
fi

which rar > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "Make sure you have tool \"rar\""
    exit 4
fi

if [ "z${1}" != "z" ] ; then
	DST_RAR_NAME="${1}.rar"
else
	DST_RAR_NAME="Firefly-RK3399_Android$(get_build_var PLATFORM_VERSION)_$(date -d today +%y%m%d).rar"
fi
DST_RAR_NAME=$(echo $DST_RAR_NAME | sed s/[[:space:]]//g)

cp "$SRC_RAR_FILE" "$IMG_ROOT/$DST_RAR_NAME"
cd $IMG_ROOT

# put all the *img and update log into rockdev/Image/
rar a -ap"rockdev/Image/" $DST_RAR_NAME $IMG_LIST

echo -e "\nMaking AndroidTools.rar:\n$(readlink -f ${IMG_ROOT}/${DST_RAR_NAME})\n"
