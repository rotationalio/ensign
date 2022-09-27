---
title: "Configuration"
weight: 10
bookFlatSection: false
bookToc: true
bookHidden: false
bookCollapseSection: false
bookSearchExclude: false
---

# Configuration

Ensign services are primarily configured using environment variables and will respect [dotenv files](https://github.com/joho/godotenv) in the current working directory. The canonical reference of the configuration for an Ensign service is the `config` package of that service (described below). This documentation enumerates the most important configuration variables, their default values, and any hints or warnings about how to use them.

{{< hint info >}}
**Required Configuration**<br />
If a configuration parameter does not have a default value that means it is required and must be specified by the user! If the configuration parameter does have a default value then that environment variable does not have to be set.
{{< /hint >}}

<!--more-->

## Ensign

The Ensign node is a replica of the Ensign eventing system. Its environment variables are all prefixed with the `ENSIGN_` tag. The primary configuration is as follows:

| EnvVar             | Type   | Default | Description                                                                                                    |
|--------------------|--------|---------|----------------------------------------------------------------------------------------------------------------|
| ENSIGN_MAINTENANCE | bool   | false   | Sets the node to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| ENSIGN_LOG_LEVEL   | string | info    | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                             |
| ENSIGN_CONSOLE_LOG | bool   | false   | If true will print human readable logs instead of JSON logs for machine consumption.                           |
| ENSIGN_BIND_ADDR   | string | :5356   | The address and port the Ensign service will listen on.                                                        |

### Sentry

Ensign uses [Sentry](https://sentry.io/) to assist with error monitoring and performance tracing. Configure Ensign to use Sentry as follows:

| EnvVar                          | Type    | Default     | Description                                                                                       |
|---------------------------------|---------|-------------|---------------------------------------------------------------------------------------------------|
| ENSIGN_SENTRY_DSN               | string  |             | The DSN for the Sentry project. If not set then Sentry is considered disabled.                    |
| ENSIGN_SENTRY_SERVER_NAME       | string  |             | Optional - a server name to tag Sentry events with.                                               |
| ENSIGN_SENTRY_ENVIRONMENT       | string  |             | The environment to report (e.g. development, staging, production). Required if Sentry is enabled. |
| ENSIGN_SENTRY_RELEASE           | string  | {{version}} | Specify the release of Ensign for Sentry tracking. By default this will be the package version.   |
| ENSIGN_SENTRY_TRACK_PERFORMANCE | bool    | false       | Enable performance tracing to Sentry with the specified sample rate.                              |
| ENSIGN_SENTRY_SAMPLE_RATE       | float64 | 0.2         | The percentage of transactions to trace (0.0 to 1.0).                                             |
| ENSIGN_SENTRY_DEBUG             | bool    | false       | Set Sentry to debug mode for testing.                                                             |

Sentry is considered **enabled** if a DSN is configured. Performance tracing is only enabled if Sentry is enabled *and* track performance is set to true. If Sentry is enabled, an environment is required, otherwise the configuration will be invalid.

Generally speaking, Ensign should enable Sentry for panic reports but should not enable performance tracing as this slows down the server too much.

### Monitoring

Ensign uses Prometheus for metrics and observability. The prometheus metrics server is configured as follows:

| EnvVar                      | Type   | Default | Description                                             |
|-----------------------------|--------|---------|---------------------------------------------------------|
| ENSIGN_MONITORING_ENABLED   | bool   | true    | If true, the Prometheus metrics server is served.       |
| ENSIGN_MONITORING_BIND_ADDR | string | :1205   | The address and port the metrics server will listen on. |
| ENSIGN_MONITORING_NODE_ID   | string |         | Optional - a server name to tag metrics with.           |

## Tenant

The Tenant API powers the user front-end for tenant management and configuration. Its environment variables are all prefixed with the `TENANT_` tag. The primary configuration is as follows:

| EnvVar               | Type   | Default               | Description                                                                                                      |
|----------------------|--------|-----------------------|------------------------------------------------------------------------------------------------------------------|
| TENANT_MAINTENANCE   | bool   | false                 | Sets the server to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| TENANT_BIND_ADDR     | string | :8080                 | The address and port the Tenant service will listen on.                                                          |
| TENANT_MODE          | string | release               | Sets the Gin mode, one of debug, release, or test.                                                               |
| TENANT_LOG_LEVEL     | string | info                  | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                               |
| TENANT_CONSOLE_LOG   | bool   | false                 | If true will print human readable logs instead of JSON logs for machine consumption.                             |
| TENANT_ALLOW_ORIGINS | string | http://localhost:3000 | A comma separated list of allowed origins for CORS. Set to "*" to allow all origins.                             |

### Sentry

Tenant uses [Sentry](https://sentry.io/) to assist with error monitoring and performance tracing. Configure Tenant to use Sentry as follows:

| EnvVar                          | Type    | Default     | Description                                                                                       |
|---------------------------------|---------|-------------|---------------------------------------------------------------------------------------------------|
| TENANT_SENTRY_DSN               | string  |             | The DSN for the Sentry project. If not set then Sentry is considered disabled.                    |
| TENANT_SENTRY_SERVER_NAME       | string  |             | Optional - a server name to tag Sentry events with.                                               |
| TENANT_SENTRY_ENVIRONMENT       | string  |             | The environment to report (e.g. development, staging, production). Required if Sentry is enabled. |
| TENANT_SENTRY_RELEASE           | string  | {{version}} | Specify the release of Ensign for Sentry tracking. By default this will be the package version.   |
| TENANT_SENTRY_TRACK_PERFORMANCE | bool    | false       | Enable performance tracing to Sentry with the specified sample rate.                              |
| TENANT_SENTRY_SAMPLE_RATE       | float64 | 0.2         | The percentage of transactions to trace (0.0 to 1.0).                                             |
| TENANT_SENTRY_DEBUG             | bool    | false       | Set Sentry to debug mode for testing.                                                             |

Sentry is considered **enabled** if a DSN is configured. Performance tracing is only enabled if Sentry is enabled *and* track performance is set to true. If Sentry is enabled, an environment is required, otherwise the configuration will be invalid.

## Quarterdeck

The Quarterdeck API handles authentication and authorization as well as API keys and billing management for the Ensign managed service. Its environment variables are all prefixed with the `QUARTERDECK_` tag. The primary configuration is as follows:

| EnvVar                    | Type   | Default               | Description                                                                                                      |
|---------------------------|--------|-----------------------|------------------------------------------------------------------------------------------------------------------|
| QUARTERDECK_MAINTENANCE   | bool   | false                 | Sets the server to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| QUARTERDECK_BIND_ADDR     | string | :8088                 | The address and port the Quarterdeck service will listen on.                                                     |
| QUARTERDECK_MODE          | string | release               | Sets the Gin mode, one of debug, release, or test.                                                               |
| QUARTERDECK_LOG_LEVEL     | string | info                  | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                               |
| QUARTERDECK_CONSOLE_LOG   | bool   | false                 | If true will print human readable logs instead of JSON logs for machine consumption.                             |
| QUARTERDECK_ALLOW_ORIGINS | string | http://localhost:3000 | A comma separated list of allowed origins for CORS. Set to "*" to allow all origins.                             |

### Sentry

Quarterdeck uses [Sentry](https://sentry.io/) to assist with error monitoring and performance tracing. Configure Quarterdeck to use Sentry as follows:

| EnvVar                          | Type    | Default     | Description                                                                                       |
|---------------------------------|---------|-------------|---------------------------------------------------------------------------------------------------|
| QUARTERDECK_SENTRY_DSN               | string  |             | The DSN for the Sentry project. If not set then Sentry is considered disabled.                    |
| QUARTERDECK_SENTRY_SERVER_NAME       | string  |             | Optional - a server name to tag Sentry events with.                                               |
| QUARTERDECK_SENTRY_ENVIRONMENT       | string  |             | The environment to report (e.g. development, staging, production). Required if Sentry is enabled. |
| QUARTERDECK_SENTRY_RELEASE           | string  | {{version}} | Specify the release of Ensign for Sentry tracking. By default this will be the package version.   |
| QUARTERDECK_SENTRY_TRACK_PERFORMANCE | bool    | false       | Enable performance tracing to Sentry with the specified sample rate.                              |
| QUARTERDECK_SENTRY_SAMPLE_RATE       | float64 | 0.2         | The percentage of transactions to trace (0.0 to 1.0).                                             |
| QUARTERDECK_SENTRY_DEBUG             | bool    | false       | Set Sentry to debug mode for testing.                                                             |

Sentry is considered **enabled** if a DSN is configured. Performance tracing is only enabled if Sentry is enabled *and* track performance is set to true. If Sentry is enabled, an environment is required, otherwise the configuration will be invalid.

# Development

{{< hint danger >}}
**Keep up to Date!**<br />
It is essential that we keep this configuration documentation up to date. The devops team uses it to ensure its services are configured correctly. Any time a configuration is changed ensure this documentation is also updated!
{{< /hint >}}

TODO: this section will discuss envconfig, how to interpret environment variables from the configuration struct, how to test configuration, and how to add and change configuration variables. This section should also discuss dotenv files, docker compose, and all of the places where configuration can be influenced (e.g. GitHub actions for React builds).