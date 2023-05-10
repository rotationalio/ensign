# Dynamic Builds
ARG BUILDER_IMAGE=node:16
ARG FINAL_IMAGE=nginx:stable

# Build stage
FROM ${BUILDER_IMAGE} AS builder

WORKDIR /usr/src/app

# Install dependencies
COPY web/beacon-app/package*.json ./
COPY web/beacon-app/yarn.lock ./
RUN yarn install --immutable

# Set environment variables
ARG REACT_APP_TENANT_BASE_URL="https://api.rotational.app/v1/"
ARG REACT_APP_QUARTERDECK_BASE_URL="https://auth.rotational.app/v1/"
ARG REACT_APP_ANALYTICS_ID
ARG REACT_APP_VERSION_NUMBER
ARG REACT_APP_GIT_REVISION
ARG REACT_APP_SENTRY_DSN
ARG REACT_APP_SENTRY_ENVIRONMENT="production"
ARG REACT_APP_USE_DASH_LOCALE="false"

# Build app with react-scripts
COPY web/beacon-app/ ./
RUN yarn build

# Final stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="Rotational Labs <support@rotational.io>"
LABEL description="The Ensign user interface"

COPY --from=builder /usr/src/app/build /usr/share/nginx/html
COPY containers/beacon/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD ["nginx", "-g", "daemon off;"]