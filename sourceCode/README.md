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
Password: xxx
```
After that, you should turn into the directory of `INFless` project and follow the subsequent instructions to reproduce the experimental results.
```bash
# Successfully login the private server and turn into INFless workspace
$ cd /home/tank/1_yanan/INFless/ 
$ ls
configuration  LICENSE  profiler   scripts     workload
developer      models   README.md  sourceCode
```

## 2. Build INFless and Deploy functions
INFless has been deployed on the private machine. Please check its running state using the following commands,
```bash
# You should firstly switch to the root user
$ sudo su
  [sudo] password for tank: tanklab
# List the components of INFless
$ kubectl get all -n openfaasdev 
NAME                                               READY   STATUS             RESTARTS   AGE
pod/basic-auth-plugindev-6bbffdd8c7-q8swp          1/1     Running            0          13h
pod/cpuagentcontroller-deploy-0-6687bc6f4b-47j57   0/1     Pending            0          13h
pod/cpuagentcontroller-deploy-1-75588ccd9b-9kg8x   0/1     Pending            0          13h
pod/gatewaydev-bdb695ff4-rpnjp                     2/2     Running            0          13h
pod/prometheusdev-7cb4464767-kf7v5                 1/1     Running            0          13h
...

# List the deployed inference functions in INFless
$ kubectl get all -n openfaasdev-fn
NAME                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/mobilenet   ClusterIP   10.102.207.241   <none>        8080/TCP   2m5s
service/resnet-50   ClusterIP   10.97.239.55     <none>        8080/TCP   2m17s
service/ssd         ClusterIP   10.102.74.237    <none>        8080/TCP   2m11s
```
  
## 3. Start Workload Generator
Start the load generator using the following command,

```bash
$ cd /home/tank/1_yanan/INFless/workload/
# stop the workload 
$ jps -l |grep Load |awk '{print $1}' |xargs kill -9
# start the workload
$ sh start_load.sh 192.168.1.109 22222
```
> Notice: The `start_load.sh` will run as a daemon and print some log. Please start a new terminal to run the commands in step 4.

## 4. Collect the system log and Check result

The following commands will collect INFless's runtime log and parse the results for system throughput comparison between `INFless` and its baseline (`BATCH`). 
```bash
# parse results 
$ cd /home/tank/1_yanan/INFless/workload
$ sh collect_result.sh
prefixPath:/home/tank/1_yanan/INFless/workload/
Baseline: BATCH
Total statistics QPS:54084
Scaling Efficiency: 0.5156927583326659
Throughput Efficiency: 8.112629432619582E-4
---------------------------
Baseline: INFless
Total statistics QPS:11967
Scaling Efficiency: 0.8333333333333334
Throughput Efficiency: 0.0019974242290713607
```

The result shows that `INFless` achieves 2.5x higher throughput than BATCH as in Figure 11 (0.0019 v.s. 0.0008).