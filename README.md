# EC2 Benchmark - Tool to deploy and test EC2 instances
## Why?
I wanted to test the performance of EC2 instances. It is hard to find differences in performance for CPU for different instance types. This tool helps to deploy and test EC2 instances. It uses sysbench as benchmark. 
There is in memory database which keeps data.

## How to use?
Add IAM key and secret to file ./credentials.
Change the region in the file ./config.
Add your private key, ami (because of using rpm and yum it must be some RHEL based ami) and subnet, security group to the file config.
```