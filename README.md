# simplecontroller

## Overview

A simple Kubernetes controller that:

1.  Watches for pods in selected namespace with the label:

    ```yaml
    simplecontroller.io/watch: "true"
    ```

2.  Patches them with the label:

    ```yaml
    simplecontroller.io/color: blue
    ```

## Usage

```bash
$ simplecontroller --help
Simple controller for labeling Pods

Usage:
  simplecontroller [flags]

Flags:
  -h, --help                    help for simplecontroller
      --kubeconfig string       Paths to a kubeconfig. Only required if out-of-cluster.
  -n, --namespace string        Namespace to watch. If empty, watch all.
  -v, --v Level                 number for the log level verbosity of the testing logger
      --vmodule pattern=N,...   comma-separated list of pattern=N log level settings for files matching the patterns
```

## Build

### Go

```bash
make {build|fmt|test|vet}
```

### Docker

```bash
make docker-build
```

### Kind

```bash
make {kind-test|kind-delete}
```
