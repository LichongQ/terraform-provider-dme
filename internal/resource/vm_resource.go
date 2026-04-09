// Package resource implements the VM resource for the eDME Terraform provider.
package resource

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"terraform-provider-dme/internal/api"
	"terraform-provider-dme/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceVM returns the schema and CRUD operations for the dme_vm resource.
func ResourceVM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVMCreate,
		ReadContext:   resourceVMRead,
		UpdateContext: resourceVMUpdate,
		DeleteContext: resourceVMDelete,
		Description:   "Manages a virtual machine on eDME. The VM is created from a template but not started by default.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The template ID to clone from (e.g. `urn:sites:XXXXX:vms:i-XXXXXX`)",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Virtual machine name, 1~256 characters",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Virtual machine description, 0~1024 characters",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID or host ID where the VM will be deployed (e.g. `urn:sites:XXXXX:clusters:173`)",
			},
			"is_binding_host": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to bind the VM to a specific host",
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of CPU cores, range [1,255]",
			},
			"memory_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Memory size in MB, range [128,4194304]",
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operating system type: Windows, Linux, or Other",
			},
			"os_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operating system version number",
			},
			"guest_os_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Guest OS name",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname for the VM",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Administrator password for the VM",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Static IPv4 address for the VM",
			},
			"netmask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 subnet mask",
			},
			"gateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 gateway",
			},
			"port_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network port group ID for the VM NIC",
			},
			"disks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Disk configuration for the VM",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sequence_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Disk slot number",
						},
						"datastore_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Datastore ID for the disk",
						},
						"capacity_gb": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Disk capacity in GB",
						},
						"pci_type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Bus type: IDE, SCSI, or VIRTIO",
							ValidateFunc: validation.StringInSlice([]string{"IDE", "SCSI", "VIRTIO"}, false),
						},
						"indip_disk": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether this is an independent disk",
						},
						"persistent_disk": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether this is a persistent disk",
						},
						"is_thin": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether to use thin provisioning",
						},
					},
				},
			},
			// Computed attributes
			"vm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The actual VM ID assigned by eDME after creation (e.g. `urn:sites:XXXXX:vms:i-XXXXXX`)",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Desired VM status: running, stopped, or restarted. Use 'restarted' to trigger a reboot. Defaults to stopped on create.",
			},
			"ip_address_actual": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Actual IP address of the VM as reported by eDME",
			},
			"vr_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtualization platform type: FUSIONCOMPUTE, VMWARE, or HCS",
			},
			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Host ID where the VM resides",
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Host name where the VM resides",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID where the VM resides",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster name where the VM resides",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Datacenter ID where the VM resides",
			},
			"datacenter_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Datacenter name where the VM resides",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VM creation time",
			},
			"arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CPU architecture: x86 or arm",
			},
		},
	}
}

