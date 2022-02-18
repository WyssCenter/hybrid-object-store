FROM golang:alpine AS builder

# Add label to make it easier to remove image
LABEL stage=builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh wget

WORKDIR /go/src/github.com/gigantum/libs
COPY ./libs ./

WORKDIR /go/src/github.com/gigantum/hoss-core
COPY ./core ./

# Use swaggo to build swagger docs package
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.8
RUN swag init

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a



FROM alpine AS final
ARG COMMIT_HASH
ARG VERSION

LABEL maintainer="Gigantum <support@gigantum.com>"

ARG STAGE=dev
ENV STAGE=$STAGE

# Install mc
WORKDIR /usr/local/bin
RUN wget https://dl.min.io/client/mc/release/linux-amd64/mc && \
    chmod +x mc && env

WORKDIR /opt/hoss-core

# Save Discovery Information
RUN echo $COMMIT_HASH > /opt/hoss-core/discover_build_hash
RUN echo $VERSION > /opt/hoss-core/discover_version

# Run container as non-root user
RUN adduser -D -u 1001 gig
RUN chown -R gig:gig /opt/hoss-core
USER 1001

COPY --from=builder /go/src/github.com/gigantum/hoss-core/hoss-core .
CMD [ "./hoss-core" ]
