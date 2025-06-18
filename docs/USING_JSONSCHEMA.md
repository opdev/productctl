# Enabling IDE Integration

The `productctl` CLI allows you to configure your Red Hat Product Listings (and
their associated components) using a YAML declaration on your disk. These YAML
declarations have a lot of fields that you can populate, and so we highly
recommend configuring your editor to use our *jsonschema* to guide you as your
edit your Product Listing.

> [!Warning]
> 
> The jsonschema under active expansion. Not all fields have been included in
> the schema yet. If you notice a field is missing, you can look at the upstream
> [object
> reference](https://catalog.redhat.com/api/containers/docs/objects/ProductListing.html?tab=Fields)
> for guidance, and please open an issue to let us know something is missing!

You can generate the schema using the following command:

```bash
productctl product jsonschema > productctl.productlisting.schema.json
```

You should then be able to configure YAML extensions to your favorite IDE to
utilize the jsonschema to provide intellisense while you write your Product
Listing declaration.

## E.g., for Visual Studio Code

Visual Studio Code allows you to register jsonschema's for YAML
documents using the [YAML
extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml).

This extension can be configured to point to a local schemajson on your disk.
Here's an example configuration, which should go in your Visual Studio Code
**settings.json**:

```json
    // ... other configuration in your settings.json
    "yaml.schemas": {
        "/home/user/jsonschemas/productctl.productlisting.schema.json": [
            "*.product.yaml",
            "*.product.yml",
            "*.productlisting.yaml",
            "*.productlisting.yml",
        ]
    },
    // ... other configuration in your settings.json
```

The extension should then automatically detect your product listings if your
file name matches any of the listings above.

From there, you can use the "Trigger Suggestion" command from the [Command
Palette](https://code.visualstudio.com/api/ux-guidelines/command-palette) within
Visual Studio Code, or whatever key chord you have bound for the same
functionality (e.g. Ctrl+Space). Some fields also have "Quick Info" available by
pressing the "Trigger Suggestion" key chord twice (or by using the default
binding for this, e.g. `Ctrl+K, Ctrl+I`).