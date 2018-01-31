package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const WaitForOrchestrationActiveTimeout = time.Duration(3600 * time.Second)
const WaitForOrchestrationDeleteTimeout = time.Duration(3600 * time.Second)

// OrchestrationsClient is a client for the Orchestration functions of the Compute API.
type OrchestrationsClient struct {
	ResourceClient
}

// Orchestrations obtains an OrchestrationsClient which can be used to access to the
// Orchestration functions of the Compute API
func (c *ComputeClient) Orchestrations() *OrchestrationsClient {
	return &OrchestrationsClient{
		ResourceClient: ResourceClient{
			ComputeClient:       c,
			ResourceDescription: "Orchestration",
			ContainerPath:       "/platform/v1/orchestration/",
			ResourceRootPath:    "/platform/v1/orchestration",
		}}
}

type OrchestrationDesiredState string

const (
	// * active: Creates all the orchestration objects defined in the orchestration.
	OrchestrationDesiredStateActive OrchestrationDesiredState = "active"
	// * inactive: Adds the orchestration to Oracle Compute Cloud Service, but does not create any of the orchestration
	OrchestrationDesiredStateInactive OrchestrationDesiredState = "inactive"
	// * suspended: Suspends all orchestration objects defined in the orchestration
	OrchestrationDesiredStateSuspend OrchestrationDesiredState = "suspend"
)

type OrchestrationStatus string

const (
	OrchestrationStatusActive       OrchestrationStatus = "active"
	OrchestrationStatusInactive     OrchestrationStatus = "inactive"
	OrchestrationStatusSuspend      OrchestrationStatus = "suspend"
	OrchestrationStatusActivating   OrchestrationStatus = "activating"
	OrchestrationStatusDeleting     OrchestrationStatus = "deleting"
	OrchestrationStatusError        OrchestrationStatus = "terminal_error"
	OrchestrationStatusStopping     OrchestrationStatus = "stopping"
	OrchestrationStatusSuspending   OrchestrationStatus = "suspending"
	OrchestrationStatusStarting     OrchestrationStatus = "starting"
	OrchestrationStatusDeactivating OrchestrationStatus = "deactivating"
	OrchestrationStatusSuspended    OrchestrationStatus = "suspended"
)

type OrchestrationType string

const (
	OrchestrationTypeInstance OrchestrationType = "Instance"
)

type OrchestrationRelationshipType string

const (
	OrchestrationRelationshipTypeDepends OrchestrationRelationshipType = "depends"
)

// OrchestrationInfo describes an existing Orchestration.
type Orchestration struct {
	// The default Oracle Compute Cloud Service account, such as /Compute-acme/default.
	Account string `json:"account"`
	// Description of this orchestration
	Description string `json:"description"`
	// The desired_state specified in the orchestration JSON file. A unique identifier for this orchestration.
	DesiredState OrchestrationDesiredState `json:"desired_state"`
	// Unique identifier of this orchestration
	ID string `json:"id"`
	// The three-part name of the Orchestration (/Compute-identity_domain/user/object).
	Name string `json:"name"`
	// List of orchestration objects
	Objects []Object `json:"objects"`
	// Current status of this orchestration
	Status OrchestrationStatus `json:"status"`
	// Strings that describe the orchestration and help you identify it.
	Tags []string `json:"tags"`
	// Time the orchestration was last audited
	TimeAudited string `json:"time_audited"`
	// The time when the orchestration was added to Oracle Compute Cloud Service.
	TimeCreated string `json:"time_created"`
	// The time when the orchestration was last updated in Oracle Compute Cloud Service.
	TimeUpdated string `json:"time_updated"`
	// Unique Resource Identifier
	URI string `json:"uri"`
	// Name of the user who added this orchestration or made the most recent update to this orchestration.
	User string `json:"user"`
	// Version of this orchestration. It is automatically generated by the server.
	Version int `json:"version"`
}

// CreateOrchestrationInput defines an Orchestration to be created.
type CreateOrchestrationInput struct {
	// The default Oracle Compute Cloud Service account, such as /Compute-acme/default.
	// Optional
	Account string `json:"account,omitempty"`
	// Description of this orchestration
	// Optional
	Description string `json:"description,omitempty"`
	// Specify the desired state of this orchestration: active, inactive, or suspend.
	// You can manage the state of the orchestration objects by changing the desired state of the orchestration.
	// * active: Creates all the orchestration objects defined in the orchestration.
	// * inactive: Adds the orchestration to Oracle Compute Cloud Service, but does not create any of the orchestration
	// objects defined in the orchestration.
	// Required
	DesiredState OrchestrationDesiredState `json:"desired_state"`
	// The three-part name of the Orchestration (/Compute-identity_domain/user/object).
	// Object names can contain only alphanumeric characters, hyphens, underscores, and periods. Object names are case-sensitive.
	// Required
	Name string `json:"name"`
	// The list of objects in the orchestration. An object is the primary building block of an orchestration.
	// An orchestration can contain up to 100 objects.
	// Required
	Objects []Object `json:"objects"`
	// Strings that describe the orchestration and help you identify it.
	Tags []string `json:"tags,omitempty"`
	// Version of this orchestration. It is automatically generated by the server.
	Version int `json:"version,omitempty"`
	// Time to wait for an orchestration to be ready
	Timeout time.Duration `json:"-"`
}

