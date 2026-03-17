import { AppBridge } from '@kdex-tech/ui';

export class KdexDocsApp extends AppBridge(HTMLElement) {
  constructor() {
    super();
    this.attachShadow({ mode: 'open' });
    this._docs = {};
    this._ingressPath = '/-/a/app-docs';
  }

  async onRouteActivated(path) {
    path = path.substring(1);
    await this._loadManifest();
    const lang = document.documentElement.lang || 'en';
    const langRoot = this._docs[lang];

    if (!langRoot || !langRoot.children) {
      this._renderContent('<p>No documentation found for this language.</p>');
      return;
    }

    this._render(langRoot.children, path);

    if (path.length > 0) {
      await this._loadDoc(lang, path);
    } else {
      const first = this._findFirstLeaf(langRoot.children);
      if (first) {
        await this._loadDoc(lang, first.path);
      }
    }
  }

  _findFirstLeaf(nodes) {
    for (const node of nodes) {
      if (node.path) return node;
      if (node.children) {
        const first = this._findFirstLeaf(node.children);
        if (first) return first;
      }
    }
    return null;
  }

  async onRouteDeactivated() { }

  async _loadManifest() {
    try {
      const response = await fetch(`${this._ingressPath}/manifest.json`);
      if (!response.ok) throw new Error('Failed to load manifest');
      this._docs = await response.json();
    } catch (err) {
      console.error('Error loading manifest:', err);
      this._docs = {};
    }
  }

  async _loadDoc(lang, path) {
    try {
      const response = await fetch(`${this._ingressPath}/${lang}/${path}`);
      if (!response.ok) throw new Error(`Failed to load ${path}`);
      const text = await response.text();
      this._renderContent(text);

      this.shadowRoot.querySelectorAll('nav a').forEach((a) => {
        const isActive = a.dataset.path === path;
        a.classList.toggle('active', isActive);

        // Auto-expand parent sections if active
        if (isActive) {
          let parent = a.closest('.nav-section');
          while (parent) {
            parent.classList.add('open');
            parent = parent.parentElement.closest('.nav-section');
          }
        }
      });

      // Update active state for section headers
      this.shadowRoot.querySelectorAll('.nav-section-header').forEach((header) => {
        const isActive = header.dataset.path === path;
        header.classList.toggle('active', isActive);
        if (isActive) {
          header.parentElement.classList.add('open');
        }
      });

      this._highlightCode();
    } catch (err) {
      this._renderContent(`<p style="color: #f87171;">Error: ${err.message}</p>`);
    }
  }

