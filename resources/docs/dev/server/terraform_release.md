# Terraform Release Process

The terraform module ([https://github.com/WyssCenter/terraform-hoss-aws](https://github.com/WyssCenter/terraform-hoss-aws)) is available to make it easier to deploy and manage required AWS infrastructure if deploying with S3 and an AWS EC2 instance for the server.

Typically users will just include the git repo as the source of the module in their terraform. To cut a new release, simply push to `main`.

It is best practice to make sure the README is up-to-date, especially if changing the input or output. The release feature of GitHub can also be used to include a changelog for users, indicating what has been updated.