type Object struct {
	// The default Oracle Compute Cloud Service account, such as /Compute-acme/default.
	// Optional
	Account string `json:"account,omitempty"`
	// Description of this orchestration
	// Optional
	Description string `json:"description,omitempty"`
	// The desired state of the object
	// Optional
	DesiredState OrchestrationDesiredState `json:"desired_state,omitempty"`
	// Dictionary containing the current state of the object
	Health Health `json:"health,omitempty"`
	// A text string describing the object. Labels can't include spaces. In an orchestration, the label for
	// each object must be unique. Maximum length is 256 characters.
	// Required
	Label string `json:"label"`
	// The four-part name of the object (/Compute-identity_domain/user/orchestration/object). If you don't specify a name
	// for this object, the name is generated automatically. Object names can contain only alphanumeric characters, hyphens,
	// underscores, and periods. Object names are case-sensitive. When you specify the object name, ensure that an object of
	// the same type and with the same name doesn't already exist. If such a object already exists, then another
	// object of the same type and with the same name won't be created and the existing object won't be updated.
	// Optional
	Name string `json:"name,omitempty"`
	// The three-part name (/Compute-identity_domain/user/object) of the orchestration to which the object belongs.
	// Required
	Orchestration string `json:"orchestration"`
	// Specifies whether the object should persist when the orchestration is suspended. Specify one of the following:
	// * true: The object persists when the orchestration is suspended.
	// * false: The object is deleted when the orchestration is suspended.
	// By default, persistent is set to false. It is recommended that you specify true for storage
	// volumes and other critical objects. Persistence applies only when you're suspending an orchestration.
	// When you terminate an orchestration, all the objects defined in it are deleted.
	// Optional
	Persistent bool `json:"persistent,omitempty"`
	// The relationship between the objects that are created by this orchestration. The
	// only supported relationship is depends, indicating that the specified target objects must be created first.
	// Note that when recovering from a failure, the orchestration doesn't consider object relationships.
	// Orchestrations v2 use object references to recover interdependent objects to a healthy state. SeeObject
	// References and Relationships in Using Oracle Compute Cloud Service (IaaS).
	Relationships []Relationship `json:"relationships,omitempty"`
	// The template attribute defines the properties or characteristics of the Oracle Compute Cloud Service object
	// that you want to create, as specified by the type attribute.
	// The fields in the template section vary depending on the specified type. See Orchestration v2 Attributes
	// Specific to Each Object Type in Using Oracle Compute Cloud Service (IaaS) to determine the parameters that are
	// specific to each object type that you want to create.
	// For example, if you want to create a storage volume, the type would be StorageVolume, and the template would include
	// size and bootable. If you want to create an instance, the type would be Instance, and the template would include
	// instance-specific attributes, such as imagelist and shape.
	// Required
	Template interface{} `json:"template"`
	// Specify one of the following object types that you want to create.
	// The only allowed type is Instance
	// Required
	Type OrchestrationType `json:"type"`
	// Version of this object, generated by the server
	// Optional
	Version int `json:"version,omitempty"`
}

type Health struct {
	// The status of the object
	Status OrchestrationStatus `json:"status,omitempty"`
	// What caused the status of the object
	Cause string `json:"cause,omitempty"`
	// The specific details for what happened to the object
	Detail string `json:"detail,omitempty"`
	// Any errors associated with creation of the object
	Error string `json:"error,omitempty"`
}

type Relationship struct {
	// The type of Relationship
	// The only type is depends
	// Required
	Type OrchestrationRelationshipType `json:"type"`
	// What objects the relationship depends on
	// Required
	Targets []string `json:"targets"`
}

