// Package api defines all types used for eDME API request/response payloads.
package api

// AuthRequest is the request body for the session authentication API.
// PUT /rest/plat/smapp/v1/sessions
type AuthRequest struct {
	GrantType string `json:"grantType"` // Authentication type, always "password"
	UserName  string `json:"userName"`  // eDME northbound user name
	Value     string `json:"value"`     // eDME northbound user password
}

// AuthResponse is the response body from the session authentication API.
type AuthResponse struct {
	AccessSession  string            `json:"accessSession"`  // Token for subsequent API calls
	RoaRand        string            `json:"roaRand"`        // Random string for signing
	Expires        int               `json:"Expires"`        // Token TTL in seconds
	AdditionalInfo map[string]string `json:"additionalInfo"` // Extra info (e.g. passwdStatus)
}

// AsyncTaskResponse is returned by async operations (clone, update, power-on/off, etc.).
type AsyncTaskResponse struct {
	TaskID string `json:"task_id"` // Task ID for tracking async operations
}

// Vm represents a virtual machine as returned by the GET /rest/vmmgmt/v1/vms/{id} API.
type Vm struct {
	ID             string        `json:"id,omitempty"`
	ProjectID      string        `json:"project_id,omitempty"`
	PID            string        `json:"pid,omitempty"`
	Name           string        `json:"name,omitempty"`
	Description    string        `json:"description,omitempty"`
	Status         string        `json:"status,omitempty"`
	IPAddress      string        `json:"ip_address,omitempty"`
	SiteID         string        `json:"site_id,omitempty"`
	ClusterID      string        `json:"cluster_id,omitempty"`
	ClusterName    string        `json:"cluster_name,omitempty"`
	HostID         string        `json:"host_id,omitempty"`
	HostName       string        `json:"host_name,omitempty"`
	CreateTime     string        `json:"create_time,omitempty"`
	DiskTotalSize  int64         `json:"disk_total_size,omitempty"`
	DiskUsedSize   int64         `json:"disk_used_size,omitempty"`
	DiskNum        int           `json:"disk_num,omitempty"`
	IsTemplate     bool          `json:"is_template,omitempty"`
	Arch           string        `json:"arch,omitempty"`
	VrType         string        `json:"vr_type,omitempty"`
	DatacenterID   string        `json:"datacenter_id,omitempty"`
	DatacenterName string        `json:"datacenter_name,omitempty"`
	SnapshotNum    int64         `json:"snapshot_num,omitempty"`
	IsBindingHost  bool          `json:"is_binding_host,omitempty"`
	CPU            *VMCpu        `json:"cpu,omitempty"`
	Memory         *VMMemory     `json:"memory,omitempty"`
	VmNics         []VmNicDetail `json:"vm_nics,omitempty"`
	VmDisks        []VDisk       `json:"vm_disks,omitempty"`
	PvDriverStatus string        `json:"pv_driver_status,omitempty"`
	OsOptions      *OsOptions    `json:"os_options,omitempty"`
}

// VMCpu represents CPU configuration for a VM.
type VMCpu struct {
	Quantity    int `json:"quantity,omitempty"`     // Total CPU cores, range [1,255]
	HotPlug     int `json:"cpu.hot_plug,omitempty"` // 0=disabled, 1=enabled (default 1)
	Weight      int `json:"weight,omitempty"`       // CPU shares, range [2,255000]
	Reservation int `json:"reservation,omitempty"`  // Reserved CPU MHz
	Limit       int `json:"limit,omitempty"`        // Max CPU MHz
}

// VMMemory represents memory configuration for a VM.
type VMMemory struct {
	QuantitySize int `json:"quantity_size,omitempty"` // Total memory in MB, range [128,4194304]
	HotPlug      int `json:"mem.hot_plug,omitempty"`  // 0=disabled, 1=enabled (default 1)
	Reservation  int `json:"reservation,omitempty"`   // Reserved memory MB
	Limit        int `json:"limit,omitempty"`         // Max memory MB
}

