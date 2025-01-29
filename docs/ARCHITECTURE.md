# Architecture Overview

This tool is effectively a GraphQL client interacting with the Red Hat Catalog
API available at https://catalog.redhat.com/api/containers/graphql/. It will
perform predefined GraphQL queries and mutations to define Product Listings and
their associated Components[1].

## Generating GraphQL Client Code

This project uses a code generator
**[genqlient](https://github.com/Khan/genqlient)** in order to generate client
code for the operations used by this tool.

This code generator generates strongly-typed code based on the schema definition
provided by the GraphQL endpoint. All generated code lives in the
**internal/genpyxis/** directory, and is wired to generate this code when `go
generate` is called.

The configuration for **genqlient** can be found in the same directory.

The **genqlient** generator works by accepting the GraphQL schema as well as a
defined set of operations needed by your code, and generating just the necessary
functions and methods.

## Introspecting the GraphQL schema for the Catalog API

Most GraphQL endpoints support the use of introspection to enumerate schema, and
the Python reference implementation **gql** provides hooks to both run the
introspection and convert it to GraphQL Schema Definition Language ("SDL")
format.

For this reason, this repository includes a script to generate the schema at will.

The [Makefile](../Makefile) provides hooks for generating all of the necessary
bits required to generate a fresh schema and fresh client code, and serves as
the documentation for how this process is executed.

The Make targets in scope (e.g. `make generate`) should generate all code
required to run. Python code will be scoped to a virtualenv local to the
repository to run the script with the appropriate dependencies.

## Footnotes:

1. The term "Component" is interchangeable with the term "Certification
   Project", which is the historical term for the object stored in the Red Hat
   Catalog API. This tool should use the term "Component", while plumbing this
   tool leverages may use the previous term "Certification Project".
