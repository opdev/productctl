#!/usr/bin/env python3

import os
from sys import stderr

import gql
from gql.transport.aiohttp import AIOHTTPTransport
from graphql import build_client_schema, print_schema, get_introspection_query

def main():
    try:
        # We'll hit the the Catalog API unless specified otherwise. 
        url = "https://catalog.redhat.com/api/containers/graphql/"
        env_url = os.getenv("PYXIS_URL")
        if env_url:
            print(f"Using alternate URL {env_url}", file=stderr)
            url = env_url

        # Staging or Development endpoints may require proxy support. We're
        # respecting just this proxy variable until additional ones are
        # required.
        session_args = None
        p = os.getenv("HTTPS_PROXY")
        if p:
            print(f"Respecting HTTPS_PROXY value.", file=stderr)
            session_args = {"proxy": p}


        transport = AIOHTTPTransport(
            url=url,
            client_session_args=session_args
        )

        # Introspect
        q = get_introspection_query(descriptions=True)
        res = gql.Client(transport=transport, fetch_schema_from_transport=True).execute(gql.gql(q))

        # Build schema from introspection
        schema = build_client_schema(res)

        # Convert schema to SDL
        sdl = print_schema(schema)
        print(sdl)
    except Exception as e:
        raise(e)

if __name__ == "__main__":
    main()

