# Pingdom StatusPage integration

## Introduction

This chart installs a controller for [Traefik](https://traefik.io/) global [error pages](https://docs.traefik.io/configuration/backends/kubernetes/#global-default-backend-ingresses). This service creates a global ingress rule which match all unmatched ingress rules. You can use any container you want. The container must expose the error page at the `/` path. The only error which is handled by traefik over this service by default is the `404` error.


## Prerequisites

-  kubernetes 1.11+

## Installing the Chart

To install the chart with the release name `pingdom-statuspage-integration`:

```bash
$ helm install docplanner/pingdom-statuspage-integration --name pingdom-statuspage-integration --values=my-values.yaml
```

## Uninstalling the Chart

To uninstall/delete the `pingdom-statuspage-integration` deployment:

```bash
$ helm delete pingdom-statuspage-integration
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the pingdom-statuspage-integration chart and their default values.

| Parameter                                  | Description                               | Default                            |
| ------------------------------------------ | ----------------------------------------- | ---------------------------------- |
| `replicaCount` | replica count | `1`|
| `updateStrategy` | Deployment update strategy | `type: RollingUpdate` |
| `image.repository` | Docker image repo | `docplanner/pingdom-statuspage-integration`|
| `image.tag` | Docker image tag | `latest`|
| `image.pullPolicy` | Docker image pull policy| `IfNotPresent`|
| `resources` | Resources | `{}`|
| `nodeSelector` | NodeSelector | `{}`|
| `tolerations` | Tolerations | `[]`|
| `affinity` | Affinity | `{}`|
| `ingress.enabled` | Ingress enabled | `false` |
| `ingress.annotations` | Ingress annotations | `{}` |
| `ingress.path` | Ingress path | `/` |
| `ingress.hosts` | Ingress hosts | `[]` |
| `ingress.tls` | Ingress TLS | `[]` |
| `config.secret` | secret used in communication from Pingdom | `""` |
| `config.statuspageToken` | StatusPage API key | `""` |
| `config.maxRetries` | Number of retries | `"2"` |
| `config.retryInterval` | Numer of seconds between retries | `"10"` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```bash
$ helm install --name pingdom-statuspage-integration --set ingress.tag=v1.0.0 docplanner/pingdom-statuspage-integration
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart.