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

