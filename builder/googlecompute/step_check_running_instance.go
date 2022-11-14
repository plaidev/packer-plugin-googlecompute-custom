package googlecompute

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
)

// StepCheckRunningInstance represents a Packer build step that creates GCE instances.
type StepCheckRunningInstance struct {
	Debug         bool
	GeneratedData *packerbuilderdata.GeneratedData
}

// Check instance is running.
func (s *StepCheckRunningInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)
	d := state.Get("driver").(Driver)

	ui := state.Get("ui").(packersdk.Ui)

	name := c.InstanceName

	status, err := d.GetInstanceState(c.Zone, c.RunningInstanceName)
	if err != nil {
		err := fmt.Errorf("Error getting instance state for check running: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if status != "RUNNING" {
		err := fmt.Errorf("Instance %s is not running.\n"+
			"Instance set in RunningInstanceName must be running.", c.RunningInstanceName)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Instance is running")

	// Things succeeded, store the name so we can remove it later
	state.Put("instance_name", name)
	// instance_id is the generic term used so that users can have access to the
	// instance id inside of the provisioners, used in step_provision.
	state.Put("instance_id", name)

	return multistep.ActionContinue
}
