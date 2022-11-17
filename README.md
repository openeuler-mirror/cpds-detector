# cpds-detector

#### 介绍
cpds-detector是为CPDS(Container Problem Detect System)容器故障检测系统开发的异常检测组件

本组件根据cpds-analyzer(容器故障/亚健康诊断组件)下发的异常规则，对集群各节点原始数据进行分析，检测节点是否存在异常

#### 从源码编译

`cpds-detector`只支持 Linux，必须使用 Go 版本 1.18 或更高版本构建。

```bash
# create a 'gitee.com/cpds' in your GOPATH/src
cd $GOPATH/gitee.com/cpds
git clone https://gitee.com/openeuler/cpds-detector.git
cd cpds-detector

make
```

编译完成后的`cpds-detector`在`out`目录中

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

