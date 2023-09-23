# 数据生成

生成的数据呈现树形结构,其文件结构大致如下所示

```
.
├── mripe-ncc-ta
├── mripe-ncc-ta_children
│   ├── m102ed0852ea4b4700eae91a42e9f7e0fe53497e6
│   ├── m102ed0852ea4b4700eae91a42e9f7e0fe53497e6_children
│   ├── ma4250f3c598917bc119f7ddd423595bb7251181c
│   ├── ma4250f3c598917bc119f7ddd423595bb7251181c_children
│   │   ├── m1jEEXgH74HRy26tqz6rGi2p4R28
│   │   ├── mCnL-b5z5_hDy2tw1BsZhkUG21hY
│   │   ├── mNtmoRKu7irFWj9OV4aAdDNTs5lU
│   │   ├── mjJSdv3GcxMvclDgsa6pNVv1_v1w
│   │   ├── mjt86YvmixlvlDi5KNkl_9WRWSkk
│   │   └── mxEVPS6F2oPWaqya7YmwaYMAFK4s
│   └── maa2304f03c4e807dae51fd33008c93692b973c40
│       ├── mzpyS25FDQF186nbbBJlhIVgSs3E
│       └── mzxwWMquEOVl0Q9l4xgbK1-rubso
└── readme.md
```

编写成树形结构主要的目的就是为了能够更好的利用多线程能力，加快运行
