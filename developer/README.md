# Decription
The directory of `developer` includes the inference model files and scripts for uploading and deploying functions with INFless.
# Instructions
## Faas-cli login 
Use the following commands to login INFless, please make sure that `faasdev-cli` executable binary file exists in directory of `/usr/local/bin` of master server.
```bash 
$ faasdev-cli login --password admin
WARNING! Using --password is insecure, consider using: cat ~/faas_pass.txt | faas-cli login -u user --password-stdin
Calling the OpenFaaS server to validate the credentials...
WARNING! Communication is not secure, please consider using HTTPS. Letsencrypt.org offers free SSL/TLS certificates.
credentials saved for admin http://127.0.0.1:31212
```
## Building inference functions
Use the following commands to build the docker image for the inference model.
``` nash
$ cd /INFless/developer/servingFunctions
$ faasdev-cli build -f resnet-50.yml
..................build begin........................
Removing intermediate container 9c772362382b
 ---> 3dd464dd7882
Step 24/24 : CMD ["sh ./tfserving_entrypoint.sh"]
 ---> Running in f8bbb78253ff
Removing intermediate container f8bbb78253ff
 ---> 8f76822ab420
Successfully built 8f76822ab420
Successfully tagged resnet-50:latest

# list the function images
$ docker images |grep resnet
resnet-50   latest   10b7bd7e8d66   9 days ago   2.67GB
```
> Notice: Here we just give an example to show how to build function images with faasdev-cli tool. For the other inference models, just replace the .yml file name in these commands.

## Deploying inference functions
When the function images is build , use the following commands to deploy function with INFless.
``` bash
$ faasdev-cli deploy -f resnet-50.yml
Deploying: resnet-50.
WARNING! Communication is not secure, please consider using HTTPS. Letsencrypt.org offers free SSL/TLS certificates.

Deployed. 202 Accepted.
URL: http://127.0.0.1:31212/function/resnet-50
```
If you see the above message, this means the inference function of resnet-50 has been deployed successfully.