// resourceVMCreate creates a new VM by cloning from a template.
// The VM is created with is_auto_boot=false (not started by default).
func resourceVMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, err := getClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	templateID := d.Get("template_id").(string)

	cloneReq := api.CloneRequest{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Location:      d.Get("location").(string),
		IsBindingHost: d.Get("is_binding_host").(bool),
		IsAutoBoot:    d.Get("status").(string) == "running", // Auto-boot only if desired status is running
		VmNum:         1,
	}

	// Set OS options if provided.
	if osType := d.Get("os_type").(string); osType != "" {
		cloneReq.OsOptions = &api.OsOptions{
			OsType:      osType,
			OsVersion:   d.Get("os_version").(string),
			GuestOSName: d.Get("guest_os_name").(string),
		}
	}

	// Set guest customization if hostname or password is provided.
	if hostname := d.Get("hostname").(string); hostname != "" {
		guestCustom := &api.GuestOSCustomization{
			Hostname:           hostname,
			Password:           d.Get("password").(string),
			IsUpdateVmPassword: true,
		}
		if osType := d.Get("os_type").(string); osType != "" {
			guestCustom.OsType = osType
		}

		// Set NIC specification if IP/port_group is provided.
		if portGroupID := d.Get("port_group_id").(string); portGroupID != "" {
			guestCustom.NicSpecification = []api.NicInfo{
				{
					SequenceNum: "0",
					IpVersion:   4,
					Ip:          d.Get("ip_address").(string),
					Netmask:     d.Get("netmask").(string),
					Gateway:     d.Get("gateway").(string),
				},
			}
		}
		cloneReq.GuestCustomization = guestCustom
	}

	// Set VM config (CPU, memory, disks, NICs) if provided.
	vmConfig := &api.VmSimpleConfig{}
	hasConfig := false

	if cpuCores := d.Get("cpu_cores").(int); cpuCores > 0 {
		vmConfig.CPU = &api.VMCpu{Quantity: cpuCores}
		hasConfig = true
	}
	if memMB := d.Get("memory_mb").(int); memMB > 0 {
		vmConfig.Memory = &api.VMMemory{QuantitySize: memMB}
		hasConfig = true
	}

	// Set disks.
	if diskList := d.Get("disks").([]interface{}); len(diskList) > 0 {
		for _, dItem := range diskList {
			diskMap := dItem.(map[string]interface{})
			vmConfig.Disks = append(vmConfig.Disks, api.VMDisk{
				SequenceNum:    diskMap["sequence_num"].(int),
				DatastoreID:    diskMap["datastore_id"].(string),
				Quantity:       diskMap["capacity_gb"].(int),
				PciType:        diskMap["pci_type"].(string),
				IndepDisk:      diskMap["indip_disk"].(bool),
				PersistentDisk: diskMap["persistent_disk"].(bool),
				IsThin:         diskMap["is_thin"].(bool),
			})
		}
		hasConfig = true
	}

	// Set NICs.
	if portGroupID := d.Get("port_group_id").(string); portGroupID != "" {
		vmConfig.Nics = []api.VMNic{{PortGroupID: portGroupID}}
		hasConfig = true
	}

	if hasConfig {
		cloneReq.VmConfig = vmConfig
	}

	// Clone the VM from template.
	path := fmt.Sprintf("/vmmgmt/v1/vms/%s/clone", templateID)
	var cloneResp api.AsyncTaskResponse
	if err := cl.Do(ctx, "POST", path, cloneReq, &cloneResp); err != nil {
		return diag.Errorf("failed to clone VM from template %s: %s", templateID, err)
	}

	log.Printf("[INFO] Clone task %s initiated for VM %s, waiting for completion...", cloneResp.TaskID, cloneReq.Name)

	// Wait for the VM to be created and get its actual ID.
	vmID, err := waitForVMCreated(ctx, cl, templateID, cloneResp.TaskID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("failed to confirm VM creation: %s", err)
	}

	// Use a deterministic local ID for Terraform state.
	localID := fmt.Sprintf("dme-vm-%x", sha256.Sum256([]byte(templateID+cloneResp.TaskID)))[:16]
	d.SetId(localID)
	// Store the actual eDME VM ID as a computed attribute.
	_ = d.Set("vm_id", vmID)

	log.Printf("[INFO] VM %s created with eDME ID %s", cloneReq.Name, vmID)

	return resourceVMRead(ctx, d, meta)
}

