import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'go-vibe',
  description: 'Production-ready Go microservice template',
  base: '/go-vibe/',

  appearance: 'dark',

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/go-vibe/favicon.svg' }],
    ['meta', { name: 'theme-color', content: '#00ACD7' }],
  ],

  themeConfig: {
    logo: { light: '/logo-light.svg', dark: '/logo-dark.svg', alt: 'go-vibe' },

    nav: [
      { text: 'Home', link: '/' },
      {
        text: 'Guide',
        activeMatch: '/guide/',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'Architecture', link: '/guide/architecture' },
          { text: 'Features', link: '/guide/features' },
          { text: 'API Reference', link: '/guide/api' },
          { text: 'Deployment', link: '/guide/deployment' },
          { text: 'Observability', link: '/guide/observability' },
          { text: 'CI/CD', link: '/guide/ci-cd' },
        ],
      },
      { text: 'API Reference', link: '/guide/api' },
      { text: 'Deployment', link: '/guide/deployment' },
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          collapsed: false,
          items: [
            { text: 'Introduction', link: '/guide/getting-started' },
            { text: 'Architecture', link: '/guide/architecture' },
          ],
        },
        {
          text: 'Features',
          collapsed: false,
          items: [
            { text: 'Feature Overview', link: '/guide/features' },
            { text: 'API Reference', link: '/guide/api' },
          ],
        },
        {
          text: 'Deployment',
          collapsed: false,
          items: [
            { text: 'Docker & Kubernetes', link: '/guide/deployment' },
          ],
        },
        {
          text: 'Observability',
          collapsed: false,
          items: [
            { text: 'Metrics & Logging', link: '/guide/observability' },
          ],
        },
        {
          text: 'CI/CD',
          collapsed: false,
          items: [
            { text: 'GitHub Actions', link: '/guide/ci-cd' },
          ],
        },
      ],
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/huberp/go-vibe' },
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024 go-vibe contributors',
    },

    editLink: {
      pattern: 'https://github.com/huberp/go-vibe/edit/main/docs-site/:path',
      text: 'Edit this page on GitHub',
    },

    search: {
      provider: 'local',
    },
  },

  markdown: {
    theme: {
      light: 'github-light',
      dark: 'one-dark-pro',
    },
    lineNumbers: true,
  },
})
