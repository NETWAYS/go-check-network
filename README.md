# go-check-network

go-check-network is a collection of modules for the development of monitoring plugins using network protocols.

See also:

* https://github.com/NETWAYS/go-check

We decided to create a dedicated collection for this code to keep the `go-check` module small and focused.

## http

The `checkhttp` module provides packages for the HTTP protocol.

### config

The go-check-network/http/config package provides helpers to configure HTTP connections (e.g. RoundTrippers, TLSConfig, etc.)

### mock

The go-check-network/http/config package provides additions to the jarcoal/httpmock module.

# License

Copyright (c) 2023 [NETWAYS GmbH](mailto:info@netways.de)

This library is distributed under the GPL-2.0 or newer license found in the [COPYING](./COPYING)
file.
