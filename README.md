# simplecontroller

## Overview

A simple Kubernetes controller that:

1.  Watches for pods across all namespaces with the label:

    ```yaml
    simplecontroller.io/watch: "true"
    ```

2.  Patches them with the label:

    ```yaml
    simplecontroller.io/color: blue
    ```

## Build

### Docker

```bash
docker build -t simplecontroller:0.1.0 .
```

## Test

### kind

```bash
kind create cluster
kind load docker-image simplecontroller:0.1.0
kubectl apply -f config/default
kind delete cluster
```
