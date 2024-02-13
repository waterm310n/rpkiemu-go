# tar命令

`kubectl cp `系列命令是通过`tar`实现的，因此它要求容器必须包含`tar`命令



## 下载文件
```bash
tar cf - <srcfile> | cat
```
以一个文本文件hello.txt为例，其中的内容如下
```txt
hello！this is hello.txt and line 1
hello！this is hello.txt and line 2

```
然后执行上述命令可以得到如下内容
```bash
master@master:~/test$ tar cf - hello.txt | cat
hello.txt0000664000175000017500000000011214516457177012307 0ustar  mastermasterhello！this is hello.txt and line 1
hello！this is hello.txt and line 2
```
可以看到该命令能够将源文件打包并输出的stdout端。因此通过在容器执行tar，然后通过重定向的方式，即可在本地获取输出，进而进行解压缩包处理。以此实现文件下载功能。

## 上传文件
```bash
tar -xmf -
```
基本原理同上。

## 额外知识
在shell中，单独的短横线`-`，表示stdout
```bash
echo hello world | cat -
```