#!/bin/bash

set -e

usage()
{
cat << EOF
usage:
    $(basename $0) [-n update_img_name] [-l lunch] 
    -l: lunch name when make android
    -n: dest update.img name, if have not this arg, there are an default name, like:
        Firefly-RK3399_Android6.0_YYMMDD

NOTE: Run in the path of SDKROOT
EOF

if [ ! -z $1 ] ; then
    exit $1
fi
}

while getopts "hn:l:" arg
do
	case $arg in
		 l)
			USER_LUNCH=$OPTARG
			;;
		 n)
			UPDATE_USER_NAME=$OPTARG
			;;
		 h)
			usage 0
			;;
		 ?) 
			usage 1
			;;
	esac
done

. build/envsetup.sh >/dev/null && setpaths
if [ ! -z "$USER_LUNCH" ] ; then
    lunch "$USER_LUNCH"
fi

#set -x
PRODUCT_FIREFLY_NAME=`get_build_var PRODUCT_FIREFLY_NAME`
PRODUCT_FIREFLY_NAME=${PRODUCT_FIREFLY_NAME:="DEFAULT"}
echo "PRODUCT_FIREFLY_NAME=$PRODUCT_FIREFLY_NAME"

TARGET_PRODUCT=`get_build_var TARGET_PRODUCT`
echo -e "TARGET_PRODUCT=$TARGET_PRODUCT\n"

TARGET_VERSION=`get_build_var PLATFORM_VERSION`

CMD_PATH=$(dirname $0)
IMAGE_SRC_CUR_PATH=rockdev/Image-$TARGET_PRODUCT
IMAGE_SRC_PATH=$(readlink -f $IMAGE_SRC_CUR_PATH)
LOADER_PATH=$(awk /bootloader/'{print $2}' ${CMD_PATH}/package-file | tr -d '\r\n')
LINK_IMAGE_PATH=$(echo ${LOADER_PATH} | awk -F/ '{print $1}')

if [ -z "$UPDATE_USER_NAME" ] ; then
    UPDADE_NAME="Firefly-RK3399_Android${TARGET_VERSION}_${PRODUCT_FIREFLY_NAME}_$(date -d today +%y%m%d)"
else
    UPDADE_NAME="$UPDATE_USER_NAME"
fi

if [ ! -d ${IMAGE_SRC_CUR_PATH} ] ; then
    echo "[ERROR]: Can't find image path: ${IMAGE_SRC_CUR_PATH}"
    exit 1
fi

pushd $CMD_PATH > /dev/null

ln -sf ${IMAGE_SRC_PATH} ${LINK_IMAGE_PATH}
if [ ! -f ${LOADER_PATH} ]; then
    echo "[ERROR]: Can't find loader: ${LOADER_PATH}"
    exit 2
fi

./afptool -pack ./ ${LINK_IMAGE_PATH}/tmp_${UPDADE_NAME}.img
./rkImageMaker -RK330C ${LOADER_PATH} ${LINK_IMAGE_PATH}/tmp_${UPDADE_NAME}.img "${LINK_IMAGE_PATH}/${UPDADE_NAME}.img" -os_type:androidos
rm -f ${LINK_IMAGE_PATH}/tmp_${UPDADE_NAME}.img
rm -f ${LINK_IMAGE_PATH}

popd > /dev/null

echo ""
echo "Making updateimg: ${IMAGE_SRC_CUR_PATH}/${UPDADE_NAME}.img"

exit 0
