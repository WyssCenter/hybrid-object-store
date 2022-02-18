# FRONTEND BUILD CONTAINER
FROM ubuntu:18.04 as ui-builder
LABEL stage=ui-builder

# Install system level dependencies
RUN apt-get -y update && apt-get -y install curl git wget
RUN curl -sL https://deb.nodesource.com/setup_14.x |  bash -
RUN apt-get -y install nodejs

RUN mkdir /opt/ui
COPY ./ui/package.json /opt/ui

# Install global requirements for frontend build
WORKDIR /opt/ui

# Install nodejs packages
RUN npm install -g yarn && \
    yarn install

# Build frontend, including static pre-rendering
COPY ui /opt/ui
RUN yarn build




# GO BUILD CONTAINER
FROM golang:alpine AS go-builder

# Add label to make it easier to remove image
LABEL stage=go-builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash openssh wget make build-base

WORKDIR /go/src/github.com/gigantum/hoss-ui
COPY . ./
RUN make build




# Render Service Image
FROM alpine

WORKDIR /renderer
RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates \
    && addgroup -S renderservice && adduser -S renderservice -G renderservice \
    && update-ca-certificates 2>/dev/null || true \
    && rm -rf /var/cache/apk/*

# Copy compiled frontend from build container
COPY --from=ui-builder /opt/ui/build ui/build

COPY --from=go-builder /go/src/github.com/gigantum/hoss-ui/render .

RUN chown -R renderservice:renderservice /renderer

USER renderservice

ENTRYPOINT [ "/renderer/render" ]
