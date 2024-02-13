#!/bin/bash
#set -e指遇到错误就停止执行脚本
#set -u指遇到不存在的变量就报错，默认情况是忽视
#set -o pipefail指输出从右往左的管道中第一个非零返回值
set -e -u -o pipefail 

DATA_DIR=/home/routinator
TAL_DIR=~/.rpki-cache/tals
mkdir -p ${DATA_DIR}
mkdir -p ${TAL_DIR}

cp /opt/routinator.conf ~/.routinator.conf
cp /opt/exceptionSlurm.json ~/exceptionSlurm.json
export BANNER="Routinator setup for Krill"
source /opt/my_funcs.sh

#使用-O的目的是为了支持覆盖文件的能力
OLD_IFS="$IFS"
IFS=";"
arr=($SRC_TALS)
IFS="$OLD_IFS"

for index in ${!arr[@]}
do 
arr2=(${arr[$index]//// })
# my_retry 12 5 wget --no-check-certificate ${arr[$index]} -O ${TAL_DIR}/${arr2[1]}.tal > /dev/null
install_tal ${arr[$index]} ${TAL_DIR}/${arr2[1]}.tal
done


cd ${DATA_DIR}

my_log "Launching Routinator"
routinator \
    --strict \
    --fresh \
    --config ~/.routinator.conf \
    --rrdp-root-cert=/opt/rootCA.crt \
    -vvv \
    server \
    --rtr 0.0.0.0:3323 --http 0.0.0.0:9556 \
    --refresh 10