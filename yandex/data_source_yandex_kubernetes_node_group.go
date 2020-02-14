package yandex

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/k8s/v1"
	"github.com/yandex-cloud/go-sdk/sdkresolvers"
)

func dataSourceYandexKubernetesNodeGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceYandexKubernetesNodeGroupRead,
		Schema: map[string]*schema.Schema{
			"node_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_template": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"cores": {
										Type:     schema.TypeInt,
										Computed: true,
									},

									"core_fraction": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"boot_disk": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"platform_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"scheduling_policy": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"preemptible": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"scale_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fixed_scale": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"auto_scale": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"initial": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"allocation_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"instance_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_info": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_revision_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"new_revision_summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_deprecated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"maintenance_policy": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_upgrade": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_repair": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"maintenance_window": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      dayOfWeekHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"duration": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceYandexKubernetesNodeGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	err := checkOneOf(d, "node_group_id", "name")
	if err != nil {
		return err
	}

	nodeGroupID := d.Get("node_group_id").(string)
	_, nodeGroupNameOk := d.GetOk("name")

	if nodeGroupNameOk {
		nodeGroupID, err = resolveObjectID(ctx, config, d, sdkresolvers.KubernetesNodeGroupResolver)
		if err != nil {
			return fmt.Errorf("failed to resolve data source node-group by name: %v", err)
		}
	}

	ng, err := config.sdk.Kubernetes().NodeGroup().Get(ctx, &k8s.GetNodeGroupRequest{
		NodeGroupId: nodeGroupID,
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Kubernetes node-group with ID %q", nodeGroupID))
	}

	err = flattenNodeGroupSchemaData(ng, d)
	if err != nil {
		return fmt.Errorf("failed to fill Kubernetes node-group shema: %v", err)
	}

	d.Set("node_group_id", ng.Id)
	return nil
}