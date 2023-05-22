---
title: "Staging"
weight: 90
date: 2023-05-17T17:03:41-04:00
---

*Note: This page is for internal Ensign development and will probably not be very useful to Ensign users. The staging environment has the latest code deployed frequently, may introduce breaking changes, and has it's data routinely deleted.*

## Staging Environment

Ensign developers can access the staging environment in order to perform testing and development or to QA release candidates before they are deployed.

To get started, make sure that you've created an API Key in the staging environment using the Beacon UI at [https://ensign.world](https://ensign.world). Once you've obtained those credentials, add the following environment variables so that your script can access the credentials:

- `$ENSIGN_CLIENT_ID`
- `$ENSIGN_CLIENT_SECRET`

If you're working on the Go SDK in staging, make sure you have the latest version from the commit rather than the latest tagged version so that your client code is up to date with what is in staging:

```bash
$ go get github.com/rotationalio/go-ensign@main
```

By default the Ensign client connects to the Ensign production environment. To connect to Staging you need to specify the staging endpoints in your credentials:

```go
client, err := ensign.New(&ensign.Options{
    Endpoint: "staging.ensign.world:443",
    ClientID: os.GetEnv("ENSIGN_CLIENT_ID"),
    ClientSecret: os.GetEnv("ENSIGN_CLIENT_SECRET"),
    AuthURL: "https://auth.ensign.world",
})
```

If you're feeling extra, you can also use the `ensign.ninja:443` endpoint which is an alias for `staging.ensign.world:443`.