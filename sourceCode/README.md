# Description  
The source code of INFless implementation and evaluation consists of three parts:  

1. `Go`: The source code of INFless system, which is fully implemented in [OpenFaaS](https://docs.openfaas.com/deployment/kubernetes/). OpenFaaS includes three important components: faas-cli, gateway and faas-netes. They have been modified in INFless to be highly adaptive with AI inference. Guidelines to build and install INFless are listed as below:
- Guideline to build and install `faasdev-cli` tool is available  [here](https://github.com/TankLabTJU/INFless/tree/main/sourceCode/Go/src/github.com/openfaas/faas-cli/README.md). 
- Guideline to build and install gateway is available  [here](https://github.com/TankLabTJU/INFless/tree/main/sourceCode/Go/src/github.com/openfaas/faas/gateway/README.md). 
- Guideline to build and install faas-netes is available  [here](https://github.com/TankLabTJU/INFless/tree/main/sourceCode/Go/src/github.com/openfaas/faas-netes/README.md). 
2. `Java`: The source code of workload generator, which is written with Java 1.8.0. When inference functions are deployed successfully with INFless, the workload generator could be used to generator requests as virtual clients to invoke them.
- Guideline to build and install `LoadGen` is available [here](https://github.com/TankLabTJU/INFless/tree/main/sourceCode/Java/LoadGen).
3. `Matlab`: The source code of evaluation results plots, which is written with `MATLAB` 2014a.  
- Guideline to run the code files is available at directory `/INFless/sourceCode/Matlab/INFless/`.


# Instructions
We provide the steps to build and launch INFless system as below:
## 1. Prepare the test environment.
There are two ways to choose:
1. Using the remote test environment in our private cluster. We have configured a SSH reverse proxy between one public cloud server and the physical machine in our private cluster, which allows you to access the test environment through the public server. To do it, you should firstly login the public server.
```bash
# Try to login the public server
$ ssh root@47.106.xxx.xxx
Password: xxxx
```
When you login the public server successfully, you should then  access our private machine.

```bash
# Successfully login the public server and try to login the private server
$ ssh tank@localhost -p 8387
Password: xxxx
```
After that, you should turn into the directory of INFless project and follows the subsequent instructions on Github.
```bash
# Successfully login the private server and turn into INFless workspace
$ cd /home/tank/1_yanan/INFless/ 
$ ls
configuration  LICENSE  profiler   scripts     workload
developer      models   README.md  sourceCode
```
2. For anther way, we recommend you to prepare the specified hardwares in a local bare-metal server cluster, and follow the instructions for building and running the system. The detailed documents are available on Github. 
 
## 2. Build and launch INFless framework
INFless is fully implemented within OpenFaaS, which is a FaaS platform runs on Kubernetes. To install INFless, firstly, you should compile and build the docker images for each component. Using the following commands to compile codes for faasdev-cli, faas-gateway and faas-netes.
``` bash
# build faasdev-cli
$ cd sourceCode/Go/src/github.com/openfaas/
$ ls
faas  faas-cli  faas-idler  faas-netes

# build faasdev-cli
$ cd INFless/sourceCode/Go/src/github.com/openfaas/faas-cli
$ make 
$ cp faasdev-cli /usr/local/bin 
$ chmod 777 /usr/local/bin/faasdev-cli 

# build gateway
$ cd INFless/sourceCode/Go/src/github.com/openfaas/faas/gateway
$ make

# build faas-netes
$ cd INFless/sourceCode/Go/src/github.com/openfaas/faas-netes
$ kubectl create -f namespace.yml
$ make
```
The prepared image list should be like this:
``` bash
$ docker images |grep dev
openfaas/faas-netes  latest-dev  8f76822ab420   2 days ago   65.6MB
openfaas/gateway     latest-dev  ce08c7020a45   12 days ago  30MB
openfaas/faas-cli    latest-dev  2e71371d741a   7 weeks ago  31.2MB
```
Then, deploy the INFless system on top of Kubernetes cluster.
```bash
# cluster configuration files
$ cd INFless/sourceCode/Go/src/github.com/openfaas/faas-netes
$ cp yml/clusterCapConfig-dev.yml /root/yaml

# model profiler files
$ mkdir /root/yaml
$ cd INFless/
$ cp -r profiler/ /root/yaml/

# create namespace
$ kubectl create -f namespace.yml
# create basic-auth secret
$ kubectl -n openfaasdev create secret generic basic-auth \
--from-literal=basic-auth-user=admin \
--from-literal=basic-auth-password=admin

# deploy components
$ kubectl delete -f yml/inuse
$ kubectl apply -f yml/inuse
```

## 3. Deploy infererence functions
The inference model files are stored in directory of `INFless/developer/servingFunctions/`
```bash
$ cd INFless/developer/servingFunctions/
$ faasdev-cli build -f resnet-50.yml
$ faasdev-cli deploy -f resnet-50.yml
```
