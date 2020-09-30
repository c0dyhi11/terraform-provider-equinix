package equinix

import (
	"fmt"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var neDevicePlatformSchemaNames = map[string]string{
	"DeviceTypeCode":  "device_type",
	"Flavor":          "flavor",
	"CoreCount":       "core_count",
	"Memory":          "memory",
	"MemoryUnit":      "memory_unit",
	"PackageCodes":    "packages",
	"ManagementTypes": "management_types",
	"LicenseOptions":  "license_options",
}

func dataSourceNeDevicePlatform() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNeDevicePlatformRead,
		Schema: map[string]*schema.Schema{
			neDevicePlatformSchemaNames["DeviceTypeCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			neDevicePlatformSchemaNames["Flavor"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"small", "medium", "large", "xlarge"}, false),
			},
			neDevicePlatformSchemaNames["CoreCount"]: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			neDevicePlatformSchemaNames["Memory"]: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			neDevicePlatformSchemaNames["MemoryUnit"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDevicePlatformSchemaNames["PackageCodes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			neDevicePlatformSchemaNames["ManagementTypes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"EQUINIX-CONFIGURED", "SELF-CONFIGURED"}, false),
				},
			},
			neDevicePlatformSchemaNames["LicenseOptions"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"BYOL", "Sub"}, false),
				},
			},
		},
	}
}

func dataSourceNeDevicePlatformRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	typeCode := d.Get(neDevicePlatformSchemaNames["DeviceTypeCode"]).(string)
	platforms, err := conf.ne.GetDevicePlatforms(typeCode)
	if err != nil {
		return err
	}
	var filtered []ne.DevicePlatform
	for _, platform := range platforms {
		if v, ok := d.GetOk(neDevicePlatformSchemaNames["Flavor"]); ok && platform.Flavor != v.(string) {
			continue
		}
		if v, ok := d.GetOk(neDevicePlatformSchemaNames["CoreCount"]); ok && platform.CoreCount != v.(int) {
			continue
		}
		if v, ok := d.GetOk(neDevicePlatformSchemaNames["PackageCodes"]); ok {
			pkgCodes := expandSetToStringList(v.(*schema.Set))
			if !stringsFound(pkgCodes, platform.PackageCodes) {
				continue
			}
		}
		if v, ok := d.GetOk(neDevicePlatformSchemaNames["ManagementTypes"]); ok {
			mgmtTypes := expandSetToStringList(v.(*schema.Set))
			if !stringsFound(mgmtTypes, platform.ManagementTypes) {
				continue
			}
		}
		if v, ok := d.GetOk(neDevicePlatformSchemaNames["LicenseOptions"]); ok {
			licOptions := expandSetToStringList(v.(*schema.Set))
			if !stringsFound(licOptions, platform.LicenseOptions) {
				continue
			}
		}
		filtered = append(filtered, platform)
	}
	if len(filtered) < 1 {
		return fmt.Errorf("device platform query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return fmt.Errorf("device platform query returned more than one result, please try more specific search criteria")
	}
	return updateNeDevicePlatformResource(filtered[0], typeCode, d)
}

func updateNeDevicePlatformResource(platform ne.DevicePlatform, typeCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", typeCode, platform.Flavor))
	if err := d.Set(neDevicePlatformSchemaNames["Flavor"], platform.Flavor); err != nil {
		return fmt.Errorf("error reading Flavor: %s", err)
	}
	if err := d.Set(neDevicePlatformSchemaNames["CoreCount"], platform.CoreCount); err != nil {
		return fmt.Errorf("error reading CoreCount: %s", err)
	}
	if err := d.Set(neDevicePlatformSchemaNames["Memory"], platform.Memory); err != nil {
		return fmt.Errorf("error reading Memory: %s", err)
	}
	if err := d.Set(neDevicePlatformSchemaNames["MemoryUnit"], platform.MemoryUnit); err != nil {
		return fmt.Errorf("error reading MemoryUnit: %s", err)
	}
	if err := d.Set(neDevicePlatformSchemaNames["PackageCodes"], platform.PackageCodes); err != nil {
		return fmt.Errorf("error reading PackageCodes: %s", err)
	}
	if err := d.Set(neDevicePlatformSchemaNames["ManagementTypes"], platform.ManagementTypes); err != nil {
		return fmt.Errorf("error reading ManagementTypes: %s", err)
	}
	if err := d.Set(neDevicePlatformSchemaNames["LicenseOptions"], platform.LicenseOptions); err != nil {
		return fmt.Errorf("error reading LicenseOptions: %s", err)
	}
	return nil
}