package googlecompute

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
	"io/ioutil"
	"log"
	"strings"
)

// SetMetadataToRunningInstance represents a Packer build step that creates GCE instances.
type SetMetadataToRunningInstance struct {
	Debug         bool
	GeneratedData *packerbuilderdata.GeneratedData
}

func (c *Config) createInstanceMetadataSSHKeyMerged(sshPublicKey string) (map[string]string, error) {

	instanceMetadata := make(map[string]string)

	var err error
	var errs *packersdk.MultiError

	// Copy metadata from config.
	for k, v := range c.Metadata {
		instanceMetadata[k] = v
	}
	log.Printf("aaa")

	// Merge any existing ssh keys with our public key, unless there is no
	// supplied public key. This is possible if a private_key_file was
	// specified.
	if sshPublicKey != "" {
		sshMetaKey := "ssh-keys"
		sshPublicKey = strings.TrimSuffix(sshPublicKey, "\n")
		sshKeys := fmt.Sprintf("%s:%s %s", c.Comm.SSHUsername, sshPublicKey, c.Comm.SSHUsername)
		if confSSHKeys, exists := instanceMetadata[sshMetaKey]; exists {
			sshKeys = fmt.Sprintf("%s\n%s", sshKeys, confSSHKeys)
		}
		instanceMetadata[sshMetaKey] = sshKeys
	}
	log.Printf("bbbb")

	// If UseOSLogin is true, force `enable-oslogin` in metadata
	// In the event that `enable-oslogin` is not enabled at project level
	//if c.UseOSLogin {
	//	instanceMetadata[EnableOSLoginKey] = "TRUE"
	//}

	for key, value := range c.MetadataFiles {
		var content []byte
		content, err = ioutil.ReadFile(value)
		if err != nil {
			errs = packersdk.MultiErrorAppend(errs, err)
		}
		instanceMetadata[key] = string(content)
	}
	log.Printf("ccc")

	if errs != nil && len(errs.Errors) > 0 {
		return instanceMetadata, errs
	}
	return instanceMetadata, nil
}

// Check instance is running.
func (s *SetMetadataToRunningInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)

	ui := state.Get("ui").(packersdk.Ui)

	name := c.RunningInstanceName
	ui.Say("Set Metadata to running instance...")

	ui.Say("Set Metadata to running instance3...")

	// Things succeeded, store the name so we can remove it later
	state.Put("instance_name", name)
	// instance_id is the generic term used so that users can have access to the
	// instance id inside of the provisioners, used in step_provision.
	state.Put("instance_id", name)

	return multistep.ActionContinue
}

func (s *SetMetadataToRunningInstance) Cleanup(state multistep.StateBag) {
}
