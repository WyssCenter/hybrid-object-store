# Hoss Goals and Design Motivations

Object storage for data science applications has very appealing features:

- Durable
- Scalable both in storage capacity and IO
- Object-level metadata
- Storage tiering

But it has many challenges:

- Permissions and access control systems (IAM) are very complex
- Large learning curve for data scientists
- Credential management is complicated
- No good user interface for novices
- Metadata is not easily searched or discovered
- Typically cloud only (e.g. S3)
- Egress costs can add up without caching


The Hoss was built to address these challenges and bring as many benefits to data science applications as possible. The goal is to build a system that:

- Is simple to administer (e.g. no kubernetes!)
- Makes it easy for users of diverse skill levels work together on data in an object store
- Makes it easy to search and discover data & metadata
- Can be made highly available
  - Even if not supported yet, architecturally should be compatible with an HA deployment
- Can scale to high IO throughput

This resulted in a system that:
- Is built on Docker and Docker Compose
- Is composed of open source tools and custom services
  - Custom services are broken apart based on how they would need to deploy/scale, with a preference on minimizing the number of services
- Abstracts and automates IAM policies
- Manages temporary IAM credentials
- Integrates with existing auth systems, including enterprise systems like Azure AD
- Provides the ability to sync data between servers to build "hybrid cloud" workflows
- Provides a web-based UI and file browser for "easy" use by novices and experts alike
- Provides a Python-based client library or REST API for skilled users