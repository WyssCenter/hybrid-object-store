FROM golang:alpine AS builder

# Add label to make it easier to remove image
LABEL stage=builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /go/src/github.com/gigantum/libs
COPY ./libs ./

WORKDIR /go/src/github.com/gigantum/hoss-auth
COPY ./auth ./

# Use swaggo to build swagger docs package
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.8
RUN swag init

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine AS final
LABEL maintainer="Gigantum <support@gigantum.com>"

ARG STAGE=dev
ENV STAGE=$STAGE

RUN apk update && apk upgrade && \
    apk add --no-cache su-exec openldap-clients

VOLUME /secrets

# Run container as non-root user
WORKDIR /opt/hoss-auth
RUN adduser -D -u 1001 gig # Entrypoint will clean up permissions on /secrets and su-exec as this user
RUN chown -R gig:gig /opt/hoss-auth

COPY auth/docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

COPY --from=builder /go/src/github.com/gigantum/hoss-auth/hoss-auth .
CMD [ "./hoss-auth" ]