// VmNicDetail represents a NIC on a VM as returned in the GET response.
type VmNicDetail struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	PID               string `json:"pid,omitempty"`
	Type              int    `json:"type,omitempty"`
	Vringbuf          int    `json:"vringbuf,omitempty"`
	Queues            int    `json:"queues,omitempty"`
	VmID              string `json:"vm_id,omitempty"`
	VmName            string `json:"vm_name,omitempty"`
	MacAddress        string `json:"mac_address,omitempty"`
	IPAddress         string `json:"ip_address,omitempty"`
	PortGroupID       string `json:"port_group_id,omitempty"`
	PortGroupName     string `json:"port_group_name,omitempty"`
	DvsID             string `json:"dvs_id,omitempty"`
	DvsName           string `json:"dvs_name,omitempty"`
	DvsType           string `json:"dvs_type,omitempty"`
	BootOrder         int    `json:"boot_order,omitempty"`
	SiteID            string `json:"site_id,omitempty"`
	VspID             string `json:"vsp_id,omitempty"`
	EnableSecurityGrp bool   `json:"enable_security_group,omitempty"`
	SecurityGroupID   int    `json:"security_group_id,omitempty"`
	SecurityGroupName string `json:"security_group_name,omitempty"`
	NicType           int    `json:"nic_type,omitempty"`
}

// VDisk represents a virtual disk on a VM.
type VDisk struct {
	ID                 string  `json:"id,omitempty"`
	PID                string  `json:"pid,omitempty"`
	Status             string  `json:"status,omitempty"`
	Type               string  `json:"type,omitempty"`
	Name               string  `json:"name,omitempty"`
	TotalCapacity      float64 `json:"total_capacity,omitempty"`
	UsedCapacity       float64 `json:"used_capacity,omitempty"`
	ProvisionCapacity  float64 `json:"provision_capacity,omitempty"`
	VmID               string  `json:"vm_id,omitempty"`
	VmName             string  `json:"vm_name,omitempty"`
	PciType            string  `json:"pci_type,omitempty"`
	ScsiCmdPassthrough bool    `json:"scsi_command_passthrough,omitempty"`
	VdiskSettingsType  int     `json:"vdisk_settings_type,omitempty"`
	IsSystemVdisk      bool    `json:"is_system_vdisk,omitempty"`
	IsThin             bool    `json:"is_thin,omitempty"`
	DatastoreID        string  `json:"datastore_id,omitempty"`
	DatastoreType      string  `json:"datastore_type,omitempty"`
	VmNum              int     `json:"vm_num,omitempty"`
	IsIndepDisk        bool    `json:"is_indep_disk,omitempty"`
	IsPersistentDisk   bool    `json:"is_persistent_disk,omitempty"`
	IOmode             string  `json:"io_mode,omitempty"`
	SequenceNum        int  `json:"sequence_num,omitempty"`
}

// OsOptions represents OS configuration for a VM.
type OsOptions struct {
	OsVersion   string `json:"os_version,omitempty"`    // OS version number
	OsType      string `json:"os_type,omitempty"`       // Windows, Linux, Other
	GuestOSName string `json:"guest_os_name,omitempty"` // Guest OS name
}

// CloneRequest is the request body for cloning or creating VM from template.
// POST /rest/vmmgmt/v1/vms/{id}/clone
type CloneRequest struct {
	Name               string                `json:"name"`                         // Required: VM name, 1~256 chars
	Description        string                `json:"description,omitempty"`        // Optional: VM description
	Location           string                `json:"location"`                     // Required: cluster_id or host_id
	IsBindingHost      bool                  `json:"is_binding_host,omitempty"`    // Default: false
	VmNum              int                   `json:"vm_num,omitempty"`             // Default: 1, range [1,50]
	NamePostfixIndex   int                   `json:"name_postfix_index,omitempty"` // Default: 0
	IsAutoBoot         bool                  `json:"is_auto_boot,omitempty"`       // Default: false (not started after creation)
	GuestCustomization *GuestOSCustomization `json:"guest_customization,omitempty"`
	VmConfig           *VmSimpleConfig       `json:"vm_config,omitempty"`
	OsOptions          *OsOptions            `json:"os_options,omitempty"`
}

