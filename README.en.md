# cpds-detector

#### Description

cpds-detector is an anomaly detection component developed for the CPDS (Container Problem Detect System) container fault detection system.

This component analyzes the raw data of each node in the cluster according to the exception rules issued by cpds-analyzer (container fault/sub-health diagnosis component), and detects whether there is an exception in the node.


#### Build from source

`cpds-detector` is only supported on Linux and must be built with Go version 1.18 or higher.

```bash
# create a 'gitee.com/cpds' in your GOPATH/src
cd $GOPATH/gitee.com/cpds
git clone https://gitee.com/openeuler/cpds-detector.git
cd cpds-detector

make
```
Finally, the compiled `cpds-detector` is in the `out` directory.


#### Contribution

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request
