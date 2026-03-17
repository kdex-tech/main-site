---
title: "KDex Function / FaaS Integration"
weight: 17
---


## KDex Function / FaaS Integration

Beyond serving static and dynamic web experiences through [KDexApp](051_glossary.md#kdexapp) and frontend module assembly, K-CNAS offers an integrated pathway for serverless backend logic through the **[KDex Function](040_function-lifecycle.md)** ecosystem.

### Overview

While frontend developers build elements using ES modules and `KDexApp`, backend developers can focus purely on business logic. The [KDexFunction](051_glossary.md#kdexfunction) CRD and associated tooling provide a mechanism to describe backend requirements and APIs, which K-CNAS translates into executable containers and routable endpoints.

**Key capabilities include:**
- Representing concise units of logic as declarative [Kubernetes Custom Resources](050_api-crd-reference.md) (`KDexFunction`).
- Managing the full lifecycle from OpenAPI specification to deployed container.
- Supporting native deployment models and bridging to external Function-as-a-Service (FaaS) providers (e.g., Knative).

### Deep Dive

For advanced topics on how K-CNAS manages function state machines and performs automated code generation, check the dedicated **[KDex Function Lifecycle Architecture](040_function-lifecycle.md)** section later in this documentation.
