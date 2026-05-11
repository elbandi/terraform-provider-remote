package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"remote": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
	"remotehost": func() (*schema.Provider, error) {
		provider := New("dev")()
		configureProvider := provider.ConfigureContextFunc
		provider.ConfigureContextFunc = func(c context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
			if err := rd.Set("conn", []interface{}{
				map[string]interface{}{
					"host":     "remotehost",
					"user":     "root",
					"password": "password",
					"port":     22,
				},
			}); err != nil {
				return nil, diag.FromErr(err)
			}
			return configureProvider(c, rd)
		}
		return provider, nil
	},
	"remotehost2": func() (*schema.Provider, error) {
		provider := New("dev")()
		configureProvider := provider.ConfigureContextFunc
		provider.ConfigureContextFunc = func(c context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
			if err := rd.Set("conn", []interface{}{
				map[string]interface{}{
					"host":     "remotehost2",
					"user":     "root",
					"password": "password",
					"port":     22,
				},
			}); err != nil {
				return nil, diag.FromErr(err)
			}
			return configureProvider(c, rd)
		}
		return provider, nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func checkFileExists(host string, destinationFilePath string, should bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		exists, err := checkExist(host, destinationFilePath, "f")
		if err != nil {
			return fmt.Errorf("Error occurred while retrieving file info at path: %s\n, error: %s\n",
				destinationFilePath, err)
		}

		if exists != should {
			return fmt.Errorf(
				"File not exists.\nexpected:%t\ngot: %t\n",
				should, exists)
		}

		return nil
	}
}

func checkDirExists(host string, destinationFilePath string, should bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		exists, err := checkExist(host, destinationFilePath, "d")
		if err != nil {
			return fmt.Errorf("Error occurred while retrieving directory info at path: %s\n, error: %s\n",
				destinationFilePath, err)
		}

		if exists != should {
			return fmt.Errorf(
				"Directory not exists.\nexpected:%t\ngot: %t\n",
				should, exists)
		}

		return nil
	}
}
