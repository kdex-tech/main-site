---
title: "Repository Functional Map"
weight: 15
---


## Repository Functional Map

This document provides a comprehensive map of the KDex platform, identifying the functional responsibilities of each repository and the relationships between architectural components.

### Core Architectural Components

#### 1. API & Schema Definition (`kdex-crds`)
- **Responsibility:** Defines the source of truth for the KDex platform via Kubernetes Custom Resource Definitions (CRDs).
- **Group:** `kdex.dev/v1alpha1`
- **Key Resources:** `KDexHost`, `KDexPage`, `KDexApp`, `KDexTheme`, `KDexFunction`, `KDexScriptLibrary`.

#### 2. Core Control Plane (`kdex-nexus-manager`)
- **Responsibility:** Implements the primary controllers for the KDex resource lifecycle.
- **Key Features:** Ingress/Gateway management, cross-cutting script assembly, resource resolution, and RBAC.

#### 3. Host & Internal Management (`kdex-host-manager`)
- **Responsibility:** Handles host-level operations, internal resource mapping (`KDexInternal*`), and authentication (OIDC/LDAP).
- **Key Features:** Token management, internal package reference resolution, and FaaS integration logic.

#### 4. FaaS Workflow (`kdex-fngogen`, `kdex-knative-deployer`)
- **kdex-fngogen:** Provides code generation tools to transform OpenAPI specs into Go implementation code. It is an example of a code generator implementation managed by a `KDexFaaSAdaptor`. Other FaaS adaptors or languages may require different code generator implementations.
- **kdex-knative-deployer:** Orchestrates the deployment of `KDexFunction` resources as Knative services.

#### 5. Frontend Runtime (`kdex-ui`)
- **Responsibility:** Provides the browser-side runtime, routing, and metadata management for KDexApp Web Components.

---

### Utility Libraries

- **kdex-dmapper:** Data mapping utility using Common Expression Language (CEL).
- **kdex-entitlements:** User and role entitlement management library.

---

### Tooling & Infrastructure

- **kdex-cli-tools:** Alpine-based image containing common CLI tools (`git`, `kubectl`, `jq`, `curl`, etc.) for syncing generated code to the repository and debugging.
- **kdex-node-tools:** Node.js-based tool image.
- **kdex-backend-static:** Caddy-based static resource server.

---

### Site & Documentation

- **kdex-main-site:** Defines the KDex company site, including the documentation application (`docs-app`), visual theme (`kdex-main-theme`), and production Kubernetes manifests.

---

### Functional Relationship Summary

| Resource | Primary Controller | Key Collaborators |
| :--- | :--- | :--- |
| `KDexHost` | `kdex-nexus-manager` | `kdex-host-manager` (Auth/Internal mapping) |
| `KDexPage` | `kdex-nexus-manager` | `kdex-ui` (Routing) |
| `KDexApp` | `kdex-nexus-manager` | `kdex-ui` (Runtime) |
| `KDexFunction` | `kdex-host-manager` | `kdex-knative-deployer`, `kdex-fngogen` |
| `KDexTheme` | `kdex-nexus-manager` | OCI Registries (via `kdex-host-manager` logic) |
