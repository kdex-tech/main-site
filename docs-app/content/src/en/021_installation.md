---
title: "INSTALLATION"
weight: 21
---


## Installing the K-CNAS Operator

Requirements:
- a Kubernetes Cluster
- the `helm` command
- (optional) the `kubectl` client

**K-CNAS** is a [Kubernetes operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) deployed via a [Helm](https://helm.sh/) chart.

The chart installs the **[Nexus Manager](051_glossary.md#nexus-manager)**, which is responsible for reconciling KDex CRDs and setting up the cluster environment.

Helm uses the term _"release"_ for the instance of a chart running in a Kubernetes cluster. Since a chart can be instantiated multiple times in a cluster each instance needs a unique name.

Installing a chart uses the following command:

```shell
helm install <release_name> <chart>
```

Installing K-CNAS operator chart uses the following command:

```shell
helm install kcnas-operator kdex-tech/kcnas-operator
```

The `helm` command has many flags so check out the documentation.


### Prerequisites

- A running Kubernetes cluster (v1.24+ recommended).
- `helm` CLI installed.
- (Optional but recommended) `cert-manager` for [automatic TLS certificate provisioning](031_security-and-access-control.md#configuring-cert-manager).

### Helm Installation

1. Navigate to the chart directory (assuming you have the source code or an artifact):
   ```bash
   cd kdex-nexus-manager/dist/chart
   ```

2. Install the `kdex-nexus` release into your cluster:
   ```bash
   helm upgrade --install kdex-nexus . \
     --namespace kdex-system \
     --create-namespace
   ```

### What gets installed?

The Helm chart provisions:
- **Custom Resource Definitions (CRDs)**: Defines `KDexHost`, `KDexPage`, etc.
- **Nexus Manager Controller**: The operator pod that watches and reconciles resources.
- **Webhook Capabilities**: For validating CRD configurations upon creation or update.
- **Bundled Default Resources**: The chart installs default cluster-scoped components (themes, archetypes) from `kdex-nexus-manager/config/bundled` out of the box.

Once the operator pods are running, your cluster is ready to host K-CNAS sites.
