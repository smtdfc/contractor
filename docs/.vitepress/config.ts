import { defineConfig } from "vitepress";

export default defineConfig({
  title: "Contractor",
  description: "Type-Safe IDL & Code Generation Toolchain",
  base: "/",
  themeConfig: {
    nav: [
      { text: "Guide", link: "/guide/getting-started" },
      { text: "GitHub", link: "https://github.com/smtdfc/contractor" },
    ],
    sidebar: [
      {
        text: "Introduction",
        items: [
          { text: "Getting Started", link: "/guide/getting-started" },
        ],
      },
      {
        text: "Fundamentals",
        collapsed: true,
        items: [
          { text: "Models", link: "/guide/language-basics/models" },
          { text: "Enums", link: "/guide/language-basics/enums" },
          { text: "Errors", link: "/guide/language-basics/errors" },
          { text: "Events", link: "/guide/language-basics/events" },
          { text: "REST Endpoints", link: "/guide/language-basics/rest" },
        ],
      },
      {
        text: "Advanced",
        collapsed: true,
        items: [
          { text: "Validation", link: "/guide/validation" },
          { text: "Code Generation", link: "/guide/code-generation" },
        ],
      },
      {
        text: "Typescript",
        collapsed: true,
        items: [
          { text: "Code Generation", link: "/guide/typescript/code-generation" },
        ],
      },
    ],
    socialLinks: [
      { icon: "github", link: "https://github.com/smtdfc/contractor" },
    ],
  },
});
