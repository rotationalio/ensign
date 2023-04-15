---
title: "Configuration"
weight: 10
bookFlatSection: false
bookToc: true
bookHidden: false
bookCollapseSection: false
bookSearchExclude: false
---

*Note: This page is for internal Ensign development and will probably not be very useful to Ensign users.*

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

### Monitoring

Ensign uses Prometheus for metrics and observability. The Prometheus metrics server is configured as follows:

| EnvVar                      | Type   | Default | Description                                             |
|-----------------------------|--------|---------|---------------------------------------------------------|
| ENSIGN_MONITORING_ENABLED   | bool   | true    | If true, the Prometheus metrics server is served.       |
| ENSIGN_MONITORING_BIND_ADDR | string | :1205   | The address and port the metrics server will listen on. |
| ENSIGN_MONITORING_NODE_ID   | string |         | Optional - a server name to tag metrics with.           |

### Storage

The Ensign storage configuration defines on disk where Ensign keeps its data. Configure storage as follows:

| EnvVar                   | Type   | Default | Description                                                                                              |
|--------------------------|--------|---------|----------------------------------------------------------------------------------------------------------|
| ENSIGN_STORAGE_READ_ONLY | bool   | false   | If true then no writes or deletes will be allowed to the database.                                       |
| ENSIGN_STORAGE_DATA_PATH | string |         | The path to a directory on disk where Ensign will store its meta and event data.                         |
| ENSIGN_STORAGE_TESTING   | bool   | false   | If true then a mock store will be opened rather than a leveldb store (should not be used in production). |

Note that the Ensign data path must be to a directory. If the directory does not exist, it is created. An error occurs if the path is to a file or the process doesn't have permissions to access the directory. Ensign will open two different data stores in the data path: one for metadata and the other to store event data locally.

If the testing flag is set to true, a mock store is created that can be used in unit and integration tests.

### Authentication

Ensign uses Quarterdeck to authenticate and authorize requests. This configuration defines how Ensign accesses public keys for JWT verification and how the authentication interceptor behaves.

| EnvVar                           | Type     | Default                                           | Description                                                                       |
|----------------------------------|----------|---------------------------------------------------|-----------------------------------------------------------------------------------|
| ENSIGN_AUTH_KEYS_URL             | string   | https://auth.rotational.app/.well-known/jwks.json | The path to the public keys used to verify JWT tokens.                            |
| ENSIGN_AUTH_AUDIENCE             | string   | https://ensign.rotational.app                     | The audience to verify inside of JWT tokens.                                      |
| ENSIGN_AUTH_ISSUER               | string   | https://auth.rotational.app                       | The issuer to verify inside of JWT tokens.                                        |
| ENSIGN_AUTH_MIN_REFRESH_INTERVAL | duration | 5m                                                | The minimum time to wait before refreshing the JWKS public keys in the validator. |

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

Generally speaking, Ensign should enable Sentry for panic reports but should not enable performance tracing as this slows down the server too much. Note also that the `sentry.Config` object has a field `Repanic` that should not be set by the user. This field is used to manage panics in chained interceptors.

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

### Authentication

Tenant uses Quarterdeck to authenticate and authorize requests. The following configuration defines how Tenant validates JWT tokens passed in from the user that were created by Quarterdeck and secures cookies:

| EnvVar                    | Type   | Default                                           | Description                                                                         |
|---------------------------|--------|---------------------------------------------------|-------------------------------------------------------------------------------------|
| TENANT_AUTH_KEYS_URL      | string | https://auth.rotational.app/.well-known/jwks.json | The path to the public keys used to verify JWT tokens.                              |
| TENANT_AUTH_AUDIENCE      | string | https://rotational.app                            | The audience to verify inside of JWT tokens.                                        |
| TENANT_AUTH_ISSUER        | string | https://auth.rotational.app                       | The issuer to verify inside of JWT tokens.                                          |
| TENANT_AUTH_COOKIE_DOMAIN | string | rotational.app                                    | Defines the domain that is set on cookies including CSRF cookies and token cookies. |

### Database

Tenant connects to a replicated trtl database for its data storage; the trtl connection is configured as follows:

