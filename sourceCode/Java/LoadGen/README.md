# LoadGen
The load generator is a maven project writtern with Java 1.8, it works in open-loop[1] mode and uses httpclient and threadpool as the key components in its implementation. It also has a web GUI to show the realtime latency changes (99th tail-latency, QPS, RPS, etc.). Meanwhile, the load generator provides RMI interface for external call and supports dynamic worload during the evaluation.

* support thousands of concurrent requests per second in single node
* support web GUI for real-time latency monitoring
* support dynamic workload in open-loop
* support RMI interface for external call

##  Code architecture

```bash
/LoadGen
 |----/build/sdcloud.war   # the executable war package
 |----/src/main
           |----/java/scs/
           |          |----/controller/*      # MVC controller layer
           |          |----/pojo/*	 # entity bean layer
           |          +----/util
           |                |----/format/*  # format time and data
           |                |----/loadGen
           |                |     |----/driver/*  # drivers for generating request workloads
           |                |     |----/recordDriver/* # drivers for record request metrics
           |                |     +----/strategy/* 
           |                |     +----/threads/* # threads tools
           |                |----/respository/*   # in-memory data storage  
           |                |----/rmi/*       # RMI service and interfaces   
           |                +----/tools/*   # some tools
           |----/resources/*  # configuration files
           +----/webapp/*     # GUI pages
```

## Hardware environment
In our experiment environment, the configuration of nodes are shown as below:

| hostname | description | IP | role |
| ---- | ---- | ---- | ----|
| master-node | The master node of Kubernetes cluster, which has deployed INFless system. Inference functions that deployed with INFless could be invoked by HTTP interface, e.g., http://192.168.3.100:31212/resnet-50 | 192.168.3.100 | Inference service |
| agent-node1 | The LoadGen tool is deployed in this node | 192.168.3.130 | load generator |

##  Build load generator
The load generator is writen in Java, it can be deployed in container or host machine, and we need install Java JDK and apache tomcat before using it 
### Step 1: Modify the propority file

Firstly, you need to modify the configuration file in `LoadGen`
```bash
$ vi LoadGen/src/main/resources/conf/sys.properties
```
Modify the content of `sys.properties` as the configruations in your local environment:
```python
56 # URL of web inference service that can be accessed by http 
57 exampleURL=http://www.github.com 
58 # node IP that deployed load generator
59 serverIp=192.168.3.130     
60 # the port of RMI service provided by load generator 
61 rmiServiceEnable=false
62 rmiPort=22222
63 # window size of latency recorder, which can be seen from GUI page
64 windowSize=60
65 # record interval of latency, default: 1000ms
66 recordInterval=1000
```

Secondly, you need to modify the source code in directory `/LoadGen/src/main/java/scs/util/loadGen/driver/serverless/`.
This directory includes the driver class files of the deployed inference functions, and can be extended to other inference models easily. We list the support inference models as below:
```bash
$ ls LoadGen/src/main/java/scs/util/loadGen/driver/serverless/
CatdogFaasTFServingDriver.java 
LstmFaasTFServingDriver.java
MobileNetFaasTFServingDriver.java
ResNet50FaasTFServingDriver.java
SsdFaasTFServingDriver.java
Textcnn20FaasTFServingDriver.java
Textcnn69FaasTFServingDriver.java
YamNetFaasTFServingDriver.java
```
Then, you need to modify the .java file for each inference model to overwrite the default ip address (xxx.xxx.xxx.xxx) of the `faas-gateway` to that in your deployment environment.

```pyhton
31	@Override
32	protected void initVariables() {
33		httpClient=HttpClientPool.getInstance().getConnection();
34		queryItemsStr=Repository.resNet50FaasBaseURL;
35		jsonParmStr=Repository.resNet50ParmStr; 
36		queryItemsStr=queryItemsStr.replace("Ip","xxx.xxx.xxx.xxx");
37		queryItemsStr=queryItemsStr.replace("Port","31212");
38	}
```
> Note: Before using the loadGen, please be sure that the gateway ip addresses in the driver class files are correct, otherwise the loadGen cannot work.


### Step 2: Install Java JDK
> Note: if you have installed Java JDK in your local environment, this step can be skipped.

