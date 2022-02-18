# Note: when updating this image, you must also update it in
# server/core/docker-compose-testing.yaml and in .github/workflows/main.yaml
# so that unit tests use the same minio image (and it's not actually built)
# This is needed to prevent the auth service check from running.
FROM minio/minio:RELEASE.2021-11-09T03-21-45Z
ARG AUTH_ENDPOINT="http://auth:8080/v1"
ENV AUTH_ENDPOINT=${AUTH_ENDPOINT}

RUN microdnf install netcat

COPY datastore/wait-for-services.sh /opt/wait-for-services.sh

RUN chmod +x /opt/wait-for-services.sh

# This calls the wait
ENTRYPOINT /opt/wait-for-services.sh $AUTH_ENDPOINT 12 -- $0 $@
