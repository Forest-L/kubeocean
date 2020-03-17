# KubeOcean
Deploy a Kubernetes Cluster flexibly and easily
## Quick Start
### build
```shell script
git clone https://github.com/pixiake/kubeocean.git
cd kubeocean
./build
```
> Note: Docker needs to be installed before building.

### Usage
* Deploy a Allinone cluster
```shell script
./kubeocean create cluster
```
* Deploy a MultiNodes cluster
```shell script
# Create a cluster config file
./kubeocean create config
# Deploy cluster
./kubeocean create cluster --cluster-info ./cluster-info.yaml
```
###Supported
* Deploy allinone cluster
* Deploy multinodes cluster
* Add nodes (masters and nodes)
