# Content Kit

A CLI tool which to edit contents of Shopify's store through Admin API that inspired by Theme Kit.

## Usage

1. Download latest binary from [Release page](https://github.com/kogai/k9bookshelf/releases).
2. Add configuration file like below
3. Execute `./path/to/bin deploy|download --config path/to/config.yml` and take a look your result!

```yml
# ./path/to/config.yml
content:
  domain: your-shop.myshopify.com
  key: ${MARKDOWN_APP_KEY}
  secret: ${MARKDOWN_APP_SECRET}
  token: ${MARKDOWN_APP_SECRET}
  dir: ${PWD}/contents
```

## Blog articles

Markdown files under `contents/blogs/<blog-handle>/<article-handle>.md` are deployed as blog articles.

A file may start with YAML front matter.

```md
---
title: "Article title"
published_at: "2016-12-22T19:31:47+09:00" # RFC3339
tags:
  - TypeScript
---

Body of the article...
```

- If the article already exists on the store, its body is updated, along with any fields given in front matter.
- If the article doesn't exist, it's created. Creating one requires `title`.
- If the blog doesn't exist, it's created with its handle as the title.
- `download` writes `title`, `published_at` and `tags` of blog articles as front matter.
