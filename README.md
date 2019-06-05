# packer-s3-provisioner

packer-s3-provisioner is a Packer provisioner like `file` that can fetch data
from an S3 bucket

## Install

go build

## Usage

## Configuration Options

### Required options

- `bucket` -- the S3 bucket to fetch from
- `key` -- the path in the bucket of the file to fetch
- `local_path` the local path on the provisioning machine to store the file