// CreateOrchestration creates a new Orchestration with the given name, key and enabled flag.
func (c *OrchestrationsClient) CreateOrchestration(input *CreateOrchestrationInput) (*Orchestration, error) {
	var createdOrchestration Orchestration

	input.Name = c.getQualifiedName(input.Name)
	for _, i := range input.Objects {
		i.Orchestration = c.getQualifiedName(i.Orchestration)
		if i.Type == OrchestrationTypeInstance {
			instanceClient := c.ComputeClient.Instances()
			instanceInput := i.Template.(*CreateInstanceInput)
			instanceInput.Name = c.getQualifiedName(instanceInput.Name)

			qualifiedSSHKeys := []string{}
			for _, key := range instanceInput.SSHKeys {
				qualifiedSSHKeys = append(qualifiedSSHKeys, c.getQualifiedName(key))
			}

			instanceInput.SSHKeys = qualifiedSSHKeys

			qualifiedStorageAttachments := []StorageAttachmentInput{}
			for _, attachment := range instanceInput.Storage {
				qualifiedStorageAttachments = append(qualifiedStorageAttachments, StorageAttachmentInput{
					Index:  attachment.Index,
					Volume: c.getQualifiedName(attachment.Volume),
				})
			}
			instanceInput.Storage = qualifiedStorageAttachments

			instanceInput.Networking = instanceClient.qualifyNetworking(instanceInput.Networking)

		}
	}

	if err := c.createResource(&input, &createdOrchestration); err != nil {
		return nil, err
	}

	// Call wait for orchestration ready now, as creating the orchestration is an eventually consistent operation
	getInput := &GetOrchestrationInput{
		Name: createdOrchestration.Name,
	}

	if input.Timeout == 0 {
		input.Timeout = WaitForOrchestrationActiveTimeout
	}

	// Wait for orchestration to be ready and return the result
	// Don't have to unqualify any objects, as the GetOrchestration method will handle that
	orchestrationInfo, orchestrationError := c.WaitForOrchestrationState(getInput, input.Timeout)
	if orchestrationError != nil {
		deleteInput := &DeleteOrchestrationInput{
			Name: createdOrchestration.Name,
		}
		err := c.DeleteOrchestration(deleteInput)
		if err != nil {
			return nil, fmt.Errorf("Error deleting orchestration %s: %s", getInput.Name, err)
		}
		return nil, fmt.Errorf("Error creating orchestration %s: %s", getInput.Name, orchestrationError)
	}

	return &orchestrationInfo, nil
}

// GetOrchestrationInput describes the Orchestration to get
type GetOrchestrationInput struct {
	// The three-part name of the Orchestration (/Compute-identity_domain/user/object).
	Name string `json:name`
}

// GetOrchestration retrieves the Orchestration with the given name.
func (c *OrchestrationsClient) GetOrchestration(input *GetOrchestrationInput) (*Orchestration, error) {
	var orchestrationInfo Orchestration
	if err := c.getResource(input.Name, &orchestrationInfo); err != nil {
		return nil, err
	}

	return c.success(&orchestrationInfo)
}

// UpdateOrchestrationInput defines an Orchestration to be updated
type UpdateOrchestrationInput struct {
	// The default Oracle Compute Cloud Service account, such as /Compute-acme/default.
	// Optional
	Account string `json:"account,omitempty"`
	// Description of this orchestration
	// Optional
	Description string `json:"description,omitempty"`
	// Specify the desired state of this orchestration: active, inactive, or suspend.
	// You can manage the state of the orchestration objects by changing the desired state of the orchestration.
	// * active: Creates all the orchestration objects defined in the orchestration.
	// * inactive: Adds the orchestration to Oracle Compute Cloud Service, but does not create any of the orchestration
	// objects defined in the orchestration.
	// Required
	DesiredState OrchestrationDesiredState `json:"desired_state"`
	// The three-part name of the Orchestration (/Compute-identity_domain/user/object).
	// Object names can contain only alphanumeric characters, hyphens, underscores, and periods. Object names are case-sensitive.
	// Required
	Name string `json:"name"`
	// The list of objects in the orchestration. An object is the primary building block of an orchestration.
	// An orchestration can contain up to 100 objects.
	// Required
	Objects []Object `json:"objects"`
	// Strings that describe the orchestration and help you identify it.
	Tags []string `json:"tags,omitempty"`
	// Version of this orchestration. It is automatically generated by the server.
	Version int `json:"version,omitempty"`
	// Time to wait for an orchestration to be ready
	Timeout time.Duration `json:"-"`
}

// UpdateOrchestration updates the orchestration.
func (c *OrchestrationsClient) UpdateOrchestration(input *UpdateOrchestrationInput) (*Orchestration, error) {
	var updatedOrchestration Orchestration
	input.Name = c.getQualifiedName(input.Name)
	for _, i := range input.Objects {
		i.Orchestration = c.getQualifiedName(i.Orchestration)
		if i.Type == OrchestrationTypeInstance {
			instanceInput := i.Template.(map[string]interface{})
			instanceInput["name"] = c.getQualifiedName(instanceInput["name"].(string))
		}
	}

	if err := c.updateResource(input.Name, input, &updatedOrchestration); err != nil {
		return nil, err
	}

	// Call wait for orchestration ready now, as creating the orchestration is an eventually consistent operation
	getInput := &GetOrchestrationInput{
		Name: updatedOrchestration.Name,
	}

	if input.Timeout == 0 {
		input.Timeout = WaitForOrchestrationActiveTimeout
	}

	// Wait for orchestration to be ready and return the result
	// Don't have to unqualify any objects, as the GetOrchestration method will handle that
	orchestrationInfo, orchestrationError := c.WaitForOrchestrationState(getInput, input.Timeout)
	if orchestrationError != nil {
		return nil, orchestrationError
	}

	return &orchestrationInfo, nil
}

