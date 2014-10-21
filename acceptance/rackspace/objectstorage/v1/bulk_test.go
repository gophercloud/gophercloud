// +build acceptance rackspace objectstorage v1

package v1

import (
  "fmt"
  "testing"

 "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/bulk"
  th "github.com/rackspace/gophercloud/testhelper"
)

func TestBulk(t *testing.T){
  c, err := createClient(t, false)
  th.AssertNoErr(t, err)

  options := &bulk.DeleteOpts{"container/object1"}
  res := bulk.Delete(c, options)
  fmt.Printf("res: %+v\n", res)
}
