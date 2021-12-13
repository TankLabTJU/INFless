# Description
The directory of `workload` includes trace files and scripts for evaluation of INFless with deployed functions.  

# Instructions
- LoadGenSimClient.jar : The java runnable file for workload simulation. 
- Bursty.txt : Bursty workload from Micro Azure trace.
- Periodic.txt : Periodic workload from Micro Azure trace.
- Sporadic.txt : Sporadic workload from Micro Azure trace.
- start_load.sh : Bash script for generating workload.

## Using the script
Before using the script, make sure that two prefix steps have been done:
1. An inference function has been deployed successfully in INFless, which can be invoked by visiting the faas-gateway. For example, if resnet-50 model is deployed as an inference service, the invocation interface shoule be http://192.168.3.100:31212/resnet-50 with a picture (3\*224\*224) as input data.
2. The `LoadGen.war` has been deployed with Apache Tomcat and starts successfully with port 22222 exposed.
### Check IP address
The IP address and RMI port of `LoadGen` are defined as the two input parameters of the `start_load.sh` script. 

```bash
cd /INFless/workload
vi start_load.sh
loadGen_ip=$1
loadGen_port=$2
# ssd
java -jar LoadGenSimClient.jar 600 1 9 1800 10 0 false /INFless/workload/Periodic.txt /INFless/workload/results/ $1 $2 &
# mobilenet
java -jar LoadGenSimClient.jar 250 1 9 1800 16 0 false /INFless/workload/Periodic.txt /INFless/workload/results/ $1 $2 &
# resnet-50
java -jar LoadGenSimClient.jar 350 1 9 1800 14 0 false /INFless/workload/Periodic.txt /INFless/workload/results/ $1 $2 &
```
Suppose that the `LoadGen` is deployed in one server with 192.168.3.130. To use it, start the workload generator with the following commands. 
```bash
cd /INFless/workload
sh start_load.sh 192.168.3.130 22222
```
> Notice: Once thhe script is started, it will automatically exit after the system runtime defination. 

