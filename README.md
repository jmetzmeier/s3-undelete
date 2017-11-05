# s3-undelete
Restore deleted objects from s3. This is useful in the case of accidental mass deletions

# Usage

s3-undelete --prefix some/prefix mybucket

# AWS Credentials

AWS credentials and region can be set in the same way as awscli. See the official AWS documentation: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

# Notes

This script relies on S3 versioning, which needs to have been enabled at the time the objects were deleted
