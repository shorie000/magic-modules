package firebaserules_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirebaserulesRelease_updateRulesetInPlace(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
		"allow":         "false",
	}
	updatedContext := map[string]interface{}{
		"project_name":  context["project_name"],
		"random_suffix": context["random_suffix"],
		"allow":         "request.auth != null",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirebaserulesReleaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaserulesRelease_updateRulesetInPlace(context),
			},
			{
				ResourceName:            "google_firebaserules_release.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_ruleset_in_place"},
			},
			{
				Config: testAccFirebaserulesRelease_updateRulesetInPlace(updatedContext),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_firebaserules_release.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_firebaserules_release.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_ruleset_in_place"},
			},
		},
	})
}

func testAccFirebaserulesRelease_updateRulesetInPlace(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firebaserules_release" "primary" {
  name                    = "cloud.firestore/tf-test-database%{random_suffix}"
  project                 = "%{project_name}"
  ruleset_name            = "projects/%{project_name}/rulesets/${google_firebaserules_ruleset.firestore.name}"
  update_ruleset_in_place = true
}

resource "google_firebaserules_ruleset" "firestore" {
  project = "%{project_name}"

  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if %{allow}; } } }"
      name    = "firestore.rules"
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}
`, context)
}
