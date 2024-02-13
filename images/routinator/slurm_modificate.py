import json
import argparse

if __name__ == "__main__" :
    # 使用样例
    # python3 slurm_modificate.py -f exceptionSlurm.json -A -p "10.11.1.0/24" -a "123" -m "28" -c ""
    parser = argparse.ArgumentParser(description="Modificate Slurm ")
    parser.add_argument("-f","--file",help="slurm file path")
    parser.add_argument("-p","--prefix",help="ip prefix")
    parser.add_argument("-a","--asn",help="asn")
    parser.add_argument("-m","--maxPrefixLength",help="maxPrefixLength")
    parser.add_argument("-c","--comment",help="comment")
    parser.add_argument("-A","--assertion",action="store_true",help="prefixAssertions")
    parser.add_argument("-F","--filter",action="store_true",help="prefixFilters")
    args = parser.parse_args()
    if args.assertion:
        with open(args.file,"r") as f:
            object = json.loads(f.read())
            object["locallyAddedAssertions"]["prefixAssertions"].append(
                {"prefix":args.prefix,"asn":int(args.asn),"maxPrefixLength":int(args.maxPrefixLength),"comment":args.comment})
        with open(args.file,"w") as f:
            f.write(json.dumps(object,indent=4))
    elif args.filter:
        with open(args.file,"r") as f:
            object = json.loads(f.read())
            object["validationOutputFilters"]["prefixFilters"].append(
                {"prefix":args.prefix,"asn":int(args.asn),"maxPrefixLength":int(args.maxPrefixLength),"comment":args.comment})
        with open(args.file,"w") as f:
            f.write(json.dumps(object,indent=4))