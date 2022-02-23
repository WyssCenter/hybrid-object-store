# Get current commit of the repo. This is used when building core service to encode in version info
COMMIT=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

# Get the semantic version of the server. This is used when building core service to encode in version info and
# to tag core, auth, and sync services.
VERSION=$(shell cat ./VERSION)

# Setup variables
CURRENT_DIR = $(shell pwd)
HOSS_CONFIG_DIR = ${HOME}/.hoss
ENV_FILE = ${HOSS_CONFIG_DIR}/.env
SHELL := /bin/bash
UID := $(shell id -u)
GID := $(shell id -g)

# Set Docker Compose Files and options
REQUIRED_COMPOSE := -f ./docker-compose.yaml -f ./ingress/docker-compose.yaml -f ./core/docker-compose.yaml -f ./ui/docker-compose.yaml -f ./opensearch/docker-compose.yaml
REQUIRED_COMPOSE_TLS := -f ./core/docker-compose-tls.yaml -f ./ui/docker-compose-tls.yaml
AUTH_COMPOSE := -f ./auth/docker-compose.yaml
AUTH_COMPOSE_TLS := -f ./auth/docker-compose-tls.yaml
DEX_COMPOSE := -f ./dex/docker-compose.yaml
DEX_RECAPTCHA_COMPOSE := -f ./dex/docker-compose-recaptcha.yaml
DEX_COMPOSE_TLS := -f ./dex/docker-compose-tls.yaml
LDAP_COMPOSE := -f ./ldap/docker-compose.yaml
SYNC_COMPOSE := -f ./sync/docker-compose.yaml
MINIO_COMPOSE := -f ./datastore/docker-compose.yaml
MINIO_COMPOSE_TLS := -f ./datastore/docker-compose-tls.yaml
AUTH_DEP_COMPOSE := -f ./options/auth-minio-dependency.yaml
NO_MINIO_COMPOSE := -f ./options/no-minio-redirect.yaml
HEALTH_CHECK_COMPOSE := -f ./options/health-check.yaml

TEST_COMPOSE := $(REQUIRED_COMPOSE) $(AUTH_COMPOSE) $(DEX_COMPOSE) $(LDAP_COMPOSE) $(SYNC_COMPOSE) $(MINIO_COMPOSE) -f ./core/docker-compose-testing.yaml -f ./sync/docker-compose-testing.yaml
TEST_SERVICES := minio etcd-0 reverse-proxy db rabbitmq

# This is used in the `up` target to know if the user
# wants to run the system in the foreground or background
ifndef DETACH
override DETACH:=
else
override DETACH:=-d
endif

# This is used in the `watch-logs` and `get-logs` targets
# to interact with just a single service.
ifndef SERVICE_NAME
override SERVICE_NAME:=
endif

## Load the env file if it exists
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
endif

## Build set of compose files based on enabled services and if TLS is enabled
COMPOSE_FILES := $(REQUIRED_COMPOSE)
ifeq ($(LETS_ENCRYPT_ENABLED),true)
    # TLS is enabled
	COMPOSE_FILES += $(REQUIRED_COMPOSE_TLS)

	ifneq (,$(findstring auth,$(SERVICES)))
		# Auth Service is enabled
		COMPOSE_FILES += $(AUTH_COMPOSE)
		COMPOSE_FILES += $(AUTH_COMPOSE_TLS)
	endif

	ifneq (,$(findstring dex,$(SERVICES)))
		# Dex Service is enabled
		ifneq ($(RECAPTCHA_SITE_KEY),)
			# Dex Service with the reCAPTCHA feature is enabled
			COMPOSE_FILES += $(DEX_RECAPTCHA_COMPOSE)
		else
			COMPOSE_FILES += $(DEX_COMPOSE)
		endif
		COMPOSE_FILES += $(DEX_COMPOSE_TLS)
	endif

	ifneq (,$(findstring ldap,$(SERVICES)))
		# LDAP Service is enabled
		COMPOSE_FILES += $(LDAP_COMPOSE)
	endif

	ifneq (,$(findstring minio,$(SERVICES)))
		# Minio is enabled
		COMPOSE_FILES += $(MINIO_COMPOSE)
		COMPOSE_FILES += $(MINIO_COMPOSE_TLS)
	endif

	ifneq (,$(findstring sync,$(SERVICES)))
		# Sync Service is enabled
		COMPOSE_FILES += $(SYNC_COMPOSE)
	endif

