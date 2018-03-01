package compute

import (
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const WaitForVolumeAttachmentDeletePollInterval = time.Duration(1 * time.Second)
const WaitForVolumeAttachmentDeleteTimeout = time.Duration(30 * time.Second)
const WaitForVolumeAttachmentReadyPollInterval = time.Duration(1 * time.Second)
const WaitForVolumeAttachmentReadyTimeout = time.Duration(30 * time.Second)

// StorageAttachmentsClient is a client for the Storage Attachment functions of the Compute API.
type StorageAttachmentsClient struct {
	ResourceClient
}

// StorageAttachments obtains a StorageAttachmentsClient which can be used to access to the
// Storage Attachment functions of the Compute API
func (c *ComputeClient) StorageAttachments() *StorageAttachmentsClient {
	return &StorageAttachmentsClient{
		ResourceClient: ResourceClient{
			ComputeClient:       c,
			ResourceDescription: "storage volume attachment",
			ContainerPath:       "/storage/attachment/",
			ResourceRootPath:    "/storage/attachment",
		}}
}

type StorageAttachmentState string

const (
	Attaching   StorageAttachmentState = "attaching"
	Attached    StorageAttachmentState = "attached"
	Detaching   StorageAttachmentState = "detaching"
	Unavailable StorageAttachmentState = "unavailable"
	Unknown     StorageAttachmentState = "unknown"
)

// StorageAttachmentInfo describes an existing storage attachment.
type StorageAttachmentInfo struct {
	// Name of this attachment, generated by the server.
	Name string `json:"name"`

	// Index number for the volume. The allowed range is 1-10
	// An attachment with index 1 is exposed to the instance as /dev/xvdb, an attachment with index 2 is exposed as /dev/xvdc, and so on.
	Index int `json:"index"`

	// Multipart name of the instance attached to the storage volume.
	InstanceName string `json:"instance_name"`

	// Multipart name of the volume attached to the instance.
	StorageVolumeName string `json:"storage_volume_name"`

	// The State of the Storage Attachment
	State StorageAttachmentState `json:"state"`
}

func (c *StorageAttachmentsClient) success(attachmentInfo *StorageAttachmentInfo) (*StorageAttachmentInfo, error) {
	c.unqualify(&attachmentInfo.Name, &attachmentInfo.InstanceName, &attachmentInfo.StorageVolumeName)
	return attachmentInfo, nil
}

type CreateStorageAttachmentInput struct {
	// Index number for the volume. The allowed range is 1-10
	// An attachment with index 1 is exposed to the instance as /dev/xvdb, an attachment with index 2 is exposed as /dev/xvdc, and so on.
	// Required.
	Index int `json:"index"`

	// Multipart name of the instance to which you want to attach the volume.
	// Required.
	InstanceName string `json:"instance_name"`

	// Multipart name of the volume that you want to attach.
	// Required.
	StorageVolumeName string `json:"storage_volume_name"`

	// Time to wait between polls to check volume attachment status
	PollInterval time.Duration `json:"-"`

	// Time to wait for storage volume attachment
	Timeout time.Duration `json:"-"`
}

// CreateStorageAttachment creates a storage attachment attaching the given volume to the given instance at the given index.
func (c *StorageAttachmentsClient) CreateStorageAttachment(input *CreateStorageAttachmentInput) (*StorageAttachmentInfo, error) {
	input.InstanceName = c.getQualifiedName(input.InstanceName)
	input.StorageVolumeName = c.getQualifiedName(input.StorageVolumeName)

	var attachmentInfo *StorageAttachmentInfo
	if err := c.createResource(&input, &attachmentInfo); err != nil {
		return nil, err
	}

	if input.PollInterval == 0 {
		input.PollInterval = WaitForVolumeAttachmentReadyPollInterval
	}
	if input.Timeout == 0 {
		input.Timeout = WaitForVolumeAttachmentReadyTimeout
	}

	return c.waitForStorageAttachmentToFullyAttach(attachmentInfo.Name, input.PollInterval, input.Timeout)
}

// DeleteStorageAttachmentInput represents the body of an API request to delete a Storage Attachment.
type DeleteStorageAttachmentInput struct {
	// The three-part name of the Storage Attachment (/Compute-identity_domain/user/object).
	// Required
	Name string `json:"name"`

	// Time to wait between polls to check volume attachment status
	PollInterval time.Duration `json:"-"`

	// Time to wait for storage volume snapshot
	Timeout time.Duration `json:"-"`
}

// DeleteStorageAttachment deletes the storage attachment with the given name.
func (c *StorageAttachmentsClient) DeleteStorageAttachment(input *DeleteStorageAttachmentInput) error {
	if err := c.deleteResource(input.Name); err != nil {
		return err
	}

	if input.PollInterval == 0 {
		input.PollInterval = WaitForVolumeAttachmentDeletePollInterval
	}
	if input.Timeout == 0 {
		input.Timeout = WaitForVolumeAttachmentDeleteTimeout
	}

	return c.waitForStorageAttachmentToBeDeleted(input.Name, input.PollInterval, input.Timeout)
}

// GetStorageAttachmentInput represents the body of an API request to obtain a Storage Attachment.
type GetStorageAttachmentInput struct {
	// The three-part name of the Storage Attachment (/Compute-identity_domain/user/object).
	// Required
	Name string `json:"name"`
}

// GetStorageAttachment retrieves the storage attachment with the given name.
func (c *StorageAttachmentsClient) GetStorageAttachment(input *GetStorageAttachmentInput) (*StorageAttachmentInfo, error) {
	var attachmentInfo *StorageAttachmentInfo
	if err := c.getResource(input.Name, &attachmentInfo); err != nil {
		return nil, err
	}

	return c.success(attachmentInfo)
}

// waitForStorageAttachmentToFullyAttach waits for the storage attachment with the given name to be fully attached, or times out.
func (c *StorageAttachmentsClient) waitForStorageAttachmentToFullyAttach(name string, pollInterval, timeout time.Duration) (*StorageAttachmentInfo, error) {
	var waitResult *StorageAttachmentInfo

	err := c.client.WaitFor("storage attachment to be attached", pollInterval, timeout, func() (bool, error) {
		input := &GetStorageAttachmentInput{
			Name: name,
		}
		info, err := c.GetStorageAttachment(input)
		if err != nil {
			return false, err
		}

		if info != nil {
			if info.State == Attached {
				waitResult = info
				return true, nil
			}
		}

		return false, nil
	})

	return waitResult, err
}

// waitForStorageAttachmentToBeDeleted waits for the storage attachment with the given name to be fully deleted, or times out.
func (c *StorageAttachmentsClient) waitForStorageAttachmentToBeDeleted(name string, pollInterval, timeout time.Duration) error {
	return c.client.WaitFor("storage attachment to be deleted", pollInterval, timeout, func() (bool, error) {
		input := &GetStorageAttachmentInput{
			Name: name,
		}
		_, err := c.GetStorageAttachment(input)
		if err != nil {
			if client.WasNotFoundError(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}
