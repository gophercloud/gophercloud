package testing

import (
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
)

var (
	SegmentID1 = "a5e3a494-26ee-4fde-ad26-2d846c47072e"
	SegmentID2 = "e75ff709-fd25-47f8-ad89-aa6404d74fa8"

	Segment1 = segments.Segment{
		ID:              "a5e3a494-26ee-4fde-ad26-2d846c47072e",
		NetworkID:       "f35f300f-8c26-4d51-a0b0-84e2fff14307",
		Name:            "seg1",
		Description:     "desc",
		PhysicalNetwork: "public",
		NetworkType:     "flat",
		SegmentationID:  0,
		RevisionNumber:  1,
		CreatedAt:       time.Date(2025, 6, 13, 10, 36, 52, 0, time.UTC),
		UpdatedAt:       time.Date(2025, 6, 13, 10, 36, 52, 0, time.UTC),
	}

	Segment2 = segments.Segment{
		ID:              "e75ff709-fd25-47f8-ad89-aa6404d74fa8",
		NetworkID:       "7a7f34de-4dcc-4211-85da-a3afefc8f990",
		Name:            "",
		Description:     "",
		PhysicalNetwork: "",
		NetworkType:     "geneve",
		SegmentationID:  35745,
		RevisionNumber:  0,
		CreatedAt:       time.Date(2025, 6, 13, 10, 36, 45, 0, time.UTC),
		UpdatedAt:       time.Date(2025, 6, 13, 10, 36, 45, 0, time.UTC),
	}

	SegmentsListBody = `
{
  "segments": [
    {
      "id": "a5e3a494-26ee-4fde-ad26-2d846c47072e",
      "network_id": "f35f300f-8c26-4d51-a0b0-84e2fff14307",
      "name": "seg1",
      "description": "desc",
      "physical_network": "public",
      "network_type": "flat",
      "segmentation_id": null,
      "created_at": "2025-06-13T10:36:52Z",
      "updated_at": "2025-06-13T10:36:52Z",
      "revision_number": 1
    },
    {
      "id": "e75ff709-fd25-47f8-ad89-aa6404d74fa8",
      "network_id": "7a7f34de-4dcc-4211-85da-a3afefc8f990",
      "name": null,
      "description": null,
      "physical_network": null,
      "network_type": "geneve",
      "segmentation_id": 35745,
      "created_at": "2025-06-13T10:36:45Z",
      "updated_at": "2025-06-13T10:36:45Z",
      "revision_number": 0
    }
  ]
}`

	createRequest = `{
  "segment": {
    "network_id":"f35f300f-8c26-4d51-a0b0-84e2fff14307",
    "network_type":"flat",
    "physical_network":"public",
    "name":"seg1",
    "description":"desc"
  }
}`

	createResponse = `{
  "segment": {
    "id":"a5e3a494-26ee-4fde-ad26-2d846c47072e",
    "network_id":"f35f300f-8c26-4d51-a0b0-84e2fff14307",
    "name":"seg1",
    "description":"desc",
    "physical_network":"public",
    "network_type":"flat",
    "segmentation_id":null,
    "created_at": "2025-06-13T10:36:52Z",
    "updated_at": "2025-06-13T10:36:52Z",
    "revision_number":1
  }
}
`

	updateRequest = `{
  "segment": {
    "name":"new-name",
    "description":"new-desc"
  }
}`

	updateResponse = `{
  "segment": {
    "id":"a5e3a494-26ee-4fde-ad26-2d846c47072e",
    "network_id":"f35f300f-8c26-4d51-a0b0-84e2fff14307",
    "name":"new-name",
    "description":"new-desc",
    "physical_network":"public",
    "network_type":"flat",
    "segmentation_id":null,
    "created_at": "2025-06-13T10:36:52Z",
    "updated_at": "2025-06-13T10:36:52Z",
    "revision_number":1
  }
}`
)
