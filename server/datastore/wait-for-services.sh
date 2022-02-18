#!/usr/bin/env bash

attempt_counter=0
max_attempts=$2

echo "Waiting for Auth Service at ${1}" 
until $(curl --output /dev/null --silent --fail $1/.well-known/openid-configuration); do
    if [ ${attempt_counter} -eq ${max_attempts} ];then
      echo "Max attempts reached"
      exit 1
    fi

    echo 'Not ready. Trying again in 5 seconds...'
    attempt_counter=$(($attempt_counter+1))
    sleep 5
done
echo "Auth service is ready!"

echo "Waiting for RabbitMQ Service" 
attempt_counter=0
until $(curl --output /dev/null --silent rabbitmq:15672/api/overview); do
    if [ ${attempt_counter} -eq ${max_attempts} ];then
      echo "Max attempts reached"
      exit 1
    fi

    echo 'Not ready. Trying again in 5 seconds...'
    attempt_counter=$(($attempt_counter+1))
    sleep 5
done

# Exec the normal entrypoint if you get here.
sleep 3
echo "RabbitMQ service is ready! Starting minio."
exec /usr/bin/docker-entrypoint.sh "${@:4}"
