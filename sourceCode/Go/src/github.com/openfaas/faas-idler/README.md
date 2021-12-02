# faas-idler

[![Build Status](https://travis-ci.org/openfaas-incubator/faas-idler.svg?branch=master)](https://travis-ci.org/openfaas-incubator/faas-idler)

Scale OpenFaaS functions to zero replicas after a period of inactivity

> Premise: functions (Deployments) can be scaled to 0/0 replicas from 1/1 or N/N replicas when they are not receiving traffic. Traffic is observed from Prometheus metrics collected in the OpenFaaS API Gateway.

![](./docs/faas-idler.png)

Scaling to zero requires an "un-idler" or a blocking HTTP proxy which can reverse the process when incoming requests attempt to access a given function. This is done through the OpenFaaS API Gateway through which every incoming call passes - see [Add feature: scale from zero to 1 replicas #685](https://github.com/openfaas/faas/pull/685).

faas-idler is implemented as a controller which polls Prometheus metrics on a regular basis and tries to reconcile a desired condition - i.e. zero replicas -> scale down API call.

## Building

The build requires Docker and builds a local Docker image.

```
TAG=0.2.0 make build
TAG=0.2.0 make push
```

## Usage

### Quick start

#### Swarm:

```sh
docker stack deploy func -c docker-compose.yml
```

#### Kubernetes

The faas-idler is installed as part of the [helm chart](https://github.com/openfaas/faas-netes/tree/master/chart/openfaas), make sure that you pass the argument "--set faasIdler.dryRun=false" if you want the idler to go live and make changes to the API.

#### Activating a function for scale to zero

Now decorate some functions with the label: `com.openfaas.scale.zero: "true"` and watch the idler scale them to zero. You should also change the `-dry-run` flag to `false`. For example:

```sh
faas-cli store deploy figlet --label "com.openfaas.scale.zero=true"
```

Or if using the [openfaas-operator](https://github.com/openfaas-incubator/openfaas-operator) and CRD:

```yaml
...
spec:
  labels:
    com.openfaas.scale.zero: "true"
...
```

### Configuration

* Environmental variables:

On Kubernetes the `gateway_url` needs to contain the suffix of the namespace you picked at deploy time. This is usually `.openfaas` and is pre-configured with a default.

Try using the ClusterIP/Cluster Service instead and port 8080.

| env_var               | description                                                 |
| --------------------- |----------------------------------------------------------   |
| `gateway_url`         | The URL for the API gateway i.e. http://gateway:8080 or http://gateway.openfaas:8080 for Kubernetes       |
| `prometheus_host`     | host for Prometheus |
| `prometheus_port`     | port for Prometheus |
| `inactivity_duration` | i.e. `15m` (Golang duration) |
| `reconcile_interval`  | i.e. `1m` (default value) |
| `secret_mount_path`   | default `/var/secrets/`, path from which `basic-auth-user` and `basic-auth-password` files are read |
| `write_debug`         | default `false`, set to `true` to enable verbose logging for debugging / troubleshooting |


* Command-line args

`-dry-run` - don't send scaling event 

How it works:

`gateway_function_invocation_total` is measured for activity over `duration` i.e. `1h` of inactivity (or no HTTP requests)

## Logs

You can view the logs to show reconciliation in action.

```sh
kubectl logs -n openfaas -f deploy/faas-idler
```

