FROM golang:alpine AS builder

# Add label to make it easier to remove image
LABEL stage=builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /go/src/github.com/gigantum/libs
COPY ./libs ./

WORKDIR /go/src/github.com/gigantum/hoss-sync
COPY ./sync ./
RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine AS final
LABEL maintainer="Gigantum <support@gigantum.com>"

ARG STAGE=dev
ENV STAGE=$STAGE

# Run container as non-root user
WORKDIR /opt/hoss-sync
RUN adduser -D -u 1001 gig
RUN chown -R gig:gig /opt/hoss-sync
USER 1001

COPY --from=builder /go/src/github.com/gigantum/hoss-sync/hoss-sync .
CMD [ "./hoss-sync" ]
