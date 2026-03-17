- To rebuild the app for testing in the local Kubernetes cluster use the command: `cd kdex-main-site && make build-app`

- To ensure the app is used in the local Kubernetes cluster after building use the command: `kubectl -n kdex-site apply -f kdex-main-site/k8s/100_app_docs.yaml` you may need to wait up to 10 seconds for the site to reconcile the update.