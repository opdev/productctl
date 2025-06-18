# Getting Started

This document covers how to get started defining your Red Hat Partner Product
Listings with **productctl**.

Before you begin, please make sure you have met the [prerequisites](./PREREQS.md).

### Enabling Shell Completions

If you have a supported shell (see `productctl completions`), you can generate a
shell completion configuration and use that to help you navigate the CLI.

E.g. for BASH

```
productctl completion bash > completion.sh && source completion.sh
```

### Enabling IDE Integration

We highly recommend leveraging the included jsonschema to help you write your
Product Listings. Ultimately, we're manipulating large-ish YAML documents with
some fields that have enumerations with which you may not be familiar off-hand.
We strive to give you as much insight as we can from within your editor. The
easiest way we've found to do this is to provide you with a jsonschema that can
be used with a YAML language server of your choosing.

See [Enabling IDE integration](./USING_JSONSCHEMA.md) for more information on how to configure this.

### Creating your first product listing

The **productctl** tool can be used populate a brand new product listing.

```bash
productctl product create my.product.yaml
```

This will produce a YAML declaration that can be used to configure your product
listings on disk at the specified file. Note that this does not create any
resources on your account. The file created here is just the starting point for
you to edit.

```yaml
# my.product.yaml
kind: ProductListing
spec:
  descriptions:
    long: This can contain long form content about your product.
    short: A brief synopsis
  name: My New Product
  type: container stack
with: {}
```
Feel free to modify this baseline declaration how you see fit, updating the
`name`, `descriptions`, etc.

Use the following command to apply your declaration to your partner
account

```bash
productctl product apply my.product.yaml
```

If you log into productctl, you should then see your product listing.
Continue iterating on your listing until it is ready for publishing.

### Creating a product listing from [discovered workloads](https://github.com/opdev/discover-workload)

The **productctl** tool can read the discover.json produced by
**[discover-workload](https://github.com/opdev/discover-workload)**.

The **[discover-workload](https://github.com/opdev/discover-workload)** CLI allows
you to discovery containers that are installed by your workload in select
namespaces, and with specified filters. You can then use this information to use
as a starting point for your scaffolded Product Listing. See the linked projects
for more details.

Once you've generated a discover file, you can pass this to **productctl** when
scaffolding your resource.

```bash
productctl product create my.product.yaml --from-discovery-json /path/to/discovery.json
```

### Creating and adding components to your product listing

You can also add components (previously, certification projects, or
"certprojects") to your product listing by adding them to `.with.components`.
Components that don't exist in the backend (identified by the presence of the
`_id` field in the component declaration) will be created for you, given they
have the minimum required information.

```yaml
# my.product.yaml
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

```bash
productctl product apply my.product.yaml
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
ID and update your declaration locally with the component's details when you
apply your product

```yaml
# my.product.yaml
kind: ProductListing
spec:
  name: My New Product
  # ... other fields
with: 
  components:
  - _id: 12903123123123123123
  - _id: 12273948321321321321
```

Run the `productctl product apply` command to update your declaration with
component metadata, and bind the components to your product listing.

### How to remove components from a product listing

Simply remove the `.with.components` altogether to completely remove all
components from a product listing. Other than being disassociated from the
product listing, components will not be impacted archived or changed.

```yaml
# my.product.yaml
kind: ProductListing
spec:
  _id: 81228123123123123123
  name: My New Product
  # ... other fields
with: {} # or remove .with completely
```

When your product is applied, the specified components will no longer be bound
to this product listing.

### Archiving Components / Deleting Product Listings

You can delete an entire Product Listing, as well as archive all attached
components using the `cleanup` command, and pointing to a Product Listing
reference.

```bash
productctl product cleanup my.product.yaml
```

If you need to archive or delete a component or product listing directly, you
can do so using the target's `_id` value. See the `productctl util` subcommand
for instructions on how to do this.

### Other operations against Product Listings and Components, including Publishing

The **productctl** command only allow for a subset of all operations you might
need to do against a given Product Listing or Component. For everything else,
you'll need to log into your Red Hat Partner Connect Dashboard.

## What to do Next

Creating your Product Listing is just one aspect of the Certification process.

After this, you'll need to follow the [Red Hat Software Certification
Guide](https://docs.redhat.com/en/documentation/red_hat_software_certification)
for certifying your Components, and submitting those results to Red Hat.