// resourceVMRead queries the real VM state from eDME and syncs to Terraform state.
// This forces a fresh read from the API, not from cached state.
func resourceVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, err := getClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	vmID := d.Get("vm_id").(string)
	if vmID == "" {
		// VM hasn't been created yet or was never set; nothing to read.
		return nil
	}

	var vm api.Vm
	path := fmt.Sprintf("/vmmgmt/v1/vms/%s", vmID)
	if err := cl.Do(ctx, "GET", path, nil, &vm); err != nil {
		// If the VM no longer exists, remove it from state.
		log.Printf("[WARN] Failed to read VM %s: %s", vmID, err)
		d.SetId("")
		_ = d.Set("vm_id", "")
		return nil
	}

	// Check if VM has been deleted on the server side.
	if vm.Status == "deleted" || vm.Status == "pending_delete" {
		log.Printf("[INFO] VM %s is in status %s, removing from state", vmID, vm.Status)
		d.SetId("")
		_ = d.Set("vm_id", "")
		return nil
	}

	// Sync all attributes from the real VM state.
	_ = d.Set("vm_id", vm.ID)
	_ = d.Set("name", vm.Name)
	_ = d.Set("description", vm.Description)
	_ = d.Set("status", vm.Status)
	_ = d.Set("ip_address_actual", vm.IPAddress)
	_ = d.Set("vr_type", vm.VrType)
	_ = d.Set("host_id", vm.HostID)
	_ = d.Set("host_name", vm.HostName)
	_ = d.Set("cluster_id", vm.ClusterID)
	_ = d.Set("cluster_name", vm.ClusterName)
	_ = d.Set("datacenter_id", vm.DatacenterID)
	_ = d.Set("datacenter_name", vm.DatacenterName)
	_ = d.Set("create_time", vm.CreateTime)
	_ = d.Set("arch", vm.Arch)

	// Extract CPU and memory from the real VM state.
	if vm.CPU != nil {
		_ = d.Set("cpu_cores", vm.CPU.Quantity)
	}
	if vm.Memory != nil {
		_ = d.Set("memory_mb", vm.Memory.QuantitySize)
	}

	// Extract OS options.
	if vm.OsOptions != nil {
		_ = d.Set("os_type", vm.OsOptions.OsType)
		_ = d.Set("os_version", vm.OsOptions.OsVersion)
		_ = d.Set("guest_os_name", vm.OsOptions.GuestOSName)
	}

	return nil
}

