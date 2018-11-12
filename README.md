# tfgen
[![release](http://img.shields.io/github/release/mschurenko/tfgen.svg?style=flat-square)](https://github.com/mschurenko/tfgen/releases)

`tfgen` generates some Terraform configurations to make your life easier.
All generated files are JSON and will have the file extension `.tf.json`.
Please note that this is not a Terraform wrapper.


## Installation
Download the latest release from [GitHub Releases](https://github.com/mschurenko/tfgen/releases).

## Setup
Create a `.tfgen.yml` file somewhere in your Terraform project's path.
```yaml
---
# terraform s3 remote backend
# this will be used to generate a `backend.tf.json` for each stack
s3_backend:
  aws_region: us-west-2
  bucket: my-terraform-s3-backend-bucket
  dynamodb_table: my-terraform-dynamodb-table

# your stacks must have one of the following as a parent directory
environments:
  - production
  - staging
  - dev
```

`tfgen` will walk up the directory path in your cwd until it finds a valid `.tfgen.yml` file. This allows you to have separate config files for environments, applications, etc. It all depends on the directory structure of your Terraform project.

For example:
```
terraform
    ├── .tfgen.yml
    ├── production
    │   ├── myapp
    │   ├── redis
    │   ├── mysql
    │   ├── alb
    ├── staging
    │   ├── myapp
    │   ├── redis
    │   ├── mysql
    │   ├── alb
```

## Usage
#### Adding a new stack:
```sh
cd production
mkdir new-stack
cd new-stack
tfgen init-stack
```

This will generate a `backend.tf.json` file

#### Adding a `terraform_remote_state` data source:
```
tfgen remote-state production/vpc
```

This will generate a `remote_states.tf.json` file.

Example:
```json
{
  "data": {
    "terraform_remote_state": {
      "production_vpc": {
        "backend": "s3",
        "config": {
          "bucket": "my-terraform-s3-backend-bucket",
          "dynamodb_table": "my-terraform-dynamodb-table",
          "encrypt": true,
          "key": "stacks/production/vpc/terraform.tfstate",
          "region": "us-west-2"
        }
      }
}
```

One can make use of this in a Terraform template like this:
```hcl
locals {
  vpc_id = "${data.terraform_remote_state.production_vpc.vpc_id}"
}
```
