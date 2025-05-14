# Component Mappings - Configuration Reference

Below are example/annotated component mapping configurations. These are used
with the `productctl alpha certify` subcommands in order to run certificaton
tooling against your products.

The base component mapping configuration is generated from a given product
listings. Generally, each object in these sections corresponds to a given
component. In some cases, you may want to certify a collection of tags, and so
this configuration allows you to define all tags you wish to certify for a given
component so that you can run them all at once.

Not all tooling flags are passed through to this configuration file 

## Containers

```yaml
container_components:
  # componentid, Generally, this value will be the component's _id value, and
  # is scaffolded for you when generated using
  # `productctl alpha generate-component-mappings`.
  componentid:
    # This is the image being certified for this component. No tag or digest should
    # be specified here.
    image_ref: registry.example.com/placeholder/placeholder
    # You can set each of these flag on a per-tag basis as well.
    tool_flags:
        # Whether to submit results. Remember to pass in --catalog-api-token
        # to `preflight alpha certify containers`. We do not accept that
        # credential in this file.
        submit: true
        
        # This should be within the userfiles directory, e.g. 
        # ./userfiles/auth.json, assuming you pass --userfiles-dir=$PWD/userfiles
        # when calling `productctl alpha certify containers`. The value here
        # is prepended with the value provided to that flag, so make sure that what's
        # stored here is just the filepath relative to that directory.
        docker_config: auth.json

        # Platform, for limiting the platform that preflight will run against.
        # Omit if you want preflight to take its default behavior.
        # platform: arm64

        # Pyxis Env, for users who have access to lower environments.
        # pyxis_env: stage
        #
        # Log Level, for users who need more verbosity from the certification tool.
        # values should be coming from those specified by the tool in scope (preflight).
        # loglevel: debug
    tags:
      ## A tag you want to certify.
      - tag: placeholder
      - tag: othertag
        # This tag has its own tool_flags. Note that these aren't layered
        # over the top-level tool_flags defined for this object. You must repeat
        # any flags that you also want to carry over here.
        tool_flags:
            # Don't submit this 
            submit: false
```

## Operators

The operator certification tooling is the same as the container certification
tooling. Many of the same configurations apply.

```yaml
operator_components:
  ## Imaginary Operator Component
  componentid:
    image_ref: quay.io/opdev/simple-demo-operator-bundle
    ## The operator index image (or catalog) that contains the image_ref at the
    ## listed tags. Per-tag index_images can also be configured alongside the tag
    ## definition
    index_image: 'quay.io/opdev/simple-demo-operator-catalog:latest'
    ## The tool_flags directive contains any customized flags to set on the
    ## the certification tooling.
    tool_flags:
      # kubeconfig, this path is relative to the --userfiles-dir value. That is,
      # if you specify --userfiles-dir=$PWD/userfiles, then this file must exist at path
      # $PWD/userfiles/component-kubeconfig.yaml.
      #
      kubeconfig: component-kubeconfig.yaml

      # Update your verbosity, values should match those of the underlying tool (Preflight)
      loglevel: trace

      # Below scorecard-related configuration allow the caller to adjust how scorecard
      # is executed as a part of operator certification.
      # 
      # scorecard_namespace: myscorecardns
      # scorecard_image: example.com/namespace/image:tag
      # scorecard_wait_time: 3m
      # scorecard_service_account: myscorecardsa
      #
      # channel, if you want to override the channel from the index_image that is under test
      # channel: stable
    tags:
      - tag: latest
        # You can pass through per-tag tool_flags here if you would like.
        # Remember that this is not layered on top of the top-level tool_flags, and must
        # specify all flags that are in scope for this tag if this value is set.
        # tool_flags: {}
```

## Helm Charts

The Helm certification tooling has different options than operator and container
mapping configuration. However, semantics around file references, and ensuring
they're relative to the `--userfiles-dir` value remain the same.

```yaml
helm_chart_components:
  componentid:
    chart_uri: 'https://example.com/your/path/to/your/chart/tar/is/chart-0.0.1.tgz'
    tool_flags:
      debug: true
      # kubeconfig must be relative to the --userfiles-dir path.
      kubeconfig: helm-chart-kubeconfig.yaml
      # Additional values files to pass in. Note that these must be relative to the
      # --userfiles-dir path.
      values_files:
        - my.values.yaml 
      # vendor_type, for overriding the vendor type, which controls which policy is applied
      # to your chart. Omit this for default behavior.
      # vendor_type: partner
      #
      # with_set_items allows you to define arbitrary keys and values for your chart, similar to
      # the --set flag for Helm
      with_set_items:
        - foo=bar
        - this=that
      # web_catalog_only, for partners that are not shipping their charts via traditional means, and only want
      # to be listed in the web catalog.
      web_catalog_only: true
```