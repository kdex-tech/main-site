---
title: "User Journeys"
weight: 30
---


## User Journeys

This document illustrates how different personas use KDex CRDs to collaboratively build and operate sites.

Personas:
- **Chey (Operations Technician)** – owns infrastructure, routing, certificates, resource allocation and operational posture.
- **Mark (Designer)** – owns brand, visual system and page layouts.
- **Lucy (Application Developer)** – owns dynamic applications, data services and integrations.
- **Julie (Content Manager)** – owns pages, navigation and textual content.
- **Anna (Site Manager)** – owns the site, authentication and authorization.


## Journey 1 – Chey: Provision a New Site

**Goal:** [Bring a new site online](021_bootstrapping-host.md) at `new-site.acme.com` with the default corporate theme.

**Permission:** Chey has permission to create namespaces and `KDexHost` instances and to configure host DNS and TLS. He can also grant permission to edit the host once it is created.

### Steps

1. **Create a namespace** for the site (optional but recommended):

    ```yaml
    apiVersion: v1
    kind: Namespace
    metadata:
      name: new-site
    ```

2. **Create the [KDexHost](051_glossary.md#kdexhost):**

    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexHost
    metadata:
      name: new-site
      namespace: new-site
    spec:
      brandName: "New Site"
      organization: "Acme Inc."
      routing:
        domains:
          - new-site.acme.com
      themeRef:
        kind: KDexClusterTheme
        name: acme-primary-theme
    ```

3. **Gather the IP from the KDexHost status and Configure DNS & TLS** so `new-site.acme.com` routes to the site's managed ingress endpoint.

4. The site is ready to be **handed off to Anna and Julie**.

**Result:** A new KDex-powered web property exists and is reachable at `https://new-site.acme.com`. As Chey evolves the components under his control changes are rolled out seamlessly and without need for other players to be involved.

**Use of [ExternalDNS](https://kubernetes-sigs.github.io/external-dns/latest/) and [CertManager](https://cert-manager.io/)** are both supported and encouraged. With KDex, these greatly improve and streamline automation.

**[Gateway API](https://kubernetes.io/docs/concepts/services-networking/gateway/)** is also supported here if that is preferred over [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/).


## Journey 2 – Mark: Define reusable Theme and Layout

**Goal:** Create reusable theme and structural page components that support cohesive branding for the entire oranization having a distinct operational lifecycle eliminating friction with other teams.

**Permission:** Mark has permission to create and update `KDexClusterThemes`, `KDexClusterPageArchetypes` and other cluster scoped KDex structural components. He can also create and update `KDexClusterScriptLibraries` that support design and structure components such as delivering vetted frontend component frameworks and design libraries.

### Steps

1. **Assemble static resources** into an OCI image using docker:

    ```text
    FROM scratch
    COPY --chown=65532:65532 . .
    ```

    _**Note:** The contents of this image has no formal structure other. It is a collection of static resources; CSS files, logos, icons, fonts that encapsulate the brand. It contains no logical elements._

2. **Push the static image** to a trusted registry:

    ```shell
    docker build -t registry.acme.com:5000/acme/primary-theme:1.2.3 .
    docker push registry.acme.com:5000/acme/primary-theme:1.2.3
    ```

3. **Design the [KDexClusterTheme](051_glossary.md#kdextheme)** and deploy it:

    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexClusterTheme
    metadata:
      name: acme-primary-theme
    spec:
      assets:
        - linkHref: "https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap"
          attributes:
            rel: stylesheet
        - linkHref: "/-/theme/css/main.css"
          attributes:
            rel: stylesheet
      staticImage: registry.acme.com:5000/acme/primary-theme:1.2.3
    ```

4. **Create a [KDexClusterPageArchetype](051_glossary.md#archetype)** that defines the base HTML page layout and named sections:

    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexClusterPageArchetype
    metadata:
      name: primary-page-archetype
    spec:
      content: |
        <!DOCTYPE html>
        <html lang="[[ .Language ]]">
          <head>
            [[ .Meta ]]
            <title>[[ .Title ]]</title>
            [[ .Theme ]]
          </head>
          <body>
            <header>
              [[ .Header ]]
            </header>
            <nav>
              [[ .Navigation.main ]]
            </nav>
            <main>
              [[ .Content.main ]]
            </main>
            [[- if .Content.sidebar ]]
            <aside>
              [[ .Content.sidebar ]]
            </aside>
            [[- end ]]
            <footer>
              [[ .Footer ]]
            </footer>
          </body>
        </html>
    ```

5. Optionally create:
    - `KDexClusterPageHeader` for logo and primary nav.
    - `KDexClusterPageFooter` for legal text and links.
    - `KDexClusterPageNavigation` for one or more navigation menus.

**Result:** Mark has delivered a reusable visual and structural system that other personas can compose without changing infrastructure. As Mark evolves the components under his control changes are rolled out seamlessly and without need for other players to be involved.


## Journey 3 – Lucy: Build and Deploy a Documenation App

**Goal:** Implement a dynamic `documentation` application as a Web Component, publish it as an NPM package, and make it available on the cluster.

**Permission:** Lucy has permission to create and update `KDexClusterApps`.

### Steps

1. **Implement the Web Component** in a project that emits an ES module:

    ```json
    {
      "name": "@acmi/docs-app",
      "version": "1.0.0",
      "type": "module",
      "module": "dist/docs-app.js",
      ...
    }
    ```

2. **Export a custom element** from the module, e.g. `acmi-docs-main`.

3. **Publish the NPM package** to a registry accessible by the KDex controllers.

4. **Create the [KDexClusterApp](051_glossary.md#kdexapp) resource:**

    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexClusterApp
    metadata:
      name: acmi-docs-app
    spec:
      packageReference:
        name: "@acmi/docs-app"
        version: 1.0.0
      customElements:
      - name: acmi-docs-main
    ```

5. **Coordinate with Julie** so she references `acmi-docs-app` from a **[KDexPage](051_glossary.md#kdexpage)**.

**Result:** Lucy ships a reusable, versioned application that can be wired into any host and page via CRDs, without changing web server configuration.


## Journey 4 – Julie: Create Content Pages

**Goal:** Define site pages (home, docs, etc.) for `new-site.acme.com` using static and dynamic content.

**Permission:** Julie has permission to create and update `KDexPages`, `KDexPageNavigations`,`KDexTranslations` and other host centric KDex components in the `new-site` namespace.

### Steps

1. Create the `Home` page from static content:

    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexPage
    metadata:
      name: new-site-home
      namespace: new-site
    spec:
      basePath: "/"
      contentEntries:
      - slot: main
        rawHTML: |
          <h1>[[ l10n "welcome-to-x" .BrandName ]]</h1>
      hostRef:
        name: new-site
      label: Home
      overrideNavigationRefs:
        main:
          kind: KDexPageNavigation
          name: new-site-navigation-main
        footer:
          kind: KDexPageNavigation
          name: new-site-navigation-footer
      pageArchetypeRef:
        kind: KDexClusterPageArchetype
        name: primary-page-archetype
    ```

2. Create the `Docs` page from static and dynamic content:

    Lucy created a **[KDexClusterApp](051_glossary.md#kdexapp)** called `acmi-docs-app` that Julie will reference on the `Docs` page.

    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexPage
    metadata:
      name: new-site-docs
      namespace: new-site
    spec:
      basePath: "/docs"
      contentEntries:
      - slot: sidebar
        rawHTML: |
          <h1>[[ l10n "documentation" ]]</h1>
      - slot: main
        appRef:
          kind: KDexClusterApp
          name: acmi-docs-app
        attributes:
          dataPath: /v2/docs
        customElementName: acmi-docs-main
      hostRef:
        name: new-site
      overrideNavigationRefs:
        main:
          kind: KDexPageNavigation
          name: new-site-navigation-main
        footer:
          kind: KDexPageNavigation
          name: new-site-navigation-footer
      pageArchetypeRef:
        kind: KDexClusterPageArchetype
        name: primary-page-archetype
    ```

3. Julie defines **[KDexPageNavigation](051_glossary.md#kdexpagenavigation)** that implement the menus (e.g. `main`, `footer`) that she referenced on the pages.
    
4. She adds **[KDexTranslation](051_glossary.md#kdextranslation)** resources to provide localized values for the keys she used in her content templates (e.g. `documentation` and `welcome-to-x`).

**Result:** Julie created the site pages by editing CRDs, without touching deployments or load balancers.

## Journey 5 – Anna applies Cross-Cutting JavaScript with [KDexScriptLibrary](051_glossary.md#kdexscriptlibrary) and configures Authentication

**Goal:** Anna needs to add analytics across all the pages of the host she manages. She also needs to enable authentication through the enterprise OIDC Identity Provider.

### Steps

1. Anna **defines a [KDexScriptLibrary](051_glossary.md#kdexscriptlibrary)** that pulls in analytics library:

    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexScriptLibrary
    metadata:
      name: global-analytics
      namespace: new-site
    spec:
      scripts:
        - src: "https://cdn.example.com/analytics.js"
    ```

2. Chey assists Anna by providing here an OIDC providerURL and the name of a Secret that holds the OIDC client configuration for `new-site.acme.com`:
    - provider URL  - `https://oidc.acme.com`
    - secret - `new-site-oidc`

3. Anna creates a ServiceAccount referencing the secret:

    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: new-site-service-account
    secrets:
    - name: new-site-oidc
    ```

4. Anna adjusts the KDexHost so that:
    1. OIDC is enabled using the povider URL Chey provided
    2. The script library is attached to the host making it available to every page
    3. The service account is referenced by the host and thus can locate and read the OIDC secret to get client configuration


    ```yaml
    apiVersion: kdex.dev/v1alpha1
    kind: KDexHost
    metadata:
      name: new-site
      namespace: new-site
    spec:
      auth:                                       # 1.
        oidcProvider:                             # 1.
          oidcProviderURL: https://oidc.acme.com  # 1.
      brandName: "New Site"
      organization: "Acme Inc."
      routing:
        domains:
          - new-site.acme.com
      scriptLibraryRef:                           # 2.
        kind: KDexScriptLibrary                   # 2.
        name: global-analytics                    # 2.
      serviceAccountRef:                          # 3.
        name: new-site-service-account            # 3.
      themeRef:
        kind: KDexClusterTheme
        name: acme-primary-theme
    ```

4. The **KDex host manager assembles all package references and generates import maps** when any host resources are reconciled.

**Result:** Anna has applied several cross-cutting configurations to the host and its pages without risk of exposing sensitive information or blocking any one elses efforts.


## Journey 6 – Operators and Security

**Goal:** Acme Inc. decides to evolve from their simplistic one site setup to a hardened, [multi-tenant posture](../031_security-and-access-control.md).

### Steps

1. **Start in a permissive mode** (e.g. perhaps with `devMode: true` and `modulePolicy: Loose`) while prototyping.
2. **Tighten the `modulePolicy`** to `ExternalDependencies`, `ModulesRequired` or `Strict` as security expectations increase.
3. **Introduce controlled registries and CDNs** for script, NPM dependencies and container images.
4. **Refine RBAC** so each persona only has access to the CRDs and namespaces they need.
5. **Adopt GitOps and policy-as-code** (e.g. OPA, Kyverno) to validate KDex CRs before they reach the cluster.

**Result:** The same declarative model that accelerated prototyping can be incrementally hardened into a production-ready, secure Cloud Native Application Server platform.
