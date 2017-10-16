#!/bin/bash

set -e

usage()
{
cat << EOF
usage:
    $(basename $0) [-u|k|a] [-d dts_file_name] [-l lunch] [-j make_thread]
    -u|k|a: make uboot|kernel|android alone, if this arg is not exist, make all images default
    -d: kernel dts name
    -l: lunch name when make android
    -j: make theard num, if have not this arg, default theard is 1

NOTE: Run in the path of SDKROOT
EOF

if [ ! -z $1 ] ; then
    exit $1
fi
}

MAKE_THEARD=1
KERNEL_DTS='rk3399-firefly'
USER_LUNCH='rk3399_firefly_box-userdebug'
MAKE_MODULES=''
MAKE_ALL=true

while getopts "ukahj:d:l:" arg
do
	case $arg in
		 u|k|a)
            MAKE_MODULES=$arg
            MAKE_ALL=false
			;;
		 j)
			MAKE_THEARD=$OPTARG
			;;
		 d)
			KERNEL_DTS=$OPTARG
			;;
		 l)
			USER_LUNCH=$OPTARG
			;;
		 h)
			usage 0
			;;
		 ?) 
			usage 1
			;;
	esac
done

FFTOOLS_PATH=$(dirname $0)

if $MAKE_ALL || [ $MAKE_MODULES = 'u' ]; then
    pushd u-boot/
    make rk3399_box_defconfig 
    make ARCHV=aarch64 -j $MAKE_THEARD
    popd
fi

if  $MAKE_ALL || [ $MAKE_MODULES = 'k' ]; then
    pushd kernel/
    make ARCH=arm64 firefly_defconfig
    make  ARCH=arm64 "${KERNEL_DTS}.img" -j $MAKE_THEARD
    popd
fi

if $MAKE_ALL || [ $MAKE_MODULES = 'a' ]; then
    source ${FFTOOLS_PATH}/build.sh
    lunch "$USER_LUNCH"
    make installclean
    make -j $MAKE_THEARD

    ./mkimage.sh
fi

echo "Firefly-RK3399 make images finish!"
