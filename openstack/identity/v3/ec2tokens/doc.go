/*
Package ec2tokens provides low-level access to Keystone's EC2 and S3 token
endpoints.

For normal Keystone authentication using EC2 credentials, use AuthOptionsEC2
and EC2TokenOpts in github.com/gophercloud/gophercloud/v2/auth. That package
constructs and signs the authentication request independently of these resources.

For more information, see:
https://docs.openstack.org/api-ref/identity/v2-ext/
*/
package ec2tokens