| EnvVar                    | Type   | Default               | Description                                                                                    |
|---------------------------|--------|-----------------------|------------------------------------------------------------------------------------------------|
| TENANT_DATABASE_URL       | string | trtl://localhost:4436 | The endpoint of the trtl database including the scheme (usually the k8s trtl service).         |
| TENANT_DATABASE_INSECURE  | bool   | true                  | Connect to the trtl database without TLS (default true inside of a k8s cluster).               |
| TENANT_DATABASE_CERT_PATH | string |                       | The path to the mTLS certificate of the client to connect to trtl in an authenticated fashion. |
| TENANT_DATABASE_POOL_PATH | string |                       | The path to an x509 pool file with trusted trtl servers to connect to in.                      |

### Quarterdeck

Tenant connects directly to Quarterdeck in order to send pass-through requests from the user (e.g. for registration, login, etc). This is separate from the authentication configuration used in middleware as this configuration is used to setup a Quarterdeck API client.

| EnvVar                            | Type     | Default                     | Description                                                                          |
|-----------------------------------|----------|-----------------------------|--------------------------------------------------------------------------------------|
| TENANT_QUARTERDECK_URL            | string   | https://auth.rotational.app | The Quarterdeck endpoint to create the API client on.                                |
| TENANT_QUARTERDECK_WAIT_FOR_READY | duration | 5m                          | The amount of time to wait for Quarterdeck to come online before exiting with fatal. |

### SendGrid

