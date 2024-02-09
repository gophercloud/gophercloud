package v2

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetransfers"
)

func CreateTransferRequest(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share, name string) (*sharetransfers.Transfer, error) {
	opts := sharetransfers.CreateOpts{
		ShareID: share.ID,
		Name:    name,
	}
	transfer, err := sharetransfers.Create(context.TODO(), client, opts).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to create a share transfer request: %s", err)
	}

	return transfer, nil
}

func AcceptTransfer(t *testing.T, client *gophercloud.ServiceClient, transferRequest *sharetransfers.Transfer) error {
	opts := sharetransfers.AcceptOpts{
		AuthKey:          transferRequest.AuthKey,
		ClearAccessRules: true,
	}
	err := sharetransfers.Accept(context.TODO(), client, transferRequest.ID, opts).ExtractErr()
	if err != nil {
		return fmt.Errorf("failed to accept a share transfer request: %s", err)
	}

	return nil
}

func DeleteTransferRequest(t *testing.T, client *gophercloud.ServiceClient, transfer *sharetransfers.Transfer) {
	err := sharetransfers.Delete(context.TODO(), client, transfer.ID).ExtractErr()
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return
		}
		t.Errorf("Unable to delete share transfer %s: %v", transfer.ID, err)
	}
}
