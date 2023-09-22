#!/bin/bash
#usage: ./onclick.sh handle_name parent_handle_name ipv4Resource ipv6Resource asnResource
#onclick cases: 
#./onclick.sh children1 testbed '11.14.0.0/16,12.46.0.0/16' '::/0' 'AS123'
#./onclick.sh children2 testbed '11.14.0.0/16,12.46.0.0/16' '::/0' ''
#./onclick.sh children3 testbed '11.14.0.0/16,12.46.0.0/16' '' 'AS123'
#./onclick.sh children4 testbed '11.14.0.0/16,12.46.0.0/16' '' ''
#./oneclick.sh children5 testbed  '' '::/0' 'AS123' 
#./onclick.sh children6 testbed  '' '::/0' '' 
#./onclick.sh children7 testbed  '' '' 'AS123'
#./onclick.sh children8 testbed  '' '' ''

handle_name=$1
parent_handle_name=$2
ipv4Resource=$3
ipv6Resource=$4
asnResource=$5
if [  -z "$ipv4Resource" ] && [  -z "$ipv6Resource" ] && [  -z "$asnResource" ] ;
then 
    #不符合规范
    exit 1
fi
# add ca handle
krillc add --ca "$handle_name"

krillc repo request --ca "$handle_name" > request_repo.xml
krillc pubserver publishers add --publisher "$handle_name" --request request_repo.xml > response_repo.xml
krillc repo configure --ca "$handle_name" --response response_repo.xml

krillc parents request --ca "$handle_name" > request.xml

if [ ! -z "$ipv4Resource" ] && [ ! -z "$ipv6Resource" ] && [ ! -z "$asnResource" ] ;
then 
    set -e
    krillc children add --ca "$parent_handle_name" --child "$handle_name" --ipv4 "$ipv4Resource" --ipv6 "$ipv6Resource" --asn "$asnResource" --request request.xml > response.xml
    krillc parents add --parent "$parent_handle_name" --ca "$handle_name" --response response.xml
    exit 0
fi

if [ ! -z "$ipv4Resource" ] && [ ! -z "$ipv6Resource" ] && [  -z "$asnResource" ] ;
then 
    set -e
    krillc children add --ca "$parent_handle_name" --child "$handle_name" --ipv4 "$ipv4Resource" --ipv6 "$ipv6Resource" --request request.xml > response.xml
    krillc parents add --parent "$parent_handle_name" --ca "$handle_name" --response response.xml
    exit 0
fi

if [ ! -z "$ipv4Resource" ] && [ -z "$ipv6Resource" ] && [ ! -z "$asnResource" ] ;
then 
    set -e
    krillc children add --ca "$parent_handle_name" --child "$handle_name" --ipv4 "$ipv4Resource"  --asn "$asnResource" --request request.xml > response.xml
    krillc parents add --parent "$parent_handle_name" --ca "$handle_name" --response response.xml
    exit 0
fi

if [ ! -z "$ipv4Resource" ] && [ -z "$ipv6Resource" ] && [  -z "$asnResource" ] ;
then 
    set -e
    krillc children add --ca "$parent_handle_name" --child "$handle_name" --ipv4 "$ipv4Resource"   --request request.xml > response.xml
    krillc parents add --parent "$parent_handle_name" --ca "$handle_name" --response response.xml
    exit 0
fi

if [  -z "$ipv4Resource" ] && [ ! -z "$ipv6Resource" ] && [ ! -z "$asnResource" ] ;
then 
    set -e
    krillc children add --ca "$parent_handle_name" --child "$handle_name" --ipv6 "$ipv6Resource" --asn "$asnResource" --request request.xml > response.xml
    krillc parents add --parent "$parent_handle_name" --ca "$handle_name" --response response.xml
    exit 0
fi

if [  -z "$ipv4Resource" ] && [ ! -z "$ipv6Resource" ] && [  -z "$asnResource" ] ;
then 
    set -e
    krillc children add --ca "$parent_handle_name" --child "$handle_name" --ipv6 "$ipv6Resource" --request request.xml > response.xml
    krillc parents add --parent "$parent_handle_name" --ca "$handle_name" --response response.xml
    exit 0
fi

if [  -z "$ipv4Resource" ] && [  -z "$ipv6Resource" ] && [ ! -z "$asnResource" ] ;
then 
    set -e
    krillc children add --ca "$parent_handle_name" --child "$handle_name" --asn "$asnResource" --request request.xml > response.xml
    krillc parents add --parent "$parent_handle_name" --ca "$handle_name" --response response.xml
    exit 0
fi