else
    # TLS is NOT enabled
	ifneq (,$(findstring auth,$(SERVICES)))
		# Auth Service is enabled
		COMPOSE_FILES += $(AUTH_COMPOSE)
	endif

	ifneq (,$(findstring dex,$(SERVICES)))
		# Dex Service is enabled
		ifneq ($(RECAPTCHA_SITE_KEY),)
			# Dex Service with the reCAPTCHA feature is enabled
			COMPOSE_FILES += $(DEX_RECAPTCHA_COMPOSE)
		else
			COMPOSE_FILES += $(DEX_COMPOSE)
		endif
	endif

	ifneq (,$(findstring ldap,$(SERVICES)))
		# LDAP Service is enabled
		COMPOSE_FILES += $(LDAP_COMPOSE)
	endif

	ifneq (,$(findstring minio,$(SERVICES)))
		# Minio is enabled
		COMPOSE_FILES += $(MINIO_COMPOSE)
	endif

	ifneq (,$(findstring sync,$(SERVICES)))
		# Sync Service is enabled
		COMPOSE_FILES += $(SYNC_COMPOSE)
	endif

endif

ifneq (,$(findstring auth,$(SERVICES)))
	ifneq (,$(findstring minio,$(SERVICES)))
		# Both auth and minio services are included
		COMPOSE_FILES += $(AUTH_DEP_COMPOSE)
	endif
endif

ifeq (,$(findstring minio,$(SERVICES)))
	# minio service is NOT enabled
	COMPOSE_FILES += $(NO_MINIO_COMPOSE)
endif

ifneq ($(HEALTH_CHECK_HOST),)
	# HEALTH_CHECK_HOST is set
	COMPOSE_FILES += $(HEALTH_CHECK_COMPOSE)
endif


.PHONY: env config

env:
	@echo ""
ifneq ("$(wildcard $(ENV_FILE))","")
	@echo "Using existing .env file"
else
	@echo "Creating new secrets and populating .env file at $(ENV_FILE)"
	@mkdir -p ${HOSS_CONFIG_DIR}
	@printf "SERVICES=opensearch ldap rabbitmq db dex reverse-proxy auth ldap-admin etcd-0 minio core ui sync\n" >> $(ENV_FILE)
	@printf "HOSS_DIR=${HOSS_CONFIG_DIR}\n" >> $(ENV_FILE)
	@printf "LETS_ENCRYPT_ENABLED=false\n" >> $(ENV_FILE)
	@printf "EXTERNAL_HOSTNAME=http://localhost\n" >> $(ENV_FILE)
	@printf "DOMAIN=localhost\n" >> $(ENV_FILE)
	@printf "AUTH_SERVICE_ENDPOINT=http://auth:8080/v1\n" >> $(ENV_FILE)
	@printf "AUTH_CLIENT_SECRET=" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9_' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)
	@printf "HEALTH_CHECK_HOST=\n" >> $(ENV_FILE)
	@printf "ADMIN_EMAIL=\n" >> $(ENV_FILE)
	@printf "NAS_ROOT=${HOSS_CONFIG_DIR}/data/nas\n" >> $(ENV_FILE)
	@printf "MINIO_ROOT_USER=root\n" >> $(ENV_FILE)
	@printf "MINIO_ROOT_PASSWORD=" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9_' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)

	@printf "DB_ROOT=${HOSS_CONFIG_DIR}/data/db\n" >> $(ENV_FILE)
	@printf "POSTGRES_USER=postgres\n" >> $(ENV_FILE)
	@printf "POSTGRES_PASSWORD=" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9_' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)

	@printf "DB_INIT_SCRIPTS=${HOSS_CONFIG_DIR}/core/db-init-scripts\n" >> $(ENV_FILE)
	@printf "POSTGRES_CORE_DB=hoss_core\n" >> $(ENV_FILE)
	@printf "POSTGRES_AUTH_DB=hoss_auth\n" >> $(ENV_FILE)

	@printf "RABBITMQ_COOKIE=" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9_' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)
	@printf "RABBITMQ_USER=gigantum\n" >> $(ENV_FILE)
	@printf "RABBITMQ_PASS=" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)

	@# Note, the 'hsvc_' prefix is used to identify a service access token
	@printf "SERVICE_AUTH_SECRET=hsvc_" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9_' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)
	
	@printf "LDAP_ADMIN_PASSWORD=" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)	
	@printf "LDAP_READONLY_PASSWORD=" >> $(ENV_FILE)
	@LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 48 >> $(ENV_FILE)
	@printf "\n" >> $(ENV_FILE)
	@printf "LDAP_ORGANISATION=Gigantum Inc.\n" >> $(ENV_FILE)
	@printf "LDAP_DOMAIN=hoss.gigantum.com\n" >> $(ENV_FILE)

	@printf "UID=${UID}\n" >> $(ENV_FILE)
	@printf "GID=${GID}\n" >> $(ENV_FILE)
	@printf "BACKUP_ROOT=${HOSS_CONFIG_DIR}/backup\n" >> $(ENV_FILE)
	@printf "RECAPTCHA_SITE_KEY=\n" >> $(ENV_FILE)
	@printf "RECAPTCHA_SECRET_KEY=\n" >> $(ENV_FILE)
	@printf "UI_REDIRECT_REGEX=\n" >> $(ENV_FILE)

	@echo ""
	@echo "By default the system will run on localhost. If you wish to use a FQDN you must edit the EXTERNAL_HOSTNAME and DOMAIN variables in $(ENV_FILE) before running 'make config'"
	@echo "If you wish to use Let's Encrypt to perform TLS termination at the server, the EXTERNAL_HOSTNAME variable must start with 'https://' and LETS_ENCRYPT_ENABLED=true before running 'make config'"
	@echo ""
	@echo ""