// DeleteOrchestrationInput describes the Orchestration to delete
type DeleteOrchestrationInput struct {
	// The three-part name of the Orchestration (/Compute-identity_domain/user/object).
	// Required
	Name string `json:name`
	// Timeout for delete request
	Timeout time.Duration `json:"-"`
}

// DeleteOrchestration deletes the Orchestration with the given name.
func (c *OrchestrationsClient) DeleteOrchestration(input *DeleteOrchestrationInput) error {
	if err := c.deleteOrchestration(input.Name); err != nil {
		return err
	}

	if input.Timeout == 0 {
		input.Timeout = WaitForOrchestrationDeleteTimeout
	}

	return c.WaitForOrchestrationDeleted(input, input.Timeout)
}

func (c *OrchestrationsClient) success(info *Orchestration) (*Orchestration, error) {
	c.unqualify(&info.Name)
	for _, i := range info.Objects {
		c.unqualify(&i.Orchestration)
		if OrchestrationType(i.Type) == OrchestrationTypeInstance {
			instanceInput := i.Template.(map[string]interface{})
			instanceInput["name"] = c.getUnqualifiedName(instanceInput["name"].(string))
		}
	}

	return info, nil
}

// WaitForOrchestrationActive waits for an orchestration to be completely initialized and available.
func (c *OrchestrationsClient) WaitForOrchestrationState(input *GetOrchestrationInput, timeout time.Duration) (Orchestration, error) {
	var info *Orchestration
	var getErr error
	err := c.client.WaitFor("orchestration to be ready", timeout, func() (bool, error) {
		info, getErr = c.GetOrchestration(input)
		if getErr != nil {
			return false, getErr
		}
		c.client.DebugLogString(fmt.Sprintf("Orchestration name is %v, Orchestration info is %+v", info.Name, info))
		switch s := info.Status; s {
		case OrchestrationStatusError:
			// We need to check and see if an object the orchestration is trying to create is giving us an error instead of just the orchestration as a whole.
			for _, object := range info.Objects {
				if object.Health.Status == OrchestrationStatusError {
					return false, fmt.Errorf("Error creating instance %s: %+v", object.Name, object.Health)
				}
			}
			return false, fmt.Errorf("Error initializing orchestration: %+v", info)
		case OrchestrationStatus(info.DesiredState):
			c.client.DebugLogString(fmt.Sprintf("Orchestration %s", info.DesiredState))
			return true, nil
		case OrchestrationStatusActivating:
			c.client.DebugLogString("Orchestration activating")
			return false, nil
		case OrchestrationStatusStopping:
			c.client.DebugLogString("Orchestration stopping")
			return false, nil
		case OrchestrationStatusSuspending:
			c.client.DebugLogString("Orchestration suspending")
			return false, nil
		case OrchestrationStatusDeactivating:
			c.client.DebugLogString("Orchestration deactivating")
			return false, nil
		case OrchestrationStatusSuspended:
			c.client.DebugLogString("Orchestration suspended")
			if info.DesiredState == OrchestrationDesiredStateSuspend {
				return true, nil
			} else {
				return false, nil
			}
		default:
			return false, fmt.Errorf("Unknown orchestration state: %s, erroring", s)
		}
	})
	return *info, err
}

// WaitForOrchestrationDeleted waits for an orchestration to be fully deleted.
func (c *OrchestrationsClient) WaitForOrchestrationDeleted(input *DeleteOrchestrationInput, timeout time.Duration) error {
	return c.client.WaitFor("orchestration to be deleted", timeout, func() (bool, error) {
		var info Orchestration
		if err := c.getResource(input.Name, &info); err != nil {
			if client.WasNotFoundError(err) {
				// Orchestration could not be found, thus deleted
				return true, nil
			}
			// Some other error occurred trying to get Orchestration, exit
			return false, err
		}
		switch s := info.Status; s {
		case OrchestrationStatusError:
			return false, fmt.Errorf("Error stopping orchestration: %+v", info)
		case OrchestrationStatusStopping:
			c.client.DebugLogString("Orchestration stopping")
			return false, nil
		case OrchestrationStatusDeleting:
			c.client.DebugLogString("Orchestration deleting")
			return false, nil
		case OrchestrationStatusActive:
			c.client.DebugLogString("Orchestration active")
			return false, nil
		default:
			return false, fmt.Errorf("Unknown orchestration state: %s, erroring", s)
		}
	})
}
