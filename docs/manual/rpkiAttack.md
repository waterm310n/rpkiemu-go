# rpki攻击

在完成rpki网络配置后，输入如下命令，执行攻击
```bash
./rpkiemu-go ca attack -i < 攻击文件 >
```

## 攻击文件配置

攻击文件的数据格式为json，其中为一个Attack类型的数组。

### Attack类型
|键| 说明| 类型 | 可能的值 |
|--|--|--|--|
|attack_type| 攻击类型 | 字符串|  REVOCATE\|MODIFICATE\|DELETE\|CORRUPT\|INJECT |
|attack_object| 攻击对象 |字符串| RESOURCECERT\|ROA |
|attack_data| 用于执行本次攻击所需要的数据 |AttackData类型数组| -|

### AttackData类型
|键| 说明| 类型 | 可能的值 |
|--|--|--|--|
|handle_name| handle名称 | 字符串 | mripe-ncc-ta  |
|publish_point| handle发布点 |字符串|  ripe |
|parent_handle_name| handle的上级名称 |字符串| mripe-ncc-ta |
|parent_publish_point| handle的上级发布点 |字符串| ripe |
|asn | 自治域号 |整型| 8393 |
|asnes | 自治域号数组，只在注入攻击中起作用 |整型数组| [28,8393] |
| ipv4_resource| ipv4资源 |字符串数组| ["216.205.160.0/19"] |
| ipv6_resource| ipv6资源 |字符串数组| ["2001:5000::/20"] |
| pre_bindings| 修改前的ASN-IP对 |Binding类型数组| {"ip": "193.193.235.0/24","asn": 8393} |
| bindings| ASN-IP对 | Binding类型数组| 同上 |
| after_bindings| 修改后的ASN-IP对 |Binding类型数组| 同上 |

### Binding类型
|键| 说明| 类型 | 可能的值 |
|--|--|--|--|
|ip|ip资源|字符串| "216.205.160.0/19" |
|asn|自治域号|整型| 8393|

## 示例
可见[攻击示例](../../attack_examples/kz-attack/readme.md)