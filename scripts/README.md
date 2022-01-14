
# Instructions
The following steps will reproduce the results in Figure 11 (system throuhgput comparison between INFless and its baseline).
## 1. Login the test environment
We have configured a SSH reverse proxy between one public cloud server and the physical machine in our private cluster, which has already deployed the `INFless` system and workload functions. Please login the public server using the following command,
```bash
# Try to login the public server
$ ssh root@47.106.xxx.xxx
Password: xxx
```
After you login the public server, please access our private machine using the following commands,

```bash
# Successfully login the public server and try to login the private server
$ ssh tank@localhost -p 8387
Password: tanklab
```
After that, you should turn into the directory of `INFless` project and follow the subsequent instructions to reproduce the experimental results.
```bash
# Successfully login the private server and turn into INFless workspace
$ cd /home/tank/1_yanan/INFless/ 
$ ls
configuration  LICENSE  profiler   scripts     workload
developer      models   README.md  sourceCode
```

## 2. Build and launch INFless framework
INFless is fully implemented within OpenFaaS, which is a FaaS platform runs on Kubernetes. To install INFless, firstly, you should compile and build the docker images for each component. Using the following commands to compile codes for faasdev-cli, faas-gateway and faas-netes.
```bash
# You should firstly switch to the root user
$ sudo su
  [sudo] password for tank: tanklab
  
# Compile and install INFless
$ cd /home/tank/1_yanan/INFless/sourceCode/Go/src/github.com/openfaas/
$ ls
faas  faas-cli  faas-idler  faas-netes

# Compile gateway
$ cd /home/tank/1_yanan/INFless/sourceCode/Go/src/github.com/openfaas/faas/gateway
$ make

# Compile faas-netes
$ cd /home/tank/1_yanan/INFless/sourceCode/Go/src/github.com/openfaas/faas-netes
$ make

# Install INFless on Kubernetes
$ cd /home/tank/1_yanan/INFless/sourceCode/Go/src/github.com/openfaas/faas-netes
$ kubectl apply -f yaml/

# List the components of INFless
$ kubectl get all -n openfaasdev 
NAME                                               READY   STATUS             RESTARTS   AGE
pod/basic-auth-plugindev-6bbffdd8c7-q8swp          1/1     Running            0          13h
pod/cpuagentcontroller-deploy-0-6687bc6f4b-47j57   0/1     Pending            0          13h
pod/cpuagentcontroller-deploy-1-75588ccd9b-9kg8x   0/1     Pending            0          13h
pod/gatewaydev-bdb695ff4-rpnjp                     2/2     Running            0          13h
pod/prometheusdev-7cb4464767-kf7v5                 1/1     Running            0          13h
...
```
## 3. Deploy infererence functions
The inference model files are stored in directory of `/home/tank/1_yanan/INFless/developer/servingFunctions/`
```bash
$ source /etc/profile
$ cd /home/tank/1_yanan/INFless/developer/servingFunctions/
# ssd, latency target 300ms
$ faasdev-cli deploy -f ssd.yml
# mobilenet, latency target 200ms
$ faasdev-cli deploy -f mobilenet.yml
# resnet-50, latency target 300ms
$ faasdev-cli deploy -f resnet-50.yml
```
The deployed inference functions can be listed as follows:
```

$ kubectl get all -n openfaasdev-fn
NAME                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/mobilenet   ClusterIP   10.102.207.241   <none>        8080/TCP   2m5s
service/resnet-50   ClusterIP   10.97.239.55     <none>        8080/TCP   2m17s
service/ssd         ClusterIP   10.102.74.237    <none>        8080/TCP   2m11s
```
  
## 4. Start Workload Generator
Start the load generator using the following command,

```bash
$ cd /home/tank/1_yanan/INFless/workload/
# stop the workload 
$ jps -l |grep Load |awk '{print $1}' |xargs kill -9
# start the workload
$ sh start_load.sh 192.168.1.109 22222
```
> Notice: The `start_load.sh` will run as a daemon and print some log. Please start a new terminal to run the commands in step 4.

## 5. Collect the system log and Check result

The following commands will collect INFless's runtime log and parse the results for system throughput comparison between `INFless` and its baseline (`BATCH`). 
```bash
# parse results 
$ cd /home/tank/1_yanan/INFless/workload
$ sh collect_result.sh
prefixPath:/home/tank/1_yanan/INFless/workload/
Baseline: BATCH
Total statistics QPS:52810
Scaling Efficiency: 0.498
Throughput Efficiency: 8.23954E-4
---------------------------
Baseline: INFless
Total statistics QPS:12068
Scaling Efficiency: 0.8135
Throughput Efficiency: 0.001874
```

The result shows that `INFless` achieves >2x higher throughput than BATCH as in Figure 11 (0.00187 v.s. 0.00082).