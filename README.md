**Deprecation Notice**: This project will be archived by 2026-02.

# productctl

A power tool for [Red Hat Software Certification
Partners](https://connect.redhat.com/en/benefits).

Use `productctl` to create and iterate on your Red Hat Certified Product listing
metadata from the comfort of your favorite text editor. Then, apply your changes
and validate your content from your Partner dashboard.

> [!CAUTION]
> 
> - This utility is in **very early development**.
> - Users will assume all risk of using this tooling, including overwriting product metadata.
> - This tool is not officially supported by Red Hat Support.
> - Every API used here is subject to change.

## Configuration

Create a configuration file at any of these locations (listed in order of precedence):

- $PWD/.productctl/config.yaml
- $XDG_CONFIG_DIR/productctl/config.yaml
- $HOME/.productctl/config.yaml

Example contents:

```yaml
# env: PRODUCTCTL_API_TOKEN
api-token: your-api-token

# env: PRODUCTCTL_LOG_LEVEL
log-level: "info"
```

Alternatively, you can set the environment variables mentioned in-line.

## Usage

```
Manage your Product Listing

Usage:
  productctl product [command]

Available Commands:
  apply       Apply changes to Partner product listings from the input file.
  cleanup     Detaches and archives components. Deletes the product listing. This is destructive. Use with caution.
  create      Start building a new product listing declaration on your filesystem
  fetch       Get a pre-existing product listing
  jsonschema  Generate resource jsonschema for LSPs that support it.
  sanitize    Cleans declaration for re-use and emits to stdout

Flags:
      --custom-endpoint string   Define a custom API endpoint. Supersedes predefined environment values like "prod" if set
      --env string               The catalog API environment to use. Choose from stage, prod (default "prod")
  -h, --help                     help for product

Global Flags:
      --log-level string   The verbosity of the tool itself. Ex. error, warn, info, debug (default "info")

Use "productctl product [command] --help" for more information about a command.
```

## High Level Workflow

1. Scaffold your new Product Listing (or fetch an existing one)

```bash
productctl product create [--from-discovery-json /path/to/discovery.json] my.product.yaml
```

Or fetch an existing listing:

```bash
productctl product fetch 000111222333 > my.product.yaml
```

2. Make alterations to your Product Listing, add/remove components, etc.

3. Apply your Product Listing

```bash
productctl product apply my.product.yaml
```

4. Repeat until all metadata is configured to your liking.

## Getting Started

See our [Getting Started](docs/GETTING_STARTED.md) guide.

