# Build the theme image

```shell
make docker-buildx

# to local registry
make docker-buildx PLATFORMS=linux/amd64 REPOSITORY=k3d-registry:5000
```
