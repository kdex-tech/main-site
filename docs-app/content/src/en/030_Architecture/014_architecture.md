---
title: "K-CNAS Architecture"
weight: 14
---


## High-Level Architecture

The architecture of K-CNAS is built upon [Kubernetes Custom Resource Definitions (CRDs)](../050_api-crd-reference.md).

**Key properties:**
- **Everything is a CRD**.
    - _If you don't know what a CRD is,_ see [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- A **[KDexHost](#host-centric-model)** represents a web property (site) reachable at one or more domains.
- **Pages** are defined as [KDexPage](../051_glossary.md#kdexpage) CRDs attached to a host.
- **Design** is defined by [KDexTheme](../051_glossary.md#kdextheme) CRD attached to the host.
- **Structure** is defined by [KDexPageArchetype](../051_glossary.md#archetype), [KDexPageHeader](../051_glossary.md#kdexpageheader), [KDexPageFooter](../051_glossary.md#kdexpagefooter), and [KDexPageNavigation](../051_glossary.md#kdexpagenavigation) CRDs.
- **Dynamic behavior** is provided by [KDexApp](../051_glossary.md#kdexapp) and [KDexScriptLibrary](../051_glossary.md#kdexscriptlibrary), wired using **ES modules published as NPM packages**.
- **Business Logic** is provided by [KDexFunction](../017_faas-integration-overview.md) CRDs attached to the host.


### Controller Topology

There are two main controllers: **[Nexus Manager](../051_glossary.md#nexus-manager)** and **[Host Manager](../051_glossary.md#host-manager)**.

1. **Nexus Manager (`kdex-nexus-manager`)**:
   - The cluster-level Kubernetes operator.
   - Watches all KDex CRDs across the cluster.
   - Responsible for deep validation, generation of derivative resources, and executing reconciliation loops that ensure the state of the cluster matches the declarative configuration.

2. **Host Manager (`kdex-host-manager`)**:
   - The web server and controller that physically serves the traffic.
   - Each `KDexHost` instance represents a virtual host that this controller routes requests to.
   - Handles [JWT authentication](../031_security-and-access-control.md#authentication-jwt-configuration), HTML assembly at the edge, and static asset mapping from OCI images.

### Host-Centric Model

The architecture is centered around the **`KDexHost`**:

- One `KDexHost` = one **web property** (e.g. `my.site.test`).
- One `KDexHost` = one **host manager** (i.e. compute is not shared between hosts).
- The host declares brand metadata, owning organization, and routing (domains). All host specifics are captured in the `spec`.
- Other CRDs attach to the host to define features and extensions.

### Shareable vs Host-Bound Components

To support operations at scale where multiple sites share components (such as unified corporate branding), KDex differentiates between **shareable** and **host-bound** components.

- **Shareable components (multi-site, reusable):**
  - `KDexApp`, `KDexClusterApp`
  - `KDexPageArchetype`, `KDexClusterPageArchetype`
  - `KDexPageFooter`, `KDexClusterPageFooter`
  - `KDexPageHeader`, `KDexClusterPageHeader`
  - `KDexPageNavigation`, `KDexClusterPageNavigation`
  - `KDexScriptLibrary`, `KDexClusterScriptLibrary`
  - `KDexTheme`, `KDexClusterTheme`

- **Host-bound components (per-site):**
  - `KDexHost` (The site core)
  - `KDexPage` (Pages bound to a specific host)