// resourceVMUpdate modifies the VM properties and optionally changes its power state.
// It implements idempotency: if the real state matches the desired state, no API call is made.
func resourceVMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, err := getClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	vmID := d.Get("vm_id").(string)
	if vmID == "" {
		return diag.Errorf("VM ID not set, cannot update")
	}

	// First, read the current real state.
	var vm api.Vm
	path := fmt.Sprintf("/vmmgmt/v1/vms/%s", vmID)
	if err := cl.Do(ctx, "GET", path, nil, &vm); err != nil {
		return diag.Errorf("failed to read VM %s before update: %s", vmID, err)
	}

	// Idempotency check: compare desired state with real state.
	// If nothing has changed and the real state already matches, skip API calls.
	dirty := false

	// Build the update request only with changed fields.
	updateReq := api.UpdateVmRequest{}
	if d.HasChange("name") && vm.Name != d.Get("name").(string) {
		updateReq.Name = d.Get("name").(string)
		dirty = true
	}
	if d.HasChange("description") && vm.Description != d.Get("description").(string) {
		updateReq.Description = d.Get("description").(string)
		dirty = true
	}
	if d.HasChange("cpu_cores") {
		newCores := d.Get("cpu_cores").(int)
		if vm.CPU == nil || vm.CPU.Quantity != newCores {
			updateReq.CPU = &api.VMCpu{Quantity: newCores}
			dirty = true
		}
	}
	if d.HasChange("memory_mb") {
		newMem := d.Get("memory_mb").(int)
		if vm.Memory == nil || vm.Memory.QuantitySize != newMem {
			updateReq.Memory = &api.VMMemory{QuantitySize: newMem}
			dirty = true
		}
	}

	// DME北向接口不支持查询guest_os_name，无法对此字段对比，默认不修改
	// if d.HasChange("os_type") || d.HasChange("os_version") || d.HasChange("guest_os_name") {
	if d.HasChange("os_type") || d.HasChange("os_version") {
		newOsType := d.Get("os_type").(string)
		newOsVer := d.Get("os_version").(string)
		newGuestName := d.Get("guest_os_name").(string)
		osChanged := (vm.OsOptions == nil) ||
			vm.OsOptions.OsType != newOsType ||
			vm.OsOptions.OsVersion != newOsVer ||
			vm.OsOptions.GuestOSName != newGuestName
		if osChanged {
			updateReq.OsOptions = &api.OsOptions{
				OsType:      newOsType,
				OsVersion:   newOsVer,
				GuestOSName: newGuestName,
			}
			dirty = true
		}
	}

	// Only call the update API if there are actual changes.
	if dirty {
		var updateResp api.AsyncTaskResponse
		if err := cl.Do(ctx, "PUT", path, updateReq, &updateResp); err != nil {
			return diag.Errorf("failed to update VM %s: %s", vmID, err)
		}
		if err := waitForTask(ctx, cl, updateResp.TaskID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("failed to confirm VM update task completion: %s", err)
		}
		log.Printf("[INFO] VM %s properties updated, task %s completed", vmID, updateResp.TaskID)
	}

	// Handle power state changes.
	desiredStatus := d.Get("status").(string)
	currentStatus := vm.Status

	if desiredStatus != "" && desiredStatus != currentStatus {
		// Power on: desired is "running" but current is "stopped"
		if desiredStatus == "running" && (currentStatus == "stopped" || currentStatus == "shut_off") {
			powerReq := api.PowerActionRequest{VmIDs: []string{vmID}}
			var powerResp api.AsyncTaskResponse
			if err := cl.Do(ctx, "POST", "/vmmgmt/v1/vms/batch-power-on", powerReq, &powerResp); err != nil {
				return diag.Errorf("failed to power on VM %s: %s", vmID, err)
			}
			if err := waitForTask(ctx, cl, powerResp.TaskID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return diag.Errorf("failed to confirm VM power-on task completion: %s", err)
			}
			log.Printf("[INFO] Power-on task %s completed for VM %s", powerResp.TaskID, vmID)
		}
		// Power off: desired is "stopped" but current is "running"
		if desiredStatus == "stopped" && currentStatus == "running" {
			powerReq := api.PowerActionRequest{VmIDs: []string{vmID}, Mode: "force"}
			var powerResp api.AsyncTaskResponse
			if err := cl.Do(ctx, "POST", "/vmmgmt/v1/vms/batch-power-off", powerReq, &powerResp); err != nil {
				return diag.Errorf("failed to power off VM %s: %s", vmID, err)
			}
			if err := waitForTask(ctx, cl, powerResp.TaskID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return diag.Errorf("failed to confirm VM power-off task completion: %s", err)
			}
			log.Printf("[INFO] Power-off task %s completed for VM %s", powerResp.TaskID, vmID)
		}
		// Reboot: desired is "restarted"
		if desiredStatus == "restarted" && currentStatus == "running" {
			powerReq := api.PowerActionRequest{VmIDs: []string{vmID}, Mode: "force"}
			var powerResp api.AsyncTaskResponse
			if err := cl.Do(ctx, "POST", "/vmmgmt/v1/vms/batch-reboot", powerReq, &powerResp); err != nil {
				return diag.Errorf("failed to reboot VM %s: %s", vmID, err)
			}
			if err := waitForTask(ctx, cl, powerResp.TaskID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return diag.Errorf("failed to confirm VM reboot task completion: %s", err)
			}
			log.Printf("[INFO] Reboot task %s completed for VM %s", powerResp.TaskID, vmID)
		}
	}

	return resourceVMRead(ctx, d, meta)
}

// resourceVMDelete deletes the VM using the batch-delete API.
func resourceVMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, err := getClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	vmID := d.Get("vm_id").(string)
	if vmID == "" {
		d.SetId("")
		return nil
	}

	deleteReq := api.BatchDeleteRequest{
		VmIDs: []string{vmID},
	}

	var deleteResp api.AsyncTaskResponse
	if err := cl.Do(ctx, "POST", "/vmmgmt/v1/vms/batch-delete", deleteReq, &deleteResp); err != nil {
		return diag.Errorf("failed to delete VM %s: %s", vmID, err)
	}

	if err := waitForTask(ctx, cl, deleteResp.TaskID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("failed to confirm VM delete task completion: %s", err)
	}

	log.Printf("[INFO] VM %s deleted, task %s completed", vmID, deleteResp.TaskID)
	d.SetId("")
	return nil
}