// GuestOSCustomization represents guest OS customization during VM creation.
type GuestOSCustomization struct {
	OsType             string    `json:"os_type,omitempty"`
	Hostname           string    `json:"hostname,omitempty"`
	Password           string    `json:"password,omitempty"`
	Workgroup          string    `json:"workgroup,omitempty"`
	Domain             string    `json:"domain,omitempty"`
	DomainName         string    `json:"domain_name,omitempty"`
	DomainPassword     string    `json:"domain_password,omitempty"`
	OuName             string    `json:"ou_name,omitempty"`
	IsUpdateVmPassword bool      `json:"is_update_vm_password,omitempty"`
	VmCustomizationID  int64     `json:"vm_customization_id,omitempty"`
	IsUpdateVmSID      bool      `json:"is_update_vm_sid,omitempty"`
	NicSpecification   []NicInfo `json:"nic_specification,omitempty"`
}

// NicInfo represents NIC information in guest OS customization.
type NicInfo struct {
	SequenceNum           string       `json:"sequence_num,omitempty"`
	IpVersion             int          `json:"ip_version,omitempty"`
	Ip                    string       `json:"ip,omitempty"`
	Netmask               string       `json:"netmask,omitempty"`
	Gateway               string       `json:"gateway,omitempty"`
	SetDNS                string       `json:"setdns,omitempty"`
	AddDNS                string       `json:"adddns,omitempty"`
	IsIpv6DhcpEnabled     bool         `json:"is_ipv6_dhcp_enabled,omitempty"`
	IsIpv6AutoConfEnabled bool         `json:"is_ipv6_auto_conf_enabled,omitempty"`
	Ipv6Address           []IpAddress6 `json:"ipv6_address,omitempty"`
	Ipv6Gateway           string       `json:"ipv6_gateway,omitempty"`
	Ipv6SetDNS            string       `json:"ipv6_setdns,omitempty"`
	Ipv6AddDNS            string       `json:"ipv6_adddns,omitempty"`
}

// IpAddress6 represents an IPv6 address.
type IpAddress6 struct {
	Ipv6Addr   string `json:"ipv6_addr,omitempty"`
	Ipv6Prefix int    `json:"ipv6_prefix,omitempty"`
}

// VmSimpleConfig represents VM configuration during creation.
type VmSimpleConfig struct {
	CPU    *VMCpu    `json:"cpu,omitempty"`
	Memory *VMMemory `json:"memory,omitempty"`
	Disks  []VMDisk  `json:"disks,omitempty"`
	Nics   []VMNic   `json:"nics,omitempty"`
}

// VMDisk represents a disk in VM creation config.
type VMDisk struct {
	SequenceNum    int    `json:"sequence_num"`              // Required: slot number
	DatastoreID    string `json:"datastore_id"`              // Required: datastore ID
	Quantity       int    `json:"quantity"`                  // Required: disk capacity in GB
	IndepDisk      bool   `json:"indep_disk,omitempty"`      // Default: true
	PersistentDisk bool   `json:"persistent_disk,omitempty"` // Default: true
	IsThin         bool   `json:"is_thin,omitempty"`         // Default: true
	PciType        string `json:"pci_type"`                  // Required: IDE, SCSI, VIRTIO
	VolType        int    `json:"vol_type,omitempty"`        // 0=normal, 1=delayed zero (default 0)
}

// VMNic represents a NIC in VM creation config.
type VMNic struct {
	PortGroupID string `json:"port_group_id"` // Required: port group ID
}

// UpdateVmRequest is the request body for updating a VM.
// PUT /rest/vmmgmt/v1/vms/{id}
type UpdateVmRequest struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	IsTemplate  *bool      `json:"is_template,omitempty"`
	CPU         *VMCpu     `json:"cpu,omitempty"`
	Memory      *VMMemory  `json:"memory,omitempty"`
	OsOptions   *OsOptions `json:"os_options,omitempty"`
}

// PowerActionRequest is the request body for power on/off/reboot operations.
type PowerActionRequest struct {
	VmIDs []string `json:"vm_ids"`         // Required: list of VM IDs
	Mode  string   `json:"mode,omitempty"` // Optional: "safe" or "force" (for power-off/reboot)
}

