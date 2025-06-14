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

## Usage

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
