# go-check-network

go-check-network is a collection of modules for the development of monitoring plugins using network protocols.

See also:

* https://github.com/NETWAYS/go-check

We decided to create a dedicated collection for this code to keep the `go-check` module small and focused.

## http

The `checkhttp` module provides packages for the HTTP protocol.

### config

The go-check-network/http/config package provides helpers to configure HTTP connections (e.g. RoundTrippers, TLSConfig, etc.)

Examples:

```
// Example for TLSConfig from files
tlsConfig, err := checkhttp.NewTLSConfig(&checkhttpconfig.TLSConfig{
    InsecureSkipVerify: false,
    CAFile:             myCAFile,
    KeyFile:            myKeyFile,
    CertFile:           myCertFile,
})

// Some sane defaults
var rt http.RoundTripper = &http.Transport{
    Proxy: http.ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    }).DialContext,
    TLSHandshakeTimeout: 10 * time.Second,
    TLSClientConfig:     tlsConfig,
}

// Example for Token Auth Roundtripper
rt = checkhttpconfig.NewAuthorizationCredentialsRoundTripper("Bearer", "secret-bearer-token", rt)

// Example for Basic Auth Roundtripper
rt = checkhttpconfig.NewBasicAuthRoundTripper("my-user", "password123", rt)
```

### mock

The go-check-network/http/config package provides additions to the jarcoal/httpmock module.

# License

Copyright (c) 2023 [NETWAYS GmbH](mailto:info@netways.de)

This library is distributed under the GPL-2.0 or newer license found in the [COPYING](./COPYING)
file.
