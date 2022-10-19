# Dynamic Builds
ARG BUILDER_IMAGE=node:16
ARG FINAL_IMAGE=nginx:stable

# Build stage
FROM ${BUILDER_IMAGE} as builder

WORKDIR /usr/src/app

# Node and Yarn versions
RUN corepack enable

# Dependencies
COPY web/ensign-landing-page/package.json ./
COPY web/ensign-landing-page/yarn.lock ./
RUN yarn install --immutable

# Set production environment variable for build context
ARG REACT_APP_TENANT_BASE_URL="https://api.rotational.app/v1/"
ARG REACT_APP_QUARTERDECK_BASE_URL="https://auth.rotational.app/v1/"
ARG REACT_APP_ANALYTICS_ID
ARG REACT_APP_VERSION_NUMBER
ARG REACT_APP_GIT_REVISION
ARG REACT_APP_SENTRY_DSN
ARG REACT_APP_SENTRY_ENVIRONMENT
ARG REACT_APP_USE_DASH_LOCALE

# Build app with react-scripts
COPY web/ensign-landing-page/ ./
RUN yarn build

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="Rotational Labs <support@rotational.io>"
LABEL description="Ensign Platform User Interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/landing-page/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]