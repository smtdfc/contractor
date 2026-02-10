import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Contractor",
  description: "Contractor doc page",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Documentation", link: "/docs/getting-started" },
    ],
    sidebar: {
      "/docs/": [
        {
          text: "Hướng dẫn bắt đầu",
          collapsed: false,
          items: [
            { text: "Bắt đầu", link: "/docs/getting-started" },
            { text: "Cài đặt", link: "/docs/installation" },
            { text: "Hướng dẫn cấu hình", link: "/docs/configuration" },
            { text: "Biên dịch", link: "/docs/generate" },
          ],
        },
        {
          text: "Model",
          collapsed: false,
          items: [
            { text: "Tạo model", link: "/docs/model/create" },
            { text: "Annotations", link: "/docs/model/annotations" },
            {
              text: "Field Annotations",
              link: "/docs/model/field-annotations",
            },
          ],
        },
      ],
    },

    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2026-present smtdfc",
    },
    socialLinks: [
      { icon: "github", link: "https://github.com/smtdfc/contractor" },
    ],
  },
});
