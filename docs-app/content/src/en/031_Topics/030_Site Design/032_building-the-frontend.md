---
title: "Building the Frontend"
weight: 32
---


## Building the Frontend

When composing sites in K-CNAS using structural components like Themes, Archetypes, Script Libraries, and Navigations, you'll encounter two different paradigms: **Cluster-Scoped Components** and **Namespace-Scoped Components**.

Understanding the difference helps organize components as reusable models vs specific isolated instances.

### Cluster vs Namespace Scoped Components

#### Namespace-Scoped CRDs
Examples: `KDexTheme`, `KDexPageArchetype`, `KDexPageHeader`, `KDexPageFooter`, `KDexPageNavigation`, `KDexScriptLibrary`

- Reside in an isolated namespace alongside the `KDexHost`.
- Excellent for single-site bespoke styling and configurations.
- Recommended when components shouldn't be accessible across other hosts out of the site boundary.

#### Cluster-Scoped CRDs
Examples: `KDexClusterTheme`, `KDexClusterPageArchetype`, `KDexClusterPageHeader`, `KDexClusterPageFooter`, `KDexClusterPageNavigation`, `KDexClusterScriptLibrary`

- Installed globally into the cluster and available for consumption by **any** `KDexHost` or `KDexPage` in any namespace.
- Function as the foundational libraries for the organization (e.g., standard corporate theme, standard documentation layout).
- A perfect model for reusable components across isolated multi-tenant environments.

### Component Breakdown

When building the structure of a site, you generally combine:

#### 1. Archetypes (`KDexPageArchetype` / `KDexClusterPageArchetype`)
Defines the `<html>` shell, `<head>` elements, and main content regions. Allows you to utilize placeholder fields like `[[ .Content.main ]]` and `[[ .Navigation.sidebar ]]`.

#### 2. Theme (`KDexTheme` / `KDexClusterTheme`)
Supplies the style sheets (`<style>` or `<link href="...">`), webfonts, and any static assets wrapped inside an OCI container for edge mapping. Attaches directly to `KDexHost` as a default or to specific `KDexPage`s.

#### 3. Headers & Footers (`KDexPageHeader` / `KDexClusterPageHeader`, `KDexPageFooter` / `KDexClusterPageFooter`)
Structural boilerplate providing logo, static menus, or rigid components loaded on every page utilizing the archetype's header/footer slots.

#### 4. Navigations (`KDexPageNavigation` / `KDexClusterPageNavigation`)
Defines hierarchical menus, routing configurations, and cross-links that can be rendered through navigation blocks embedded natively in KDex's runtime.

#### 5. Script Libraries (`KDexScriptLibrary` / `KDexClusterScriptLibrary`)
Injects tracking algorithms, functional JS overlays, external integrations, or foundational web components packages automatically into the specific pages or the encompassing host.

By intelligently segregating styles and capabilities between **Cluster-Scoped** models and site-specific overrides via **Namespace-Scoped** instances, large organizations can effectively govern branding while empowering developers.