  _render(nodes, activePath) {
    this.shadowRoot.innerHTML = `
      <style>
        @import url('https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css');
        :host {
          display: grid;
          grid-template-columns: 280px 1fr;
          gap: 2rem;
          min-height: 500px;
          color: #f1f5f9;
        }
        aside {
          border-right: 1px solid rgba(255, 255, 255, 0.05);
          padding-right: 1.5rem;
          height: fit-content;
          position: sticky;
          top: 2rem;
        }
        nav ul {
          list-style: none;
          padding: 0;
          margin: 0;
        }
        nav li {
          margin-bottom: 0.125rem;
        }
        .nav-section {
          margin-bottom: 0.5rem;
        }
        .nav-section-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0.5rem 0.75rem;
          cursor: pointer;
          font-size: 0.75rem;
          font-weight: 700;
          text-transform: uppercase;
          color: #64748b;
          letter-spacing: 0.05em;
          border-radius: 0.375rem;
          transition: background 0.2s, color 0.2s;
        }
        .nav-section-header:hover {
          background: rgba(255, 255, 255, 0.03);
          color: #94a3b8;
        }
        .nav-section-header.active {
          color: #3b82f6;
          background: rgba(59, 130, 246, 0.1);
        }
        .nav-section.open > .nav-section-header {
          color: #cbd5e1;
        }
        .nav-section.open > .nav-section-header.active {
          color: #3b82f6;
        }
        .nav-section-title {
          flex: 1;
        }
        .chevron {
          width: 12px;
          height: 12px;
          transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          opacity: 0.5;
        }
        .nav-section.open > .nav-section-header .chevron {
          transform: rotate(90deg);
          opacity: 1;
        }
        .nav-collapsible {
          display: grid;
          grid-template-rows: 0fr;
          transition: grid-template-rows 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          overflow: hidden;
        }
        .nav-section.open > .nav-collapsible {
          grid-template-rows: 1fr;
        }
        .nav-collapsible-inner {
          min-height: 0;
        }
        .nav-item {
          color: #94a3b8;
          text-decoration: none;
          font-size: 0.9rem;
          cursor: pointer;
          display: block;
          padding: 0.4rem 0.75rem;
          border-radius: 0.375rem;
          transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        }
        .nav-item:hover {
          color: #f1f5f9;
          background: rgba(255, 255, 255, 0.05);
        }
        .nav-item.active {
          color: #3b82f6;
          background: rgba(59, 130, 246, 0.1);
          font-weight: 600;
        }
        .nav-level-1 {
          padding-left: 1rem;
          border-left: 1px solid rgba(255, 255, 255, 0.05);
          margin-left: 0.75rem !important;
          margin-top: 0.25rem;
          margin-bottom: 0.5rem;
        }
        article {
          padding-bottom: 4rem;
          max-width: 800px;
        }
        article h1, article h2, article h3 {
          color: #f8fafc;
          margin-top: 1.5rem;
          font-weight: 700;
        }
        article h1 { font-size: 2.25rem; margin-top: 0; }
        article h2 { font-size: 1.5rem; border-bottom: 1px solid rgba(255, 255, 255, 0.1); padding-bottom: 0.5rem; }
        article p { margin: 1.25rem 0; line-height: 1.75; color: #cbd5e1; }
        article pre {
          background: #0f172a;
          border-radius: 0.75rem;
          overflow-x: auto;
          font-family: 'Fira Code', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
          border: 1px solid rgba(255, 255, 255, 0.1);
          margin: 1.5rem 0;
          box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
          line-height: 1.5;
        }
        article code {
          background: rgba(59, 130, 246, 0.1);
          color: #60a5fa;
          padding: 0.2rem 0.4rem;
          border-radius: 0.375rem;
          font-family: inherit;
          font-size: 0.875em;
        }
        article pre code {
          background: transparent;
          color: #f8fafc;
          padding: 0;
          font-size: 0.85rem;
          border-radius: 0;
        }
        article a {
          color: #3b82f6;
          text-decoration: none;
          border-bottom: 1px solid transparent;
          transition: border-color 0.2s;
        }
        article a:hover {
          border-bottom-color: #3b82f6;
        }
        article ul, article ol {
          margin: 1.25rem 0;
          padding-left: 1.5rem;
          color: #cbd5e1;
        }
        article li { margin-bottom: 0.5rem; }
        [id] {
          scroll-margin-top: 120px;
        }
      </style>
      <aside>
        <nav id="doc-nav">
          ${this._renderNav(nodes, 0, activePath)}
        </nav>
      </aside>
      <article id="content">
        <p>Loading documentation...</p>
      </article>
    `;

    // Handle link clicks
    this.shadowRoot.querySelectorAll('nav a').forEach((a) => {
      a.addEventListener('click', (e) => {
        const path = e.target.dataset.path;
        if (path) {
          this.navigate(path);
        }
      });
    });

    // Handle section toggles and navigation
    this.shadowRoot.querySelectorAll('.nav-section-header').forEach((header) => {
      header.addEventListener('click', () => {
        const path = header.dataset.path;
        if (path) {
          this.navigate(path);
        }
        const section = header.parentElement;
        section.classList.toggle('open');
      });
    });

    this._setupHighlighting();
  }

  _setupHighlighting() {
    if (!window.hljs) {
      if (!document.getElementById('hljs-script')) {
        const script = document.createElement('script');
        script.id = 'hljs-script';
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js';
        script.onload = () => {
          this._highlightCode();
        };
        document.head.appendChild(script);
      }
    } else {
      this._highlightCode();
    }
  }

  _highlightCode() {
    if (window.hljs) {
      this.shadowRoot.querySelectorAll('pre code').forEach((el) => {
        window.hljs.highlightElement(el);
      });
    }
  }

  _renderNav(nodes, level = 0, activePath = '') {
    if (!nodes || nodes.length === 0) return '';
    return `
      <ul class="nav-level-${level}">
        ${nodes
        .map((node) => {
          if (node.path && !node.children) {
            const isActive = activePath && node.path === activePath;
            return `<li><a data-path="${node.path}" class="nav-item ${isActive ? 'active' : ''}">${node.title}</a></li>`;
          } else if (node.children) {
            // Check if any child is active to pre-open the section
            const hasActiveChild = this._hasActiveChild(node.children, activePath);
            const isSelfActive = activePath && node.path === activePath;
            return `
              <li class="nav-section ${hasActiveChild || isSelfActive ? 'open' : ''}">
                <div class="nav-section-header ${isSelfActive ? 'active' : ''}" ${node.path ? `data-path="${node.path}"` : ''}>
                  <span class="nav-section-title">${node.title}</span>
                  <svg class="chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="9 18 15 12 9 6"></polyline>
                  </svg>
                </div>
                <div class="nav-collapsible">
                  <div class="nav-collapsible-inner">
                    ${this._renderNav(node.children, level + 1, activePath)}
                  </div>
                </div>
              </li>
            `;
          }
          return '';
        })
        .join('')}
      </ul>
    `;
  }

  _hasActiveChild(nodes, activePath) {
    if (!activePath) return false;
    for (const node of nodes) {
      if (node.path === activePath) return true;
      if (node.children && this._hasActiveChild(node.children, activePath)) return true;
    }
    return false;
  }

  _renderContent(content) {
    const article = this.shadowRoot.getElementById('content');
    if (article) {
      article.innerHTML = content;
    }
  }
}

customElements.define('kdex-docs-app', KdexDocsApp);
