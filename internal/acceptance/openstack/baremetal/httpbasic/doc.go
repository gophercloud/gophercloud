package httpbasic

/*
Acceptance tests for Ironic endpoints with auth_strategy=noauth.  Specify
IRONIC_ENDPOINT, OS_USERNAME and OS_PASSWORD environment variables.  For example:

  IRONIC_ENDPOINT="http://127.0.0.1:6385/v1"
  OS_USERNAME="myUser"
  OS_PASSWORD="myPassword"
  go test ./acceptance/openstack/baremetal/httpbasic/...

*/
