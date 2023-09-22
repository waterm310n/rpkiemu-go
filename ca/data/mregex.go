package data

import (
	"regexp"
)

/* python ipv6 match
IPV6SEG = r'(?:(?:[0-9a-fA-F]){1,4})'

IPV6GROUPS = (
    # 1:2:3:4:5:6:7:8
    r'(?:' + IPV6SEG + r':){7,7}' + IPV6SEG,
    # 1::                                 1:2:3:4:5:6:7::
    r'(?:' + IPV6SEG + r':){1,7}:',
    # 1::8               1:2:3:4:5:6::8   1:2:3:4:5:6::8
    r'(?:' + IPV6SEG + r':){1,6}:' + IPV6SEG,
    # 1::7:8             1:2:3:4:5::7:8   1:2:3:4:5::8
    r'(?:' + IPV6SEG + r':){1,5}(?::' + IPV6SEG + r'){1,2}',
    # 1::6:7:8           1:2:3:4::6:7:8   1:2:3:4::8
    r'(?:' + IPV6SEG + r':){1,4}(?::' + IPV6SEG + r'){1,3}',
    # 1::5:6:7:8         1:2:3::5:6:7:8   1:2:3::8
    r'(?:' + IPV6SEG + r':){1,3}(?::' + IPV6SEG + r'){1,4}',
    # 1::4:5:6:7:8       1:2::4:5:6:7:8   1:2::8
    r'(?:' + IPV6SEG + r':){1,2}(?::' + IPV6SEG + r'){1,5}',
    # 1::3:4:5:6:7:8     1::3:4:5:6:7:8   1::8
    IPV6SEG + r':(?:(?::' + IPV6SEG + r'){1,6})',
    # ::2:3:4:5:6:7:8    ::2:3:4:5:6:7:8  ::8       ::
    r':(?:(?::' + IPV6SEG + r'){1,7}|:)',
    # fe80::7:8%eth0     fe80::7:8%1  (link-local IPv6 addresses with zone index)
    r'fe80:(?::' + IPV6SEG + r'){0,4}%[0-9a-zA-Z]{1,}',
)

# Reverse rows for greedy match
IPV6ADDR = '|'.join([f'(?:{g})'for g in IPV6GROUPS[::-1]])

MINMAX_IPV6ADDR = f"Min:\s*({IPV6ADDR})\s*max:\s*({IPV6ADDR})"

print(MINMAX_IPV6ADDR)
*/

const IPV6_MINMAX_REGEX = `Min:\s*((?:fe80:(?::(?:(?:[0-9a-fA-F]){1,4})){0,4}%[0-9a-zA-Z]{1,})|(?::(?:(?::(?:(?:[0-9a-fA-F]){1,4})){1,7}|:))|(?:(?:(?:[0-9a-fA-F]){1,4}):(?:(?::(?:(?:[0-9a-fA-F]){1,4})){1,6}))|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,2}(?::(?:(?:[0-9a-fA-F]){1,4})){1,5})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,3}(?::(?:(?:[0-9a-fA-F]){1,4})){1,4})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,4}(?::(?:(?:[0-9a-fA-F]){1,4})){1,3})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,5}(?::(?:(?:[0-9a-fA-F]){1,4})){1,2})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,6}:(?:(?:[0-9a-fA-F]){1,4}))|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,7}:)|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){7,7}(?:(?:[0-9a-fA-F]){1,4})))\s*max:\s*((?:fe80:(?::(?:(?:[0-9a-fA-F]){1,4})){0,4}%[0-9a-zA-Z]{1,})|(?::(?:(?::(?:(?:[0-9a-fA-F]){1,4})){1,7}|:))|(?:(?:(?:[0-9a-fA-F]){1,4}):(?:(?::(?:(?:[0-9a-fA-F]){1,4})){1,6}))|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,2}(?::(?:(?:[0-9a-fA-F]){1,4})){1,5})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,3}(?::(?:(?:[0-9a-fA-F]){1,4})){1,4})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,4}(?::(?:(?:[0-9a-fA-F]){1,4})){1,3})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,5}(?::(?:(?:[0-9a-fA-F]){1,4})){1,2})|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,6}:(?:(?:[0-9a-fA-F]){1,4}))|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){1,7}:)|(?:(?:(?:(?:[0-9a-fA-F]){1,4}):){7,7}(?:(?:[0-9a-fA-F]){1,4})))`
const IPV4_MINMAX_REGEX = `Min:\s*((?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\s*max:\s*((?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))`
const URI_REGEX = `(https|rsync):\/\/([^\/]+)(?:\/[^\/]+)*\/([^\/]+.cer)`
const ASN_REGEX = `(\d+)`
const ASN_MINMAX_REGEX = `Min:\s*(AS\d+)\s*max:\s*(AS\d+)`

// 分组匹配URI的协议，域名，证书文件名
var URI_MATCH *regexp.Regexp

// 分组匹配ASN号
var ASN_MATCH *regexp.Regexp

// 分组匹配形如`M`in:AS12345 max:AS5465`的两个AS号
var ASN_MINMAX_MATCH *regexp.Regexp

// 分组匹配形如`Min: 140.109.0.0 max: 140.138.255.255`这样的IPV4地址
var IPV4_MINMAX_MATCH *regexp.Regexp

// 分组匹配形如`Min: 2001:7fa:3:: max: 2001:7fa:4:ffff:ffff:ffff:ffff:ffff`的IPV6
var IPV6_MINMAX_MATCH *regexp.Regexp

func init() {
	URI_MATCH = regexp.MustCompile(URI_REGEX)
	ASN_MATCH = regexp.MustCompile(ASN_REGEX)
	ASN_MINMAX_MATCH = regexp.MustCompile(ASN_MINMAX_REGEX)
	IPV4_MINMAX_MATCH = regexp.MustCompile(IPV4_MINMAX_REGEX)
	IPV6_MINMAX_MATCH = regexp.MustCompile(IPV6_MINMAX_REGEX)
}