endif
ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/core/aws_credentials)","")
	@echo "Using existing core credentials file"
else
	@echo "Creating empty aws credentials file at ${HOSS_CONFIG_DIR}/core/aws_credentials. Remember to edit the file with actual credentials if needed."
	@mkdir -p  ${HOSS_CONFIG_DIR}/core
	@touch ${HOSS_CONFIG_DIR}/core/aws_credentials
endif
ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/core/config.yaml)","")
	@echo "Using existing core configuration"
else
	@echo "Creating default core configuration at ${HOSS_CONFIG_DIR}/core/config.yaml."
	@mkdir -p  ${HOSS_CONFIG_DIR}/core
	@cp ./core/config.yaml ${HOSS_CONFIG_DIR}/core/config.yaml
endif

	@# Copying the db-init-scripts
	@cp -R ./core/db-init-scripts ${HOSS_CONFIG_DIR}/core/

ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/auth/config.yaml)","")
	@echo "Using existing auth configuration"
else
	@echo "Creating default auth configuration at ${HOSS_CONFIG_DIR}/auth/config.yaml."
	@mkdir -p  ${HOSS_CONFIG_DIR}/auth
	@cp ./auth/config.yaml ${HOSS_CONFIG_DIR}/auth/config.yaml
endif
ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/opensearch/opensearch.yaml)","")
	@echo "Using existing opensearch config"
else
	@echo "Creating default opensearch config at ${HOSS_CONFIG_DIR}/opensearch/opensearch.yaml."
	@mkdir -p  ${HOSS_CONFIG_DIR}/opensearch
	@cp ./opensearch/opensearch.yaml ${HOSS_CONFIG_DIR}/opensearch/opensearch.yaml
endif
ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/sync/aws_credentials)","")
	@echo "Using existing sync aws credentials file"
else
	@echo "Creating empty aws credentials file at ${HOSS_CONFIG_DIR}/sync/aws_credentials. Remember to edit the file with actual credentials if needed."
	@mkdir -p  ${HOSS_CONFIG_DIR}/sync
	@touch ${HOSS_CONFIG_DIR}/sync/aws_credentials
endif
ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/sync/config.yaml)","")
	@echo "Using existing sync configuration"
else
	@echo "Creating default sync configuration at ${HOSS_CONFIG_DIR}/sync/config.yaml."
	@mkdir -p  ${HOSS_CONFIG_DIR}/sync
	@cp ./sync/config.yaml ${HOSS_CONFIG_DIR}/sync/config.yaml
endif

	@mkdir -p ${HOSS_CONFIG_DIR}/ui

	@echo ""
	@echo ""

setup:
	docker network create web

config:
	@echo ""
	@echo " - Setting external hostname to '$(EXTERNAL_HOSTNAME)'"

	@if [[ $(LETS_ENCRYPT_ENABLED) == "true" ]]; then\
        echo " - Configuring ingress with TLS termination via Let's Encrypt.";\
		sed "s!{{email}}!$(ADMIN_EMAIL)!" ingress/traefik-tls.yaml.tmpl > ${HOSS_CONFIG_DIR}/traefik.yaml; \
	else \
        echo " - Configuring ingress without TLS termination.";\
		cp ./ingress/traefik.yaml ${HOSS_CONFIG_DIR}/traefik.yaml; \
	fi

	@# The dex/web dir sets the UI customization for dex
	@cp -R dex/web ${HOSS_CONFIG_DIR}/auth/

ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/auth/config-dex.yaml)","")
	@echo " - Using existing OIDC secret and Dex configuration. If you wish to regenerate this file, remove ${HOSS_CONFIG_DIR}/auth/config-dex.yaml and run 'make config' again."