Tenant uses [SendGrid](https://sendgrid.com/) to assist with email notifications. Configure Tenant to use SendGrid as follows:

| EnvVar                         | Type   | Default              | Description                                                                                          |
|--------------------------------|--------|----------------------|------------------------------------------------------------------------------------------------------|
| TENANT_SENDGRID_API_KEY        | string |                      | API Key to authenticate to SendGrid with.                                                            |
| TENANT_SENDGRID_FROM_EMAIL     | string | ensign@rotational.io | The email address in the "from" field of emails being sent to users.                                 |
| TENANT_SENDGRID_ADMIN_EMAIL    | string | admins@rotational.io | The email address to send admin emails to from the server.                                           |
| TENANT_SENDGRID_ENSIGN_LIST_ID | string |                      | A contact list to add users to if they sign up for notifications.                                    |
| TENANT_SENDGRID_TESTING        | bool   | false                | If in testing mode no emails are actually sent but are stored in a mock email collection.            |
| TENANT_SENDGRID_ARCHIVE        | string |                      | If in testing mode, specify a directory to save emails to in order to review emails being generated. |

SendGrid is considered **enabled** if the SendGrid API Key is set. The from and admin email addresses are required if SendGrid is enabled.

If the Ensign List ID is configured then Tenant will add the contact requesting private beta access to that list, otherwise it will simply add the contact to "all contacts".

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

Note also that the `sentry.Config` object has a field `Repanic` that should not be set by the user. This field is used to manage panics in chained interceptors.

## Quarterdeck

The Quarterdeck API handles authentication and authorization as well as API keys and billing management for the Ensign managed service. Its environment variables are all prefixed with the `QUARTERDECK_` tag. The primary configuration is as follows:

| EnvVar                      | Type   | Default                       | Description                                                                                                      |
|-----------------------------|--------|-------------------------------|------------------------------------------------------------------------------------------------------------------|
| QUARTERDECK_MAINTENANCE     | bool   | false                         | Sets the server to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| QUARTERDECK_BIND_ADDR       | string | :8088                         | The address and port the Quarterdeck service will listen on.                                                     |
| QUARTERDECK_MODE            | string | release                       | Sets the Gin mode, one of debug, release, or test.                                                               |
| QUARTERDECK_LOG_LEVEL       | string | info                          | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                               |
| QUARTERDECK_CONSOLE_LOG     | bool   | false                         | If true will print human readable logs instead of JSON logs for machine consumption.                             |
| QUARTERDECK_ALLOW_ORIGINS   | string | http://localhost:3000         | A comma separated list of allowed origins for CORS. Set to "*" to allow all origins.                             |
| QUARTERDECK_VERIFY_BASE_URL | string | https://rotational.app/verify | The base url to generate verify email links to direct the user to in the email verification path.                |

### SendGrid

Quarterdeck uses [SendGrid](https://sendgrid.com/) to assist with email notifications. Configure Quarterdeck to use SendGrid as follows:

| EnvVar                              | Type   | Default              | Description                                                                                          |
|-------------------------------------|--------|----------------------|------------------------------------------------------------------------------------------------------|
| QUARTERDECK_SENDGRID_API_KEY        | string |                      | API Key to authenticate to SendGrid with.                                                            |
| QUARTERDECK_SENDGRID_FROM_EMAIL     | string | ensign@rotational.io | The email address in the "from" field of emails being sent to users.                                 |
| QUARTERDECK_SENDGRID_ADMIN_EMAIL    | string | admins@rotational.io | The email address to send admin emails to from the server.                                           |
| QUARTERDECK_SENDGRID_ENSIGN_LIST_ID | string |                      | A contact list to add users to if they sign up for notifications.                                    |
| QUARTERDECK_SENDGRID_TESTING        | bool   | false                | If in testing mode no emails are actually sent but are stored in a mock email collection.            |
| QUARTERDECK_SENDGRID_ARCHIVE        | string |                      | If in testing mode, specify a directory to save emails to in order to review emails being generated. |

SendGrid is considered **enabled** if the SendGrid API Key is set. The from and admin email addresses are required if SendGrid is enabled.

If the Ensign List ID is configured then Quarterdeck will add the contact to that list to ensure they receive marketing emails about Ensign.

### Rate Limit

In order to prevent brute force attacks on the Quarterdeck system we've implemented a rate limiting middleware to prevent abuse. Rate limiting is configured as follows:

| EnvVar                            | Type     | Default | Description                                                                                                    |
|-----------------------------------|----------|---------|----------------------------------------------------------------------------------------------------------------|
| QUARTERDECK_RATE_LIMIT_PER_SECOND | float64  | 10      | The maximum number of allowed requests per second.                                                             |
| QUARTERDECK_RATE_LIMIT_BURST      | int      | 30      | Maximum number of requests that is used to track rate of requests (if zero then all requests will be rejected) |
| QUARTERDECK_RATE_LIMIT_TTL        | duration | 5m      | How long an IP address is cached for rate limiting purposes.                                                   |

It is strongly recommended that the default configuration is used.

### Database

| EnvVar                         | Type   | Default                            | Description                                      |
|--------------------------------|--------|------------------------------------|--------------------------------------------------|
| QUARTERDECK_DATABASE_URL       | string | sqlite3:////data/db/quarterdeck.db | The DSN for the sqlite3 database.                |
| QUARTERDECK_DATABASE_READ_ONLY | bool   | false                              | If true only read-only transactions are allowed. |

Quarterdeck uses a Raft replicated Sqlite3 database for authentication. The URI should have the scheme `sqlite3://` and then a path to the database. For a relative path, use `sqlite3:///path/to/relative.db` and for an absolute path use `sqlite3:////path/to/absolute.db`.

### Tokens

| EnvVar                             | Type              | Default                     | Description                                                                                         |
|------------------------------------|-------------------|-----------------------------|-----------------------------------------------------------------------------------------------------|
| QUARTERDECK_TOKEN_KEYS             | map[string]string |                             | The private keys to load into quarterdeck to issue JWT tokens with.                                 |
| QUARTERDECK_TOKEN_AUDIENCE         | string            | https://rotational.app      | The audience to add to the JWT keys for verification.                                               |
| QUARTERDECK_TOKEN_REFRESH_AUDIENCE | string            |                             | An optional additional audience to add only to refresh tokens.                                      |
| QUARTERDECK_TOKEN_ISSUER           | string            | https://auth.rotational.app | The issuer to add to the JWT keys for verification.                                                 |
| QUARTERDECK_TOKEN_ACCESS_DURATION  | duration          | 1h                          | The amount of time that an access token is valid for before it expires.                             |
| QUARTERDECK_TOKEN_REFRESH_DURATION | duration          | 2h                          | The amount of time that a refresh token is valid for before it expires.                             |
| QUARTERDECK_TOKEN_REFRESH_OVERLAP  | duration          | -15m                        | A negative duration that sets how much time the access and refresh tokens overlap in time validity. |
To create an environment variable that is a `map[string]string` use a string in the following form:

```
key1:value1,key2:value2
```

The token keys should be ULIDs keys (for ordering) and a path value to the key pair to load from disk. Generally speaking there should be two keys - the current key and the most recent previous key, though more keys can be added for verification. Only the most recent key will be used to issue tokens, however. For example, here is a valid key map:

```
01GECSDK5WJ7XWASQ0PMH6K41K:/data/keys/01GECSDK5WJ7XWASQ0PMH6K41K.pem,01GECSJGDCDN368D0EENX23C7R:/data/keys/01GECSJGDCDN368D0EENX23C7R.pem
```

{{< hint info >}}
**Future Feature**<br />
Note that in the future quarterdeck will generate its own keys and will not need them to be set as in the configuration above.
{{< /hint >}}

### Reporting

Ensign has a Daily PLG report that is sent to the Rotational admins for product led growth. This reporting tool is configured as follows:

| EnvVar                                 | Type   | Default          | Description                                           |
|----------------------------------------|--------|------------------|-------------------------------------------------------|
| QUARTERDECK_REPORTING_ENABLE_DAILY_PLG | bool   | true             | Enables the Daily PLG report scheduler.               |
| QUARTERDECK_REPORTING_DOMAIN           | string | "rotational.app" | The domain that the report is being generated for.    |
| QUARTERDECK_REPORTING_DASHBOARD_URL    | string |                  | URL to the Grafana dashboard to include in the email. |

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

Note also that the `sentry.Config` object has a field `Repanic` that should not be set by the user. This field is used to manage panics in chained interceptors.

## Beacon

A React app delivers Beacon, the Ensign UI. Its environment variables are all prefixed with the `REACT_APP` tag. The primary configuration is as follows:

| EnvVar                    | Type   | Default | Description                                                                        |
|---------------------------|--------|---------|------------------------------------------------------------------------------------|
| REACT_APP_VERSION_NUMBER  | string |         | The version number of the application build (set by tags in GitHub actions).       |
| REACT_APP_GIT_REVISION    | string |         | The git revision (short) of the application build (set by tags in GitHub actions). |
| REACT_APP_USE_DASH_LOCALE | bool   | false   | If true the "dash" language is included for i18n debugging.                        |

### API Information

Connection information for the backend is specified as follows:

| EnvVar                         | Type   | Default                        | Description                                                           |
|--------------------------------|--------|--------------------------------|-----------------------------------------------------------------------|
| REACT_APP_QUARTERDECK_BASE_URL | string | https://auth.rotational.app/v1 | The endpoint and the version of the API to connect to Quarterdeck on. |
| REACT_APP_TENANT_BASE_URL      | string | https://api.rotational.app/v1  | The endpoint and the version of the API to connect to Tenant on.      |

### Google Analytics

The React app uses [Google Analytics](https://analytics.google.com/) to monitor website traffic. Configure the React app to use Google Analytics as follows:

| EnvVar                          | Type    | Default     | Description                                                                                       |
|---------------------------------|---------|-------------|---------------------------------------------------------------------------------------------------|
| REACT_APP_ANALYTICS_ID          | string  |             | Google Analytics tracking ID for the React App.                                                   |


### Sentry

The React app uses [Sentry](https://sentry.io/) to assist with error monitoring and performance tracing. Configure the React app to use Sentry as follows:

| EnvVar                       | Type   | Default | Description                                                                                       |
|------------------------------|--------|---------|---------------------------------------------------------------------------------------------------|
| REACT_APP_SENTRY_DSN         | string |         | The DSN for the Sentry project. If not set then Sentry is considered disabled.                    |
| REACT_APP_SENTRY_ENVIRONMENT | string |         | The environment to report (e.g. development, staging, production). Required if Sentry is enabled. |
| REACT_APP_SENTRY_EVENT_ID    | string |         |                                                                                                   |

Sentry is considered **enabled** if a DSN is configured. If Sentry is enabled, an environment is strongly suggested, otherwise the `NODE_ENV` environment will be used.

# Development

{{< hint danger >}}
**Keep up to Date!**<br />
It is essential that we keep this configuration documentation up to date. The devops team uses it to ensure its services are configured correctly. Any time a configuration is changed ensure this documentation is also updated!
{{< /hint >}}

TODO: this section will discuss envconfig, how to interpret environment variables from the configuration struct, how to test configuration, and how to add and change configuration variables. This section should also discuss dotenv files, docker compose, and all of the places where configuration can be influenced (e.g. GitHub actions for React builds).