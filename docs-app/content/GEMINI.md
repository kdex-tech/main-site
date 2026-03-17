- To rebuild the image for testing in the local Kubernetes cluster use the command: `make docker-buildx PLATFORMS=linux/amd64 REPOSITORY=k3d-registry:5000`

- To ensure the latest image is used in the local Kubernetes cluster after building use the command: `kubectl -n kdex-site delete pod -l kdex.dev/backend=app-docs`