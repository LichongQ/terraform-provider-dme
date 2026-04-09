// Package datasource implements the VMs data source for the eDME Terraform provider.
package datasource

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-dme/internal/api"
	"terraform-provider-dme/internal/client"
)

// DataSourceVMs returns the schema and Read operation for the dme_vms data source.
func DataSourceVMs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVMsRead,
		Description: "Queries a list of virtual machines from eDME with optional filters.",
		Schema: map[string]*schema.Schema{
			// Filters
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Site ID to filter VMs",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster ID to filter VMs (not supported for HCS)",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Datacenter ID to filter VMs (VMWARE only)",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster name to filter VMs, supports fuzzy match (not HCS)",
			},
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host ID to filter VMs",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host name to filter VMs, supports fuzzy match",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VM name to filter, supports fuzzy match",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VM IP address to filter, supports fuzzy match",
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of VM statuses to filter (max 20 items)",
			},
			"is_template": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter by template flag: true=templates only, false=VMs only",
			},
			"os_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of OS types to filter: Windows, Linux, Other (max 3 items)",
			},
			"vr_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Virtualization platform: FUSIONCOMPUTE, VMWARE, HCS",
			},
			// Sorting
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field: name, cpu_core, memory_size, disk_total_size, create_time, ip_address",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "asc",
				Description: "Sort direction: asc or desc (default asc)",
			},
			// Pagination
			"page_no": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number (default 1)",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Items per page (default 100, max 1000)",
			},
			// Computed outputs
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matching VMs",
			},
			"vms": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of VMs matching the filters",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VM ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VM name",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VM description",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VM status",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VM IP address",
						},
						"site_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name",
						},
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host ID",
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host name",
						},
						"datacenter_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Datacenter ID",
						},
						"datacenter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Datacenter name",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"arch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CPU architecture: x86 or arm",
						},
						"vr_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtualization platform type",
						},
						"is_template": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is a template",
						},
						"disk_total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total disk space",
						},
						"disk_used_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used disk space",
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
				},
			},
		},
	}
}

func dataSourceVMsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cl, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("failed to cast provider meta to eDME client")
	}

	req := api.VmListRequest{
		SiteID:       d.Get("site_id").(string),
		ClusterID:    d.Get("cluster_id").(string),
		DatacenterID: d.Get("datacenter_id").(string),
		ClusterName:  d.Get("cluster_name").(string),
		HostID:       d.Get("host_id").(string),
		HostName:     d.Get("host_name").(string),
		Name:         d.Get("name").(string),
		IPAddress:    d.Get("ip_address").(string),
		VrType:       d.Get("vr_type").(string),
		SortKey:      d.Get("sort_key").(string),
		SortDir:      d.Get("sort_dir").(string),
		PageNo:       int32(d.Get("page_no").(int)),
		PageSize:     int32(d.Get("page_size").(int)),
	}

	if v, ok := d.GetOk("is_template"); ok {
		b := v.(bool)
		req.IsTemplate = &b
	}

	if v := d.Get("status").([]interface{}); len(v) > 0 {
		for _, s := range v {
			req.Status = append(req.Status, s.(string))
		}
	}
	if v := d.Get("os_type").([]interface{}); len(v) > 0 {
		for _, s := range v {
			req.OsType = append(req.OsType, s.(string))
		}
	}

	var resp api.VmsResponse
	if err := cl.Do(ctx, "POST", "/vmmgmt/v1/vms/query", req, &resp); err != nil {
		return diag.Errorf("failed to query VM list: %s", err)
	}

	log.Printf("[INFO] Queried VM list: total=%d, returned=%d", resp.Total, len(resp.VMs))

	// Flatten VMs into the data source.
	vms := make([]map[string]interface{}, 0, len(resp.VMs))
	for _, vm := range resp.VMs {
		vmMap := map[string]interface{}{
			"id":              vm.ID,
			"name":            vm.Name,
			"description":     vm.Description,
			"status":          vm.Status,
			"ip_address":      vm.IPAddress,
			"site_id":         vm.SiteID,
			"cluster_id":      vm.ClusterID,
			"cluster_name":    vm.ClusterName,
			"host_id":         vm.HostID,
			"host_name":       vm.HostName,
			"datacenter_id":   vm.DatacenterID,
			"datacenter_name": vm.DatacenterName,
			"create_time":     vm.CreateTime,
			"arch":            vm.Arch,
			"vr_type":         vm.VrType,
			"is_template":     vm.IsTemplate,
			"disk_total_size": int(vm.DiskTotalSize),
			"disk_used_size":  int(vm.DiskUsedSize),
			"disk_num":        vm.DiskNum,
		}
		if vm.CPU != nil {
			vmMap["cpu_cores"] = vm.CPU.Quantity
		}
		if vm.Memory != nil {
			vmMap["memory_mb"] = vm.Memory.QuantitySize
		}
		if vm.OsOptions != nil {
			vmMap["os_type"] = vm.OsOptions.OsType
			vmMap["os_version"] = vm.OsOptions.OsVersion
			vmMap["guest_os_name"] = vm.OsOptions.GuestOSName
		}
		vms = append(vms, vmMap)
	}

	// Sort vms for deterministic output.
	sort.Slice(vms, func(i, j int) bool {
		return vms[i]["id"].(string) < vms[j]["id"].(string)
	})

	_ = d.Set("total", resp.Total)
	_ = d.Set("vms", vms)

	// Generate a deterministic ID from the query parameters.
	queryStr := fmt.Sprintf("%v", req)
	d.SetId(fmt.Sprintf("dme-vms-%x", sha256.Sum256([]byte(queryStr)))[:16])

	return nil
}