// BatchDeleteRequest is the request body for batch deleting VMs.
type BatchDeleteRequest struct {
	VmIDs []string `json:"vm_ids"` // Required: list of VM IDs to delete
}

// PowerAction represents the type of power operation.
type PowerAction string

const (
	PowerOn  PowerAction = "power-on"
	PowerOff PowerAction = "power-off"
	Reboot   PowerAction = "reboot"
)

// AffectedResource represents an affected resource in a task detail response.
type AffectedResource struct {
	Operate string `json:"operate"` // CREATE, MODIFY, DELETE
	Type    string `json:"type"`    // Resource type
	ID      string `json:"id"`      // Resource ID
	Name    string `json:"name"`    // Resource name
}

// VmListRequest is the request body for listing/querying VMs.
// POST /rest/vmmgmt/v1/vms/query
type VmListRequest struct {
	SiteID       string   `json:"site_id,omitempty"`       // Site ID
	ClusterID    string   `json:"cluster_id,omitempty"`    // Cluster ID (not supported for HCS)
	DcID         string   `json:"dc_id,omitempty"`         // Datacenter ID (FusionCompute only)
	ClusterName  string   `json:"cluster_name,omitempty"`  // Cluster name, supports fuzzy match (not HCS)
	HostID       string   `json:"host_id,omitempty"`       // Host ID
	HostName     string   `json:"host_name,omitempty"`     // Host name, supports fuzzy match
	Name         string   `json:"name,omitempty"`          // VM name, supports fuzzy match
	IPAddress    string   `json:"ip_address,omitempty"`    // VM IP, supports fuzzy match
	Status       []string `json:"status,omitempty"`        // VM status list (max 20 items)
	IsTemplate   *bool    `json:"is_template,omitempty"`   // Whether it is a template
	OsType       []string `json:"os_type,omitempty"`       // OS type list (max 3 items): Windows, Linux, Other
	VrType       string   `json:"vr_type,omitempty"`       // FUSIONCOMPUTE, VMWARE, HCS
	DatacenterID string   `json:"datacenter_id,omitempty"` // Datacenter ID (VMWARE only)
	SortDir      string   `json:"sort_dir,omitempty"`      // asc or desc
	SortKey      string   `json:"sort_key,omitempty"`      // name, cpu_core, memory_size, disk_total_size, create_time, ip_address
	PageNo       int32    `json:"page_no,omitempty"`       // Page number, default 1
	PageSize     int32    `json:"page_size,omitempty"`     // Items per page, default 20, max 1000
}

// VmsResponse is the response body from listing/querying VMs.
type VmsResponse struct {
	Total int  `json:"total,omitempty"`
	VMs   []Vm `json:"vms,omitempty"`
}

// TaskDetail represents a task as returned by GET /rest/taskmgmt/v1/tasks/{task_id}.
type TaskDetail struct {
	ID                string             `json:"id,omitempty"`
	NameEn            string             `json:"name_en,omitempty"`
	NameCn            string             `json:"name_cn,omitempty"`
	Description       string             `json:"description,omitempty"`
	ParentID          string             `json:"parent_id,omitempty"`
	SeqNo             int                `json:"seq_no,omitempty"`
	Status            int                `json:"status,omitempty"` // 1=initial, 2=running, 3=success, 4=partial_success, 5=failed, 6=timeout
	Progress          int                `json:"progress,omitempty"`
	OwnerName         string             `json:"owner_name,omitempty"`
	OwnerID           string             `json:"owner_id,omitempty"`
	CreateTime        int64              `json:"create_time,omitempty"`
	StartTime         int64              `json:"start_time,omitempty"`
	EndTime           int64              `json:"end_time,omitempty"`
	DetailEn          string             `json:"detail_en,omitempty"`
	DetailCn          string             `json:"detail_cn,omitempty"`
	IsSupportfulRetry bool               `json:"is_supportful_retry,omitempty"`
	IsSupportRollback bool               `json:"is_supportrollback,omitempty"`
	Remarks           string             `json:"remarks,omitempty"`
	ExecuteHistoryCnt int32              `json:"execute_history_count,omitempty"`
	Resources         []AffectedResource `json:"resources,omitempty"`
}
