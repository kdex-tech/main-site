---
title: "Bootstrapping a Host"
weight: 22
---


## Bootstrapping a Host

Because the K-CNAS Helm chart automatically deploys the bundled default custom resources (cluster-scoped themes, archetypes, and utility pages), bootstrapping a functional host instance is incredibly straightforward.

### The Single CRD Boot

To bring up a new site, you only need to apply a single **[KDexHost](051_glossary.md#kdexhost)** resource. This host will leverage the pre-installed cluster defaults to immediately provide a themed, routed property.

```yaml
apiVersion: kdex.dev/v1alpha1
kind: KDexHost
metadata:
  name: my-first-site
  namespace: default
spec:
  brandName: "My First Site"
  organization: "My Organization"
  routing:
    domains:
      - "my-site.test"
  # Utilizing default cluster components installed by Helm
  defaultThemeRef:
    kind: KDexClusterTheme
    name: default-theme
```

Apply this file to your cluster:

```bash
kubectl apply -f my-host.yaml
```

*Note: Depending on your environment, you may need to map `my-site.test` to your local ingress controller IP or LoadBalancer IP.*

### Next Steps

With your host running, you can begin adding **[KDexPage](051_glossary.md#kdexpage)** resources that reference `my-first-site` in their `hostRef`. These pages will automatically inherit the `default-theme` and the default page archetypes bundled with K-CNAS!
