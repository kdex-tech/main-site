---
title: "KDex Function Lifecycle Architecture"
weight: 40
---


## KDexFunction Lifecycle Architecture

The **[KDexFunction](051_glossary.md#kdexfunction)** resource represents a concise unit of logic that scales in isolation. Its lifecycle is managed through a state machine that orchestrates code generation, building, and deployment.

### State Machine

The `KDexFunction` progresses through the following states:

1.  **Pending**: The initial state when a `KDexFunction` CR is created.
    - **Trigger**: Creation of the CR (e.g., via the Sniffer or manual application).
    - **Action**: The validation webhook or controller picks up the new resource.

2.  **OpenAPIValid**: The OpenAPI definition provided in `spec.API` has been validated.
    - **Trigger**: Successful validation of the OpenAPI spec against standards and internal policies.
    - **Implementation**: The **[Nexus Manager](051_glossary.md#nexus-manager)** (`kdex-nexus`) webhook validator uses the `vacuum` linter to validate the OpenAPI spec against the recommended OpenAPI 3.0.x ruleset.
    - **Action**: The system checks if the build configuration is valid.

3.  **BuildValid**: The build configuration (language, environment, dependencies) is valid.
    - **Trigger**: Verification that the `spec.Function` details (language, environment, etc.) are complete and supported.
    - **Action**: The system identifies the target FaaS Adaptor (e.g., based on config or default) and initiates the stub generation process.

4.  **SourceAvailable**: Source code stubs have been generated from the OpenAPI spec.
    - **Trigger**: Successful completion of the code generation pipeline.
    - **Dependency**: **FaaS Adaptor**. The generated stub code structure might vary depending on the target FaaS (e.g., `func.yaml` for Knative functions vs `lambda_handler` for AWS).
    - **Mechanism**: This is likely performed by a Kubernetes Job or Pod that:
        - Takes the OpenAPI spec, build details, and **Target Adaptor** as input.
        - Generates a project structure (e.g., Go module, Python package).
        - Pushes the generated code to an OCI registry or shared volume.
        - Updates `status.StubDetails` with the location of the generated code.

5.  **ExecutableAvailable**: An executable container image has been built and is available for deployment.
    - **Trigger**: Successful build of the container image.
    - **Dependency**: **FaaS Adaptor**. The container entrypoint and base image might differ based on the target runtime requirements.
    - **Mechanism**: Use of a container image build tool like **kpack** (Cloud Native Buildpacks) or Tekton.
        - Watches for changes in the generated source (or triggered by `SourceAvailable` state).
        - Builds the OCI image compliant with the target FaaS platform.
        - Updates `status.FunctionImage` (or similar) with the image reference.

6.  **FunctionDeployed**: The function has been deployed to the target FaaS runtime.
    - **Trigger**: Successful deployment by the FaaS Adaptor.
    - **Mechanism**: A FaaS Adaptor Operator:
        - Detects available FaaS adaptors (e.g., Knative, AWS Lambda, Azure Functions, OpenWhisk).
        - Selects the appropriate adaptor based on configuration.
        - Deploys the function container/artifact.
        - Updates the `KDexFunction` status with deployment details (e.g., internal address, platform-specific metadata).

7.  **Ready**: The function is verified and ready to serve traffic.
    - **Trigger**: Health checks pass on the deployed function.
    - **Mechanism**:
        - The controller verifies the deployed function is reachable and healthy.
        - Configures generic routing or ingress if necessary.
        - Updates `status.URL` with the reliable public/cluster-local endpoint.

### Architecture Components

#### Stub Generation
- **Strategy**: Kubernetes Jobs/Pods.
- **Workflow**: Controller creates a Job that mounts the necessary generators, orchestrates the generation, and outputs the artifacts.

#### Container Building
- **Strategy**: **kpack** (or comparable OCI builder).
- **Workflow**:
    - A `Image` or `Build` resource is created referencing the source provided by the Stub Generation phase.
    - kpack builds the image using appropriate buildpacks (e.g., Go, Python).
    - The resulting image digest is reported back.

#### Deployment (FaaS Adaptors)
- **Strategy**: **FaaS Adaptor Operator**.
- **Workflow**:
    - The `KDexFunction` reaches `ExecutableCreated`.
    - The FaaS Adaptor Operator evaluates the available adaptors and the function's requirements.
    - **Example (Knative)**:
        - Creates/Updates a Knative Service (KSVC) pointing to the `ExecutableCreated` image.
        - Updates status to `FunctionDeployed` with the KSVC address.
    - **Example (AWS Lambda)**:
        - Provisions a Lambda function using the container image.
        - Updates status to `FunctionDeployed` with the ARN and endpoint.

### Cross-Cutting Concerns

#### FaaS Adaptor Dependencies
The choice of FaaS Adaptor (e.g., Knative vs AWS Lambda) influences the entire pipeline, not just deployment:
- **Stub Generation**: Different platforms require different function signatures, handler files, or configuration headers (e.g., `project.toml` vs `func.yaml`).
- **Container Building**: The buildpacks or Dockerfiles used must produce an image compatible with the target runtime's invocation protocol (e.g., OCI HTTP entrypoint vs Lambda Runtime Interface Emulator).
Therefore, the **Target Adaptor** must be resolved early in the lifecycle (at `BuildValid` state) to guide subsequent generation and build steps.
