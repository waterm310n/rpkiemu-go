#!/bin/bash
#set -e指遇到错误就停止执行脚本
#set -u指遇到不存在的变量就报错，默认情况是忽视
#set -o pipefail指输出从右往左的管道中第一个非零返回值
set -e -o pipefail 
source /opt/my_funcs.sh

TA_DIR=/share/ta
mkdir -p ${TA_DIR}
cd ${TA_DIR} 

for arg in $*                                          
do
    if [[ $arg == *".cer" ]] ; then 
        my_log "Waiting for TA certificate"
        my_retry 12 5 wget --no-check-certificate $arg -O ta.cer > /dev/null
    fi
   
    if [[ $arg == *".tal" ]] ; then 
        my_log "Waiting for TA trust anchor location"
        #使用-O的目的是为了支持覆盖文件的能力
        my_retry 12 5 wget --no-check-certificate $arg -O ta.tal > /dev/null
    fi
done

my_log "Launching Rsyncd"
rsync --daemon --no-detach --log-file /dev/stdout