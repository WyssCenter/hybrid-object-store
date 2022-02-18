# TLS Configuration
The Hoss supports internal TLS termination via Let's Encrypt and [Traefik](https://doc.traefik.io/traefik/), which is used as the ingress proxy to all services.

## No TLS Enabled
By default, the system is configured to run without TLS. This is primarily for development and testing configurations. If you don't use TLS, in theory a malicious actor on your network could sniff credentials since they are sent in request headers to the Hoss API.

## Let's Encrypt
You can easily enable TLS via Let's Encrypt. When running in this configuration you must ensure that your server is reachable from the public internet on port 80. Traefik will be configured to redirect all requests, except for the ACME challenge used to validate the cert, to https/443. This means you could lock down the server on port 443 if desired.

To enable Let's Encrypt, in `~/.hoss/.env`, set:

```
LETS_ENCRYPT_ENABLED=true
EXTERNAL_HOSTNAME=https://<HOSTNAME>
```

Since you are now running on https, any other setting that include a route to the server (e.g if you are setting an auth service endpoint) must include the `https://` scheme.


## AWS ALB + Certificate Manager

If you are deploying in AWS, it is possible to place the server behind an ALB and use AWS Certificate Manager to provision certificates. Describing that process is out of scope for this documentation, but if you have successfully configured that, then in `~/.hoss/.env` you must set:

```
LETS_ENCRYPT_ENABLED=false
EXTERNAL_HOSTNAME=https://<HOSTNAME>
HEALTH_CHECK_HOST=<PRIVATE_IP>
```

By setting `HEALTH_CHECK_HOST` to the private IP address of the server, an additional route will be added to Traefik so that health checks from the ALB will succeed and the host will become healthy.

In your ALB's target group, you should use `/core/v1/discover` on port 80, expecting a 200.