Download Java JDK and install it to `/usr/local/java/`
```bash
$ wget https://download.oracle.com/otn/java/jdk/8u231-b11/5b13a193868b4bf28bcb45c792fce896/jdk-8u231-linux-x64.tar.gz
$ tar -zxvf jdk-8u231-linux-x64.tar.gz /usr/local/java/
```
Modify the `/etc/profile` file
```bash
$ vi /erc/profile
```
Config Java environment variables, append the following content into the file
```bash
export JAVA_HOME=/usr/local/java/jdk1.8.0_231
export JRE_HOME=${JAVA_HOME}/jre
export CLASSPATH=.:${JAVA_HOME}/lib:${JRE_HOME}/lib
export PATH=$PATH:${JAVA_HOME}/bin 
```
Enable the configuration
```
$ source /etc/profile
$ java -version
java version "1.8.0_231"
Java(TM) SE Runtime Environment (build 1.8.0_231-b12)
Java HotSpot(TM) 64-Bit Server VM (build 25.231-b12, mixed mode)
```

### Step 3: Compile LoadGen project with Maven
> Note: You can build the source code with `Maven` or `eclipse IDE`.

Download maven binary code and install to `/usr/local/maven/`
```bash
$ mkdir /usr/local/maven
$ cd /usr/local/maven
$ wget https://mirrors.tuna.tsinghua.edu.cn/apache/maven/maven-3/3.8.1/binaries/apache-maven-3.8.1-bin.tar.gz
$ tar -zxvf apache-maven-3.8.1-bin.tar.gz 
```

Varifing the configuration
```bash
$ /usr/local/maven/apache-maven-3.8.1/bin/mvn -v
Apache Maven 3.8.1 (05c21c65bdfed0f71a2f2ada8b84da59348c4c5d)
Maven home: /home/tank/yanan/tools/apache-maven-3.8.1
Java version: 1.8.0_265, vendor: Private Build, runtime: /usr/lib/jvm/java-8-openjdk-amd64/jre
Default locale: zh_CN, platform encoding: UTF-8
OS name: "linux", version: "4.15.0-118-generic", arch: "amd64", family: "unix"
```

Build project with `Maven`
```
$ cd LoadGen & ls -l
-rw-rw-r-- 1 tank tank 6.6K July  28 10:04 pom.xml
drwxrwxr-x 4 tank tank 4.0K July  28 09:43 src
$ /home/tank/yanan/tools/apache-maven-3.8.1/bin/mvn package -Dmaven.skip.test=true
```
then, you will see the following output for packing project with maven tool
```
[INFO] Scanning for projects...
[INFO] --------------------------< sdc.tju:loadGen >---------------------------
[INFO] Building loadGen Maven Webapp 0.0.1-SNAPSHOT
[INFO] --------------------------------[ war ]---------------------------------
Downloading from central: https://repo.maven.apache.org/maven2/net/sf/json-lib/json-lib/2.4/json-lib-2.4-jdk15.jar
Downloaded from central: https://repo.maven.apache.org/maven2/net/sf/json-lib/json-lib/2.4/json-lib-2.4-jdk15.jar (159 kB at 87 kB/s)
[INFO] --- maven-resources-plugin:2.6:resources (default-resources) @ loadGen ---
Downloading from central: https://repo.maven.apache.org/maven2/org/apache/maven/maven-plugin-api/2.0.6/maven-plugin-api-2.0.6.pom
Downloaded from central: https://repo.maven.apache.org/maven2/org/apache/maven/maven-plugin-api/2.0.6/maven-plugin-api-2.0.6.pom (1.5 kB at 3.8 kB/s)
Downloading from central: https://repo.maven.apache.org/maven2/org/apache/maven/maven/2.0.6/maven-2.0.6.pom
...
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time:  22.650 s
[INFO] Finished at: 2021-07-28T10:43:37+08:00
[INFO] ------------------------------------------------------------------------
```
when it is successfully built, the `LoadGen.war` will be generated in the targe directory
```
$ ls target
drwxr-xr-x 6 root root     4096 7月  28 10:04 classes
drwxr-xr-x 4 root root     4096 7月  28 10:43 LoadGen
-rw-r--r-- 1 root root 26066047 7月  28 10:43 LoadGen.war
drwxr-xr-x 2 root root     4096 7月  28 10:43 maven-archiver
drwxr-xr-x 3 root root     4096 7月  28 10:05 maven-status
```
### Step 4: Install Apache Tomcat
Download apache tomcat and install to `/usr/local/tomcat/`
```bash
$ wget http://mirrors.tuna.tsinghua.edu.cn/apache/tomcat/tomcat-8/v8.5.47/bin/apache-tomcat-8.5.47.tar.gz
$ tar -zxvf apache-tomcat-8.5.47.tar.gz /usr/local/tomcat/
```
### Step 5: Deploy LoadGen.war into Tomcat


