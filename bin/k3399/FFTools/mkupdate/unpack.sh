#!/bin/bash


usage()
{
cat << EOF
usage:
    $(basename $0) <image>
    image: source update.img name
EOF

if [ ! -z $1 ] ; then
    exit $1
fi
}

if [ $1 == '-h' ]; then
    usage 0
fi

if [ $# -ne 1 ]; then
    echo "[ERROR]: not image arg"
    usage 1
fi

CMD_PATH=$(dirname $0)
UPDATE_IMG_PATH="$(readlink -f $1)"
UNPACK_PATH="unpack"

if [ ! -f ${UPDATE_IMG_PATH} ]; then
    echo "[ERROR]: Can't find loader: $1"
    exit 2
fi

pushd $CMD_PATH > /dev/null

ln -sf ${UPDATE_IMG_PATH} update.img
rm -rf ${UNPACK_PATH}
mkdir -p ${UNPACK_PATH}
./rkImageMaker -unpack update.img ${UNPACK_PATH}
./afptool -unpack ${UNPACK_PATH}/firmware.img ${UNPACK_PATH}
rm -f ${UNPACK_PATH}/firmware.img
rm -f ${UNPACK_PATH}/boot.bin
rm -f update.img

popd > /dev/null

echo ""
echo "unpack updateimg: ${CMD_PATH}/${UNPACK_PATH}/"
tree ${CMD_PATH}/${UNPACK_PATH}

exit 0
