# Source Code for INFless
INFless is a domain-specific serverless platform. It caters to AI inference as BaaS offerings. Function developers can submit their inference source code to INFless; INFless accepts the function code of inference models and automates the deployment and scaling under varying access workloads; Users can get AI services from INFless. INFless guarantees subsecond latency for user requests and achieves high resource efficiency through the elaborately designed resource allocation and function management mechanisms. The current design of INFless relies on docker container for resource management and isolation among inference services, and it runs on Kuberneter cluster.
 
## Contents
- sourceCode: The source code of INFless for implementation and evaluation
- configuration: The cluster configuration file that needs to be loaded when INFless launches
- developer: The inference functions that could be deployed into INFless
- models: The metadata of inference models evaluated in INFless
- profiler: The model performance-resource profiles 
- workload: The function workload trace used for evaluation
- scripts: Some scripts used in the evaluation

## INFless Installation
Deployment guide of INFless for Kubernetes is available here: https://github.com/tjulym/INFless/sourceCode/README.md

Notice: INFless is fully implemented with OpenFaaS, which is a serverless frameworks built with Kubernetes. The installation of INFless is similar as OpenFaaS, and we recommend you have some preliminary knowledges about OpenFaaS and understand how it works with Kubernetes. The Deployment guide of OpenFaaS for Kubernetes is available here: https://docs.openfaas.com/deployment/kubernetes/.

## Inference Function Deployment
Guidance for function developer is available here: https://github.com/tjulym/INFless/developer/README.md

Notice: Once INFless is deployed successfully, developers could use faasdev-cli tools to upload their inference functions and build them as FaaS service. 
## Load Generator Installation
Guidance is available here: https://github.com/tjulym/INFless/sourceCode/Java/README.md

Notice: Be sure that the deployed inference functions and workload generator work well before evaluating INFless platform.
## Evaluation

The plotting source code is available here: https://github.com/tjulym/INFless/sourceCode/Matlab/README.md


##  Bug report & Question 
We have test the deployment guidances and fixed some bugs that have been found. If you have some questions about the reproduction process, please contact us via Email: ynyang@tju.edu.cn

## Reference
[1] Yanan Yang, Laiping Zhao, Yiming Li, Huanyu Zhang, Jie Li, Mingyang Zhao, Xingzhen Chen, Keqiu Li. INFless: A Native Serverless System for Low-latency, High-throughput Inference. The 27th ACM International Conference on Architectural Support for Programming Languages and Operating Systems (ASPLOS'22), Lausanne, Switzerland. Feb 2022.