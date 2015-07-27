package v2

// ImageStatus image statuses
//http://docs.openstack.org/developer/glance/statuses.html
type ImageStatus string

const (
	// ImageStatusQueued The image identifier has been reserved for an
	// image in the Glance registry. No image data has been uploaded to Glance
	// and the image size was not explicitly set to zero on creation.
	ImageStatusQueued ImageStatus = "queued"

	// ImageStatusActive Denotes an image that is fully available in Glance.
	// This occurs when the image data is uploaded, or the image size is
	// explicitly set to zero on creation.
	ImageStatusActive ImageStatus = "active"
	// TODO not all image statuses are defined
)

// ImageVisibility denotes an image that is fully available in Glance.
// This occurs when the image data is uploaded, or the image size
// is explicitly set to zero on creation.
// According to design
// https://wiki.openstack.org/wiki/Glance-v2-community-image-visibility-design
type ImageVisibility string

const (
	// ImageVisibilityPublic all users
	ImageVisibilityPublic ImageVisibility = "public"
	// ImageVisibilityPrivate users with tenantId == tenantId(owner)
	ImageVisibilityPrivate ImageVisibility = "private"
)
