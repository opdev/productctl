# productctl

A power tool for [Red Hat Software Certification
Partners](https://connect.redhat.com/en/benefits).

The `productctl` utility allows Red Hat Software Certification Partners to
easily iterate on the definition of their product listings in a familiar way.

This tooling is intended for Red Hat Software Certification Partners who are
working through the Software Certification Process for Containers, Helm Charts,
and Operators. It will allow you to populate much of the metadata about your
Certification Components and Product Listings from your workstation.

> [!CAUTION]
> 
> - This utility is in **very early development**.
> - Users will assume all risk of using this tooling, including overwriting product metadata.
> - This tool is not supported by Red Hat Support.
> - Every API used here is subject to change.


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

0) (Prereq) Make sure your environment has the necessary variables

```bash
export CONNECT_API_TOKEN=yourtoken
export CONNECT_ORG_ID=000000000
```

1) Scaffold your new Product Listing (or fetch an existing one)

```bash
productctl product create [--from-discovery-json /path/to/discovery.json] my.product.yaml
```

Or fetch an existing listing:

```bash
productctl product fetch 000111222333 > my.product.yaml
```

2) Make alterations to your Product Listing, add/remove components, etc. 

3) Apply your Product Listing

```bash
productctl product apply my.product.yaml
```

4) Repeat until all metadata is configured to your liking.


## Getting Started

See our [Getting Started](docs/GETTING_STARTED.md) guide.