// Package datasource implements the VM data source for the eDME Terraform provider.
package datasource

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-dme/internal/api"
	"terraform-provider-dme/internal/client"
)

// DataSourceVM returns the schema and Read operation for the dme_vm data source.
func DataSourceVM() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVMRead,
		Description: "Fetches information about an existing virtual machine from eDME.",
		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VM ID to query (e.g. `urn:sites:XXXXX:vms:i-XXXXXX`)",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtual machine name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtual machine description",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current VM status: running, stopped, creating, etc.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the VM",
			},
			"site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site ID where the VM resides",
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
			"vr_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtualization platform type: FUSIONCOMPUTE, VMWARE, or HCS",
			},
			"is_template": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a template",
			},
			"disk_total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total disk space allocated",
			},
			"disk_used_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Disk space used",
			},
			"disk_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of disks",
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of CPU cores",
			},
			"memory_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Memory size in MB",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating system type",
			},
			"os_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating system version",
			},
			"guest_os_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Guest OS name",
			},
		},
	}
}

// dataSourceVMRead fetches the VM information from eDME.
func dataSourceVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("failed to cast provider meta to eDME client")
	}

	vmID := d.Get("vm_id").(string)

	var vm api.Vm
	path := fmt.Sprintf("/vmmgmt/v1/vms/%s", vmID)
	if err := cl.Do(ctx, "GET", path, nil, &vm); err != nil {
		return diag.Errorf("failed to query VM %s: %s", vmID, err)
	}

	// Check if VM exists.
	if vm.Status == "deleted" || vm.Status == "pending_delete" {
		return diag.Errorf("VM %s has been deleted", vmID)
	}

	d.SetId(vmID)

	// Populate all attributes from the VM response.
	_ = d.Set("vm_id", vm.ID)
	_ = d.Set("name", vm.Name)
	_ = d.Set("description", vm.Description)
	_ = d.Set("status", vm.Status)
	_ = d.Set("ip_address", vm.IPAddress)
	_ = d.Set("site_id", vm.SiteID)
	_ = d.Set("cluster_id", vm.ClusterID)
	_ = d.Set("cluster_name", vm.ClusterName)
	_ = d.Set("host_id", vm.HostID)
	_ = d.Set("host_name", vm.HostName)
	_ = d.Set("datacenter_id", vm.DatacenterID)
	_ = d.Set("datacenter_name", vm.DatacenterName)
	_ = d.Set("create_time", vm.CreateTime)
	_ = d.Set("arch", vm.Arch)
	_ = d.Set("vr_type", vm.VrType)
	_ = d.Set("is_template", vm.IsTemplate)
	_ = d.Set("disk_total_size", int(vm.DiskTotalSize))
	_ = d.Set("disk_used_size", int(vm.DiskUsedSize))
	_ = d.Set("disk_num", vm.DiskNum)

	// Extract CPU and memory.
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

	log.Printf("[INFO] Data source read VM %s: name=%s, status=%s", vmID, vm.Name, vm.Status)
	return nil
}
