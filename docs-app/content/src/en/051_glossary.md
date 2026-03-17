---
title: "GLOSSARY"
weight: 51
---


## Glossary of Terms

- <a id="archetype"></a>**Archetype**: A base HTML template layout containing `[[ .Content.<slot> ]]` placeholders. Used by `KDexPageArchetype`. See [Mark's Journey](030_user-journeys.md#journey-2--mark-define-theme-and-layout).
- <a id="faas-adaptor"></a>**FaaS Adaptor**: An infrastructure adapter capable of executing a backend function on Knative, AWS Lambda, OpenWhisk, etc. See [Function Lifecycle](040_function-lifecycle.md#deployment-faas-adaptors).
- <a id="host-manager"></a>**Host Manager (`kdex-host-manager`)**: The edge web server compiling layouts, handling traffic, rendering HTML, and resolving scripts on-the-fly. See [Architecture](01_Introduction/014_architecture.md#controller-topology).
- <a id="kdexapp"></a>**KDexApp**: Packages logic written as ES modules and Web Components built via NPM. See [Lucy's Journey](030_user-journeys.md#journey-4--lucy-build-and-deploy-a-kdexapp).
- <a id="kdexfunction"></a>**KDexFunction**: Scaffolds, builds, and deploys discrete backend business functions mapped to specific FaaS runtime adapters. See [Function Lifecycle](040_function-lifecycle.md).
- <a id="kdexhost"></a>**KDexHost**: The central configuration that dictates a site, mapping a web property to specific domains and acting as the root node for pages and themes. See [Host-Centric Model](01_Introduction/014_architecture.md#host-centric-model).
- <a id="kdexpage"></a>**KDexPage**: Maps a specific URL path fragment to structured HTML sections wrapped in an Archetype. See [Julie's Journey](030_user-journeys.md#journey-3--create-content-pages).
- <a id="kdexpagefooter"></a>**KDexPageFooter**: Defines the footer content (legal, links) once and shares it across pages.
- <a id="kdexpageheader"></a>**KDexPageHeader**: Defines the global header content (logo, primary actions).
- <a id="kdexpagenavigation"></a>**KDexPageNavigation**: Maintains menu structures available to be rendered by Archetypes.
- <a id="kdexscriptlibrary"></a>**KDexScriptLibrary**: Defines cross-cutting JavaScript dependencies and analytics to be injected into pages. See [Script Library Journey](030_user-journeys.md#journey-5--cross-cutting-javascript-with-kdexscriptlibrary).
- <a id="kdextheme"></a>**KDexTheme**: Encapsulates CSS stylesheets, external links, and edge-mapped asset distributions. See [Theme Journey](030_user-journeys.md#journey-2--mark-define-theme-and-layout).
- <a id="kdextranslation"></a>**KDexTranslation**: Provides localized labels and content for specified languages.
- <a id="nexus-manager"></a>**Nexus Manager (`kdex-nexus-manager`)**: The K-CNAS operator governing cluster resources. See [Architecture](01_Introduction/014_architecture.md#controller-topology).