Deploy the web package `LoadGen.war` into tomcat webapp/
```bash
$ cd LoadGen
$ mv LoadGen/target/LoadGen.war /usr/local/tomcat/apache-tomcat-8.5.47/webapp
$ /usr/local/tomcat/apache-tomcat-8.5.47/bin/startup.sh
```
Validate if the depolyment is successful
```bash
$ curl http://localhost:8080/loadGen/
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
...
welcome to the load generator page!
...
</body>
```
 

### Step 6: Test the load generator
Open the web browser and visit url `http://192.168.3.130:8080/loadGen/`
The GUI page is shown as below:

![realtime Latency](https://github.com/yananYangYSU/book/blob/master/welcome.png?raw=true)

If you could see this page, that means the `loadGen` tool is deployed successfully.

### Step 6: Work with the LoadGenSimClient.jar 

> Notice: We have prepared available bash scripts for generating workload with `LoadGenSimClient.jar` and `LoadGen`, see the [README.md](https://github.com/TankLabTJU/INFless/tree/main/workload) in directory `/INFless/workload/`. The following instructions is unnecessary to read but could help you to understand how it works.

`LoadGenSimClient` is a loader simulator that generates diverse workload pattern, which is available in directory `/INFless/workload/LoadGenSimClient.jar`. 
`LoadGenSimClient` communicates with `LoadGen` with Java RMI interface, the compiled .jar file can be used directly.


```bash
$ cd workload
$ java -jar LoadGenSimClient.jar 600 1 9 1800 16 0 false /../workload/Periodic.txt /../simLoad/ 192.168.1.129 22222 
```
The command input includes 11 parameters, which are listed as below:
- maxSimQPS (>0) : The peak value of workload.
- simQpsPeekRate (0,1] : The scaling factor of peak request arrvial rate.
- simQpsRemainInterval (s) : The remain time for each request load point.
- systemRunTime (s) : The time length of workload persistence time.
- serviceId (int+) : The index of inference function in LoadGen.
- concurrency (0/1) : The workload concurrency flag.
- recordLatency ({true|false}) : The latency record flag.
- realQpsFilePath (/xxx/xx.txt) : The directory of workload trace file.
- resultFilePath (/xxx/xx/) : The result directory.
- rmiServerIp : The ip address of RMI server, i.e., the LoadGen.
- rmiServerPort : The port of RMI server, and the default value is 22222.
 
The index of the supported inference functions in `LoadGen` are listed as below:
```python
	/**
	 * maps the loaderIndex with the loaderDriver instance
	 * @param loaderIndex
	 * @return
	 */
	private static LoaderDriver loaderMapping(int loaderIndex){
        ...
		if(loaderIndex==10){
			return new LoaderDriver("resnet-50", ResNet50FaasTFServingDriver.getInstance());
		}
		if(loaderIndex==11){
			return new LoaderDriver("textcnn-69", Textcnn69FaasTFServingDriver.getInstance());
		}
		if(loaderIndex==12){
			return new LoaderDriver("textcnn-20", Textcnn20FaasTFServingDriver.getInstance());
		}
		if(loaderIndex==13){
			return new LoaderDriver("lstm-maxclass-2365", LstmFaasTFServingDriver.getInstance());
		}
		if(loaderIndex==14){
			return new LoaderDriver("ssd", SsdFaasTFServingDriver.getInstance());
		}
		if(loaderIndex==15){
			return new LoaderDriver("yamnet", YamNetFaasTFServingDriver.getInstance());
		}
		if(loaderIndex==16){
			return new LoaderDriver("mobilenet", MobileNetFaasTFServingDriver.getInstance());
		}
		if(loaderIndex==17){
			return new LoaderDriver("catdog", CatdogFaasTFServingDriver.getInstance());
		}
		if(loaderIndex==18){
			return new LoaderDriver("catdogAliyun-Faas", CatdogAliyunFaasTFServingDriver.getInstance());
		}
		if(loaderIndex==19){
			return new LoaderDriver("socialNetwork", SocialNetworkDriver.getInstance());
		}
		if(loaderIndex==20){
			return new LoaderDriver("solrSearch", SolrSearchDriver.getInstance());
		}
	    ...
		return null;
	} 

```