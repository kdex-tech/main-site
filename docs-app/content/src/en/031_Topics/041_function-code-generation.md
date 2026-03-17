---
title: "Function Code Generation"
weight: 41
---


## Function Code Generation

As part of the `KDexFunction` lifecycle, K-CNAS automatically generates the necessary boilerplate, stubs, and models directly from an OpenAPI specification embedded within or referenced by the CRD.

This guarantees that the logic developers write is statically typed against the declared APIs, bridging the gap between design and implementation safely.

### How It Works

1. **Extraction**: Once the OpenAPI specification is validated (`OpenAPIValid` state), the generation pipeline parses the definition.
2. **Scaffolding Pipeline**: 
   Depending on the configured language and target environment (e.g., Go, TypeScript), K-CNAS utilizes a specialized container (similar to `ghcr.io/kdex-tech/cli-tools:latest`) inside the cluster to scaffold the source directory structure.
3. **Artifact Production**:
   This automated pipeline produces:
   - Request and response models.
   - Endpoint handlers (empty stubs calling an abstract service layer).
   - Server bootstrapping logic tailored for the chosen Adaptor (HTTP wrapper, Knative binding, Lambda emulator wrapper).
4. **Final Packaging**: The compiled logic is then staged into an OCI container in the subsequent `ExecutableAvailable` state.