else
	@# We need to make multiple replacements to the dex config file. Due to incompatibilities between sed on macOS and linux,
	@# temporary files are used and then removed instead of just doing in place sed operations.
	@sed "s!{{hostname}}!$(EXTERNAL_HOSTNAME)!" dex/config-dex.yaml.tmpl > ${HOSS_CONFIG_DIR}/auth/config-dex.yaml
	@sed "s!{{ldap_readonly_password}}!$(LDAP_READONLY_PASSWORD)!" ${HOSS_CONFIG_DIR}/auth/config-dex.yaml > ${HOSS_CONFIG_DIR}/auth/config-dex2.yaml

	@$(eval LDAP_DOMAIN=$(shell echo $(LDAP_DOMAIN)) | sed s/\\./,dc=/g)
	@$(eval LDAP_BASE_DN=$(shell echo dc=$(LDAP_DOMAIN)))
	@echo " - Setting LDAP Base DN to '$(LDAP_BASE_DN)'. Modify 'LDAP_DOMAIN' in your .env file and re-run config if you wish to change this."	
	@sed "s!{{ldap_base_dn}}!$(shell echo $(LDAP_BASE_DN))!" ${HOSS_CONFIG_DIR}/auth/config-dex2.yaml > ${HOSS_CONFIG_DIR}/auth/config-dex3.yaml
	
	@echo " - Setting OIDC client secret in Auth service and Dex configs"
	@sed "s!{{auth_client_secret}}!$(shell echo $(AUTH_CLIENT_SECRET))!" ${HOSS_CONFIG_DIR}/auth/config.yaml > ${HOSS_CONFIG_DIR}/auth/config1.yaml
	@sed "s!{{auth_client_secret}}!$(shell echo $(AUTH_CLIENT_SECRET))!" ${HOSS_CONFIG_DIR}/auth/config-dex3.yaml > ${HOSS_CONFIG_DIR}/auth/config-dex4.yaml

	@# If recapcha is enabled we need to set the correct login html template.
    ifneq ($(RECAPTCHA_SITE_KEY),)
		@echo " - Configuring Dex with RECAPTCHA enabled."	
		@sed "s!{{sitekey}}!$(shell echo $(RECAPTCHA_SITE_KEY))!" dex/config-dex-captcha.tmpl >> ${HOSS_CONFIG_DIR}/auth/config-dex4.yaml
		@cp dex/password-captcha.html ${HOSS_CONFIG_DIR}/auth/web/templates/password.html
    else
		@echo " - Configuring Dex with RECAPTCHA disabled."	
		@sed "s!{{sitekey}}!$(shell echo $(RECAPTCHA_SITE_KEY))!" dex/config-dex-no-captcha.tmpl >> ${HOSS_CONFIG_DIR}/auth/config-dex4.yaml
		@cp dex/password.html ${HOSS_CONFIG_DIR}/auth/web/templates/password.html
    endif	
		
	@mv ${HOSS_CONFIG_DIR}/auth/config-dex4.yaml ${HOSS_CONFIG_DIR}/auth/config-dex.yaml
	@rm ${HOSS_CONFIG_DIR}/auth/config-dex3.yaml	
	@rm ${HOSS_CONFIG_DIR}/auth/config-dex2.yaml	
	@rm ${HOSS_CONFIG_DIR}/auth/config.yaml	
	@mv ${HOSS_CONFIG_DIR}/auth/config1.yaml ${HOSS_CONFIG_DIR}/auth/config.yaml
endif

	@# Configure the base DN for the ldap server, if enabled.
	@sed "s!{{ldap_base_dn}}!$(shell echo $(LDAP_BASE_DN))!" ldap/config-ldap.ldif.tmpl > ${HOSS_CONFIG_DIR}/auth/config-ldap.ldif
	@sed "s!{{ldap_base_dn}}!$(shell echo $(LDAP_BASE_DN))!" ldap/env.yaml.tmpl > ${HOSS_CONFIG_DIR}/auth/env.yaml

	@cp auth/openid-config.json.tmpl ${HOSS_CONFIG_DIR}/auth/openid-config.json
	
	@# Create required directories inside the hoss working dir
	@mkdir -p ${HOSS_CONFIG_DIR}/data/db
	@mkdir -p ${HOSS_CONFIG_DIR}/data/events
	@mkdir -p ${HOSS_CONFIG_DIR}/data/nas/data

	@mkdir -p ${HOSS_CONFIG_DIR}/data/opensearch

	@# Make a directory for Dex CA certs if needed
	@mkdir -p ${HOSS_CONFIG_DIR}/auth/certificates

	@mkdir -p ${BACKUP_ROOT}/backups
	@mkdir -p ${BACKUP_ROOT}/.db
	@mkdir -p ${BACKUP_ROOT}/.opensearch

	@# Setting up UI config
