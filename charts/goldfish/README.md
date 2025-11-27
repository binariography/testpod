# Goldfish
Goldfish is a web app that can be used as a playground for practicing SRE

## Installing the Chart

The Podinfo charts are published to
[GitHub Container Registry](https://github.com/binariography/goldfish/pkgs/container/charts%2Fgoldfish)
and signed with [Cosign](https://github.com/sigstore/cosign) & GitHub Actions OIDC.

To install the chart with the release name `my-release` from GHCR:

```console
$ helm upgrade -i my-release oci://ghcr.io/binariography/charts/goldfish
```

Alternatively, you can install the chart from GitHub pages:

```console
$ helm repo add goldfish https://binariography.github.io/goldfish

$ helm upgrade -i my-release goldfish/goldfish
```

The command deploys goldfish on the Kubernetes cluster in the default namespace.
The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the goldfish chart and their default values.


| Parameter                         | Default                | Description                                                                                                            |
| --------------------------------- | ---------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `replicaCount`                    | `1`                    | Desired number of pods                                                                                                 |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```console
$ helm install my-release goldfish/goldfish \
  --set=serviceMonitor.enabled=true,serviceMonitor.interval=5s
```

To add custom annotations you need to escape the annotation key string:

```console
$ helm upgrade -i my-release goldfish/goldfish \
--set podAnnotations."appmesh\.k8s\.aws\/preview"=enabled
```

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
$ helm install my-release goldfish/goldfish -f values.yaml
```

> **Tip**: You can use the default [values.yaml](values.yaml)

