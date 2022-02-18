#!/bin/sh
set -e

chown -R gig:gig /secrets
exec su-exec gig:gig "$@"