// waitForVMCreated polls the task detail API until the clone task completes,
// then extracts the created VM ID from the task's affected resources.
func waitForVMCreated(ctx context.Context, cl *client.Client, templateID, taskID string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	pollInterval := 10 * time.Second

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		time.Sleep(pollInterval)

		var tasks []api.TaskDetail
		path := fmt.Sprintf("/taskmgmt/v1/tasks/%s", taskID)
		if err := cl.Do(ctx, "GET", path, nil, &tasks); err != nil {
			log.Printf("[WARN] Failed to query task %s: %s", taskID, err)
			continue
		}

		if len(tasks) == 0 {
			log.Printf("[DEBUG] Task %s returned empty result, still running...", taskID)
			continue
		}

		task := tasks[0]
		switch task.Status {
		case 3, 4:
			// Task succeeded (3) or partially succeeded (4).
			// Extract the VM ID from affected resources.
			for _, r := range task.Resources {
				if r.Operate == "CREATE" && r.ID != "" {
					log.Printf("[INFO] Found created VM ID %s from task %s resources", r.ID, taskID)
					return r.ID, nil
				}
			}
			// Fallback: if no CREATE resource found, try any resource with an ID.
			for _, r := range task.Resources {
				if r.ID != "" {
					log.Printf("[INFO] Found resource ID %s from task %s (fallback)", r.ID, taskID)
					return r.ID, nil
				}
			}
			return "", fmt.Errorf("task %s succeeded but no affected resource found in response", taskID)
		case 5:
			return "", fmt.Errorf("task %s failed: %s", taskID, task.Remarks)
		case 6:
			return "", fmt.Errorf("task %s timed out", taskID)
		default:
			// Status 1 (initial) or 2 (running): continue polling.
			log.Printf("[DEBUG] Task %s status=%d progress=%d%%, still running...", taskID, task.Status, task.Progress)
		}
	}

	return "", fmt.Errorf("timeout waiting for task %s to complete after %s", taskID, timeout)
}

// waitForTask polls the task detail API until the task completes.
// The query task API returns an array of TaskDetail.
func waitForTask(ctx context.Context, cl *client.Client, taskID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	pollInterval := 10 * time.Second

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		time.Sleep(pollInterval)

		var tasks []api.TaskDetail
		path := fmt.Sprintf("/taskmgmt/v1/tasks/%s", taskID)
		if err := cl.Do(ctx, "GET", path, nil, &tasks); err != nil {
			log.Printf("[WARN] Failed to query task %s: %s", taskID, err)
			continue
		}

		if len(tasks) == 0 {
			log.Printf("[DEBUG] Task %s returned empty result, still running...", taskID)
			continue
		}

		task := tasks[0]
		switch task.Status {
		case 3, 4:
			log.Printf("[INFO] Task %s completed (status=%d)", taskID, task.Status)
			return nil
		case 5:
			return fmt.Errorf("task %s failed: %s", taskID, task.Remarks)
		case 6:
			return fmt.Errorf("task %s timed out", taskID)
		default:
			// Status 1 (initial) or 2 (running): continue polling.
			log.Printf("[DEBUG] Task %s status=%d progress=%d%%, still running...", taskID, task.Status, task.Progress)
		}
	}

	return fmt.Errorf("timeout waiting for task %s to complete after %s", taskID, timeout)
}

// getClient extracts the eDME client from the provider meta.
func getClient(meta interface{}) (*client.Client, error) {
	cl, ok := meta.(*client.Client)
	if !ok {
		return nil, fmt.Errorf("failed to cast provider meta to eDME client")
	}
	return cl, nil
}
