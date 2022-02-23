.. Hoss documentation master file, created by
   sphinx-quickstart on Thu Jan 13 14:59:59 2022.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Hoss Server Documentation
===============================================

The Hoss is a collection of web services and applications that *manage* S3 API compatible object stores 
(e.g. `AWS S3 <https://aws.amazon.com/s3/>`_, `minIO <https://min.io/>`_), with a focus
on data science application running on hybrid cloud architectures. It is designed to be simple
to deploy and use, with a minimal infrastructure footprint. A Hoss server:

* Provides minimal abstractions to organize data
* Automates IAM policy management
* Integrates with external authentication providers
* Automates data synchronization between buckets and sites
* Provides metadata search
* Provides familiar programming interfaces via the `hoss-client library <https://github.com/gigantum/hoss-client>`_
* Provides a web-based UI and API + Python library to meet a broad range of skill-levels

The documentation that follows provide details and instructions on how to
install, configure, update, backup, and maintain a Hoss server installation.

.. toctree::
   :maxdepth: 2
   :hidden:
   :caption: Admin - Installation

   installation/overview
   installation/prepare
   installation/install-on-prem
   installation/install-aws
   installation/install-server-pair


.. toctree::
   :maxdepth: 2
   :hidden:
   :caption: Admin - Configuration

   configuration/env-vars
   configuration/core
   configuration/auth
   configuration/sync
   configuration/ui
   configuration/tls
   configuration/captcha
   configuration/versioning

.. toctree::
   :maxdepth: 2
   :hidden:
   :caption: Admin - Maintenance

   maintenance/internal-ldap
   maintenance/backup-and-restore
   maintenance/update
   maintenance/monitor-logs
   maintenance/roles
   maintenance/revoking-access

