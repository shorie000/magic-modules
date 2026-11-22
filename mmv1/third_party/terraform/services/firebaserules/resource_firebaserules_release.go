package firebaserules

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceFirebaserulesRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirebaserulesReleaseCreate,
		Read:   resourceFirebaserulesReleaseRead,
		Update: resourceFirebaserulesReleaseUpdate,
		Delete: resourceFirebaserulesReleaseDelete,

		Importer: &schema.ResourceImporter{
			State: resourceFirebaserulesReleaseImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			tpgresource.DefaultProviderDeletionPolicy("DELETE"),
			customdiff.ForceNewIf("ruleset_name", func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
				return !d.Get("update_ruleset_in_place").(bool)
			}),
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Format: `projects/{project_id}/releases/{release_id}`\\Firestore Rules Releases will **always** have the name 'cloud.firestore'",
			},

			"ruleset_name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Name of the `Ruleset` referred to by this `Release`. The `Ruleset` must exist for the `Release` to be created. Changing this field replaces the `Release` unless `update_ruleset_in_place` is true.",
			},

			"update_ruleset_in_place": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, changing `ruleset_name` updates the `Release` in place instead of deleting and recreating it, so the `Release` never stops being served. The new `Ruleset` must exist before the `Release` is updated, so set `create_before_destroy` on the `Ruleset` it refers to. Defaults to false.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time the release was created.",
			},

			"disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Disable the release to keep it from being served. The response code of NOT_FOUND will be given for executables generated from this Release.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time the release was updated.",
			},
			//UDP schema start
			"deletion_policy": tpgresource.DeletionPolicySchemaEntry("DELETE"),
			//UDP schema end
		},
	}
}

func resourceFirebaserulesReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &Release{
		Name:        dcl.String(d.Get("name").(string)),
		RulesetName: dcl.String(d.Get("ruleset_name").(string)),
		Project:     dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLFirebaserulesClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyRelease(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Release: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Release %q: %#v", d.Id(), res)

	return resourceFirebaserulesReleaseRead(d, meta)
}

func resourceFirebaserulesReleaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &Release{
		Name:        dcl.String(d.Get("name").(string)),
		RulesetName: dcl.String(d.Get("ruleset_name").(string)),
		Project:     dcl.String(project),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLFirebaserulesClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetRelease(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("FirebaserulesRelease %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("ruleset_name", res.RulesetName); err != nil {
		return fmt.Errorf("error setting ruleset_name in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("disabled", res.Disabled); err != nil {
		return fmt.Errorf("error setting disabled in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	if err := tpgresource.DeletionPolicyReadDefault(d, config, "DELETE"); err != nil {
		return err
	}

	return nil
}

// UDP update start
func resourceFirebaserulesReleaseUpdate(d *schema.ResourceData, meta interface{}) error {
	// Besides "ruleset_name", only the root field "deletion_policy" and virtual fields are mutable
	if !d.HasChange("ruleset_name") {
		return resourceFirebaserulesReleaseRead(d, meta)
	}

	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &Release{
		Name:        dcl.String(d.Get("name").(string)),
		RulesetName: dcl.String(d.Get("ruleset_name").(string)),
		Project:     dcl.String(project),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLFirebaserulesClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyRelease(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		return fmt.Errorf("Error updating Release: %s", err)
	}

	log.Printf("[DEBUG] Finished updating Release %q: %#v", d.Id(), res)

	return resourceFirebaserulesReleaseRead(d, meta)
}

//UDP update end

func resourceFirebaserulesReleaseDelete(d *schema.ResourceData, meta interface{}) error {

	if ok, err := tpgresource.DeletionPolicyPreDelete(d); err != nil {
		return err
	} else if ok {
		return nil
	}

	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &Release{
		Name:        dcl.String(d.Get("name").(string)),
		RulesetName: dcl.String(d.Get("ruleset_name").(string)),
		Project:     dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Release %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLFirebaserulesClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteRelease(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Release: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Release %q", d.Id())
	return nil
}

func resourceFirebaserulesReleaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>.+)/releases/(?P<name>.+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/releases/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Explicitly set virtual fields to default values on import
	if err := d.Set("update_ruleset_in_place", false); err != nil {
		return nil, fmt.Errorf("Error setting update_ruleset_in_place: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
