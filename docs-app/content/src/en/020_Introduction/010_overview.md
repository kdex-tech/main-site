---
title: "What is K-CNAS?"
weight: 10
---


## What is K-CNAS?

**K-CNAS** is an **Internal Developer Platform (IDP)** that is positioned as the _first of its kind_ **Cloud Native Application Server**.

Instead of manually deploying and wiring web servers, reverse proxies, and application runtimes, K-CNAS lets you **model digital experiences as [Kubernetes Custom Resources (CRDs)](../050_api-crd-reference.md)**. The KDex operator reconciles those CRDs into ingress rules, HTML composition, theme and asset delivery, JavaScript module assembly and individually scalable business logic.

With KDex, you can bring up new sites with:
- **No containers to build** for the frontend serving
- **No deployments to manage** manually
- Fully **managed JavaScript module assembly** (handling security, caching, delivery)
- Fully **managed business logic** (handling code generation, building, and scaling)
- Fully **managed ingress and routing** (driven by declarative configuration)

---

## High-Level Architecture

### Host-Centric Model

The architecture is centered around the **`KDexHost`**:

- One `KDexHost` = one **web property** (e.g. `my.site.test`).
- The host declares brand metadata, owning organization, and routing (domains). All host specifics are captured in the `spec`.
- Other CRDs attach to the host to define:
  - Pages (`KDexPage`)
  - Themes and assets (`KDexTheme`)
  - Page structure (`KDexPageArchetype`, `KDexPageHeader`, `KDexPageFooter`, and `KDexPageNavigation`)
  - Dynamic applications (`KDexApp`)
  - Shared JavaScript modules (`KDexScriptLibrary`)
  - Localized content (`KDexTranslation`)

### Component Categories

To support operations at scale where multiple sites share components (such as unified corporate branding), KDex differentiates between **shareable** and **host-bound** components.

- **Shareable components (multi-site):**
  - `KDexApp`
  - `KDexPageArchetype`
  - `KDexPageFooter`
  - `KDexPageHeader`
  - `KDexPageNavigation`
  - `KDexScriptLibrary`
  - `KDexTheme`
  - `KDexTranslation`
  - `KDexUtilityPage`

- **Host-bound components (per-site):**
  - `KDexHost`
  - `KDexPage`
