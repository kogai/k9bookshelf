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
