# productctl

A power tool for [Red Hat Software Certification
Partners](https://access.redhat.com/documentation/en-us/red_hat_software_certification/2024/html/red_hat_software_certification_workflow_guide/index).

The `productctl` utility allows Red Hat Software Certification Partners to
easily iterate on the definition of their product listings in a familiar way.

> [!CAUTION]
> 
> - This utility is in **very early development**.
> - Users will assume all risk of using this tooling, including overwriting product metadata.
> - This tool is not supported by Red Hat Support.
> - Every API used here is subject to change.

## Prerequisites

You will need an established Red Hat Partner Connect account. This tool will not
create this for you. 

## Getting Started

You'll need to set two environment variables to use `productctl`. 

| Env Var | Description |
|-|-|
|`CONNECT_ORG_ID` |The organization you're working against. Helps in filtering queries.|
|`CONNECT_API_TOKEN`| Your API token. Used to scope requests just to your project.|

### Getting an API token:

Log into Red Hat Partner Connect and access this URL:
https://connect.redhat.com/account/api-keys

### Getting your ORG ID

Log into Red Hat Partner Connect and access this URL:
https://connect.redhat.com/account/company-profile. Your ORG ID should be listed
at the top of this UI.

## Workflow Examples

```
A basic CLI useful for helping Red Hat Certification Partners define their Product Listings, and create and manage certification projects associated with those product listings.

Usage:
  productctl [command]

Available Commands:
  alpha            Alpha commands that may be removed or modified at any point
  apply-product    Apply changes to Partner product listings from the input file.
  completion       Generate the autocompletion script for the specified shell
  fetch-product    Get a pre-existing product listing
  help             Help about any command
  new-product      Start building a new product listing.
  sanitize-product Cleans declaration for re-use and emits to stdout
  version          Prints the version information

Flags:
      --custom-endpoint string   Define a custom API endpoint. Supersedes predefined environment values like "prod" if set
      --env string               The catalog API environment to use. Choose from stage, prod (default "prod")
  -h, --help                     help for productctl
      --log-level string         The verbosity of the tool itself. Ex. error, warn, info, debug (default "info")
  -v, --version                  version for productctl

Use "productctl [command] --help" for more information about a command.
```

### Creating your first product listing

The **productctl** tool can be used populate a brand new product listing.

```bash
productctl new-product my-product.yaml
```

This will produce a Kubernetes-like declaration that can be used to configure
your product listings on disk at the specified file. Note that this does not
create any resources on your account. The file created here is just the starting
point for you to edit.

```yaml
# my-product.yaml
kind: ProductListing
spec:
  descriptions:
    long: This can contain long form content about your product.
    short: A brief synopsis
  name: My New Product
  type: container stack
with: {}
```

The object defined in `.spec` describes your product listing. A subset of the
fields defined in the [API
schema](https://catalog.redhat.com/api/containers/docs/objects/ProductListing.html?tab=Fields)
can be set in this declaration. Try out the following alpha feature to generate
a resource schemajson that may be used via LSP in your editor of choice.

```
productctl alpha lsp-completion > resource.schema.json
```

Use the **apply-product** subcommand to apply your declaration to your partner account

```
productctl apply-product my-product.yaml
```

If you log into productctl, you should then see your product listing.
Continue iterating on your listing until it is ready for publishing.

### Creating a product listing from [discovered workloads](https://github.com/opdev/discover-workload)

The **productctl** tool can read the discover.json produced by
**[discover-workload](https://github.com/opdev/discover-workload)**.

```
productctl new-product --from-discovery-json /path/to/discovery.json
```

Before doing so, the **discover.json** should be modified to ensure it contains
at most one entry for a given container image. If this is forgotten, the
produced product listing declaration can be modified to accomplish the same
outcome.

### Creating and adding components to your product listing

You can also add components (previously, certification projects, or
"certprojects") to your product listing by adding them to `.with.components`.
Components that don't exist in the backend (identified by the presence of the
`_id` field in the component declaration) will be created for you, given they
have the minimum required information.

```yaml
# my-product.yaml
kind: ProductListing
spec:
  descriptions:
    long: This can contain long form content about your product.
    short: A brief synopsis
  name: My New Product
  type: container stack
with: 
  components:
  # minimum information required for a component
  - name: My First Component
    type: Containers
    container:
      distribution_method: rhcc
      type: container
      os_content_type: Red Hat Universal Base Image (UBI)
```

With your components defined, apply your product again.

```
productctl apply-product my-product.yaml
```

If your declaration was successfully applied, your declaration will update
itself on disk, adding `_id` values and any additional server-side set default
settings.

Components are changed only by modifying `.with.components`. Changes to
`.spec.cert_projects` do not impact your product listing. This field should be
treated as read-only, representing the components currently bound to your
product listing from the server's perspective.

### Adding pre-existing components to your product listing

Adding a pre-existing component to your product listing is as simple as adding
its `_id` value to the declaration. The **productctl** tool will detect the
ID and update your declaration locally the component's details when
**apply-product** is called.

```yaml
# my-product.yaml
kind: ProductListing
spec:
  name: My New Product
  # ... other fields
with: 
  components:
  - _id: 12903123123123123123
  - _id: 12273948321321321321
```

Run the **apply-product** subcommand to update your declaration with component
metadata, and bind the components to your product listing.

### How to remove components from a product listing

Simply remove the `.with.components` altogether to completely remove all
components from a product listing. Other than being disassociated from the
product listing, components will not be impacted.

```yaml
# my-product.yaml
kind: ProductListing
spec:
  _id: 81228123123123123123
  name: My New Product
  # ... other fields
with: {} # or remove .with completely
```

When your product is applied, the specified components will no longer be bound
to this product listing.

### Publishing Products and Components

At the time of this writing, users will need to navigate to the [Red Hat Partner
Connect Dashboard](https://connect.redhat.com) to publish their components and
products. This may change in the future.

### Troubleshooting

Bump the verbosity on the logger to get more verbosity on what's happening under
the hood.

```bash
productctl --log-level debug ...
```
