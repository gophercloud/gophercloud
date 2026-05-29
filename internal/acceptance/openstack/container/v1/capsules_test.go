//go:build acceptance || containers || capsules

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/container/v1/capsules"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCapsuleBase(t *testing.T) {
	t.Skip("Currently failing in OpenLab")

	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")

	client, err := clients.NewContainerV1Client()
	th.AssertNoErr(t, err)

	template := new(capsules.Template)
	template.Bin = []byte(capsuleTemplate)

	createOpts := capsules.CreateOpts{
		TemplateOpts: template,
	}

	v, err := capsules.Create(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)
	capsule := v.(*capsules.Capsule)

	err = WaitForCapsuleStatus(client, capsule.UUID, "Running")
	th.AssertNoErr(t, err)

	pager := capsules.List(client, nil)
	err = pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		v, err := capsules.ExtractCapsules(page)
		th.AssertNoErr(t, err)
		allCapsules := v.([]capsules.Capsule)

		for _, m := range allCapsules {
			capsuleUUID := m.UUID
			if capsuleUUID != capsule.UUID {
				continue
			}
			capsule, err := capsules.Get(context.TODO(), client, capsuleUUID).ExtractBase()

			th.AssertNoErr(t, err)
			th.AssertEquals(t, capsule.MetaName, "template")

			err = capsules.Delete(context.TODO(), client, capsuleUUID).ExtractErr()
			th.AssertNoErr(t, err)

		}
		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestCapsuleV132(t *testing.T) {
	t.Skip("Currently failing in OpenLab")

	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.SkipRelease(t, "stable/rocky")
	clients.SkipRelease(t, "stable/stein")

	client, err := clients.NewContainerV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.32"

	template := new(capsules.Template)
	template.Bin = []byte(capsuleTemplate)

	createOpts := capsules.CreateOpts{
		TemplateOpts: template,
	}

	capsule, err := capsules.Create(context.TODO(), client, createOpts).ExtractV132()
	th.AssertNoErr(t, err)

	err = WaitForCapsuleStatus(client, capsule.UUID, "Running")
	th.AssertNoErr(t, err)

	pager := capsules.List(client, nil)
	err = pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		allCapsules, err := capsules.ExtractCapsulesV132(page)
		th.AssertNoErr(t, err)

		for _, m := range allCapsules {
			capsuleUUID := m.UUID
			if capsuleUUID != capsule.UUID {
				continue
			}
			capsule, err := capsules.Get(context.TODO(), client, capsuleUUID).ExtractV132()

			th.AssertNoErr(t, err)
			th.AssertEquals(t, capsule.MetaName, "template")

			err = capsules.Delete(context.TODO(), client, capsuleUUID).ExtractErr()
			th.AssertNoErr(t, err)

		}
		return true, nil
	})
	th.AssertNoErr(t, err)
}