ifneq ("$(wildcard ${HOSS_CONFIG_DIR}/ui/config.json)","")
	@echo " - Using existing ui customizations"
else
	@mkdir -p ${HOSS_CONFIG_DIR}/ui
	@sed "s!{{DOMAIN}}!$(DOMAIN)!" ui/config/config.json > ${HOSS_CONFIG_DIR}/ui/config.json
	@cp ui/config/logo.svg ${HOSS_CONFIG_DIR}/ui
	@cp ui/config/favicon.png ${HOSS_CONFIG_DIR}/ui
endif
	
	@# Generate the / to /ui redirect regex. Here we replace '.' with '\.' and make the 's' in 'https'
	@# optional because if running behind an ALB, your external hostname will be https, but requests 
	@# may be http
	@$(eval REDIRECT_DOMAIN=$(shell echo $(DOMAIN)) | sed s/\\./\\\\./g)
	@sed '/UI_REDIRECT_REGEX=/d' $(ENV_FILE) > $(ENV_FILE).tmp
	@echo -n "UI_REDIRECT_REGEX=https?://" >> $(ENV_FILE).tmp
	@echo -n ${REDIRECT_DOMAIN} >> $(ENV_FILE).tmp
	@echo -n "/" >> $(ENV_FILE).tmp
	@mv $(ENV_FILE).tmp $(ENV_FILE)

ifneq ($(UID),1000)
	@echo ""
	@echo "*** You are not running as UID 1000. Run 'make patch' to ensure bind mount permissions are set. ***"
	@echo ""
	@echo "Alternatively, run: "
	@echo " sudo chmod g+rwx ${HOSS_CONFIG_DIR}/data/opensearch"
	@echo " sudo chgrp 1000 ${HOSS_CONFIG_DIR}/data/opensearch"
	@echo " sudo chmod g+rwx ${BACKUP_ROOT}/.opensearch"
	@echo " sudo chgrp 1000 ${BACKUP_ROOT}/.opensearch"
	@echo ""
	@echo ""

endif

patch:
	@# Change the group of the opensearch data directory since it runs as 1000:1000
	@# Recommended solution from: https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#_configuration_files_must_be_readable_by_the_elasticsearch_user
	@sudo chmod g+rwx ${HOSS_CONFIG_DIR}/data/opensearch
	@sudo chgrp 1000 ${HOSS_CONFIG_DIR}/data/opensearch
	@sudo chmod g+rwx ${BACKUP_ROOT}/.opensearch
	@sudo chgrp 1000 ${BACKUP_ROOT}/.opensearch

build:
	DOCKER_BUILDKIT=1 docker-compose --env-file=$(ENV_FILE) $(COMPOSE_FILES) build

restart:
	docker-compose --env-file=$(ENV_FILE) $(COMPOSE_FILES) restart $(SERVICE_NAME)

up:
	@echo ""
	@if [[ $(LETS_ENCRYPT_ENABLED) == "true" ]]; then\
	    echo "TLS is enabled via Let's Encrypt."; \
	else \
	    echo "TLS is disabled."; \
	fi
	@echo ""
	@echo "Starting with the following services: '$(SERVICES)'"
	@echo ""
	@docker-compose --env-file=$(ENV_FILE) $(COMPOSE_FILES) up $(DETACH) $(SERVICES)

up-testing: env config
	@docker-compose --env-file=$(ENV_FILE) $(TEST_COMPOSE) up -d $(TEST_SERVICES)

down:
	@docker-compose --env-file=$(ENV_FILE) $(COMPOSE_FILES) down; \

reset:
	-docker container prune -f
	-docker volume rm server_etcd-0vol
	-docker volume rm server_auth-secrets
	-docker volume rm server_ldap-vol0
	-docker volume rm server_ldap-vol1
	-docker volume rm server_rabbitmq-data
	-rm -rf ${HOSS_CONFIG_DIR}/data
	-rm $(ENV_FILE)

watch-logs:
	@docker-compose --env-file=$(ENV_FILE) $(COMPOSE_FILES) logs -f --tail="10" $(SERVICE_NAME); \

get-logs:
	@docker-compose --env-file=$(ENV_FILE) $(COMPOSE_FILES) logs --no-color --tail="all" $(SERVICE_NAME); \

api-docs:
	@cd core && swag init
	@cd auth && swag init