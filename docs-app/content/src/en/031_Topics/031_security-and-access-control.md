---
title: "Authentication & Access Control"
weight: 31
---


## Authentication & Access Control

Securing your KDex sites involves configuring protocol-level security (TLS/Cert Manager), setting up Authentication (local / OAuth2 / OIDC) secured via JWTs, and establishing page authorization.


### Protocol-Level Security

For production use, securing the transport layer with TLS is a hard requirement. K-CNAS is built to smoothly integrate with standard Kubernetes tools like `cert-manager`.

Additionally, many Cloud Vendor services can abstract TLS entirely by terminating HTTPS at a layer 7 load balancer before traffic reaches your cluster. In those architectures, Cert-Manager may not be necessary.


#### Configuring Cert-Manager 

If you choose to use `cert-manager` inside your cluster for TLS provisioning (e.g. using Let's Encrypt), simply deploy it via:

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.19.2/cert-manager.yaml
```

Make sure to specify a `tls.secretRef` inside your **[KDexHost](051_glossary.md#kdexhost)** routing spec that corresponds to the generated cert-manager secret.


### Authentication: JWT Configuration

The KDex Web server uses **RS256 (RSA signatures)** for JWT tokens, establishing standards-compliance suitable for distributed microservices.


#### Key Features
✅ **Asymmetric Signing**: Private key for signing, public key for verification  
✅ **JWKS Endpoint**: Exposes public keys at `/.well-known/jwks.json`  
✅ **Zero Shared Secrets**: Downstream services only need the JWKS URL  
✅ **Standards Compliant**: Compatible with any OAuth2/OIDC library  


#### Setup
Ensure that the `kdex-jwt-keys` Kubernetes Secret (containing `private-key` and `public-key` in PEM format) is provisioned inside the namespace in which K-CNAS operator runs. The `kdex-host` process exposes the JWKS endpoint upon boot to allow external validations.

Example JWKS Output (`https://your-kdex-host-host/.well-known/jwks.json`):
```json
{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "alg": "RS256",
      "kid": "kdex-auth-key-1",
      "n": "base64-encoded-modulus",
      "e": "AQAB"
    }
  ]
}
```


### Page Authorization

KDex Web server supports fine-grained authorization for pages based on security requirements defined in **[KDexPage](051_glossary.md#kdexpage)** and **[KDexHost](051_glossary.md#kdexhost)** CRDs.


#### Checking Security Requirements

Security requirements can be defined at two levels:
- **Page-level**: In `KDexPage.spec.security`
- **Host-level**: In `KDexHost.spec.security` (fallback if page has no security)


#### Authorization Flow

1. **No security requirements** → Allow access (public page)
2. **Has security requirements but no JWT** → Redirect to login page
3. **Has JWT** → Validate claims:
   - Scopes are checked against the requested route.
   - Requirements are evaluated using an **OR** mechanism between blocks, and **AND** inside the specific block.

<details>
<summary>View Example YAML</summary>

```yaml
apiVersion: kdex.dev/v1alpha1
kind: KDexPage
metadata:
  name: admin-dashboard
spec:
  paths:
    basePath: "/admin"
  security:
  - bearer:  
    - "pages:admin:read"
  contentEntries:
  - slot: main
    rawHTML: "<h1>Admin Dashboard</h1>"
```

</details>


#### Scope Formats

Scopes follow the pattern: `resource:resourceName:verb`

> **Note**: If `resourceName` contains colons (`:`), it **must** be URL-encoded (using `url.PathEscape`) to prevent misinterpretation.

- `pages::read` - Read access to all pages
- `pages:home:read` - Read access to specific page "home"
- `pages:foo%3Abar:read` - Read access to page "foo:bar" (URL-encoded)


#### User Scopes & Roles

Users gain scopes through `KDexRoleBinding` linking to an Identity/Session referencing a `KDexRole`. The `KDexRole` details the resources and verbs a persona holds.
