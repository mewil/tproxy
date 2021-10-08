# TProxy

[![Docker Hub](https://img.shields.io/docker/pulls/mewil/tproxy.svg)](https://hub.docker.com/repository/docker/mewil/tproxy)


TProxy is a container-based HTTPS proxy for [Tailscale](https://tailscale.com/). It's based on [this blog post](https://rnorth.org/tailscale-docker/) by Richard North, but is designed to be used with Tailscale's [TLS certs feature](https://tailscale.com/blog/tls-certs/).

## Options

TProxy can be configured using the following environment variables:

| Variable                     | Description                                                                                                        |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `TPROXY_TARGET_ADDR`         | The target address for the proxy e.g. `http://some-container:3000`                                                 |
| `TPROXY_USER_HEADER`         | A header in which the proxy will write the Tailscale login name of a request's remote address, not used by default |
| `TPROXY_TAILSCALE_AUTH_KEY`  | An [auth key](https://tailscale.com/kb/1085/auth-keys/) from Tailscale                                             |
| `TPROXY_TAILSCALE_HOST_NAME` | The hostname for `tailscaled`                                                                                      |
