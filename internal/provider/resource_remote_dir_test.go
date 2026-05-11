package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRemoteDir(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "remote_dir" "resource_1" {
					conn {
						host = "remotehost"
						user = "root"
						sudo = true
						password = "password"
						timeout = 1000
					}
					path = "/tmp/resource_1"
					permissions = "0777"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_dir.resource_1", "id", regexp.MustCompile("remotehost:22:/tmp/resource_1")),
				),
			},
		},
	})
}

func TestAccResourceRemoteDirWithDefaultConnection(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "remote_dir" "resource_2" {
					provider = remotehost

					path = "/tmp/resource_2"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_dir.resource_2", "id", regexp.MustCompile("remotehost:22:/tmp/resource_2")),
				),
			},
		},
	})
}

func TestAccResourceRemoteDirWithAgent(t *testing.T) {
	if os.Getenv("SKIP_TEST_AGENT") == "1" {
		return
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "remote_dir" "resource_3" {
					conn {
						host = "remotehost"
						user = "root"
						agent = true
					}
					path = "/tmp/resource_3"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_dir.resource_3", "id", regexp.MustCompile("remotehost:22:/tmp/resource_3")),
				),
			},
		},
	})
}

func TestAccResourceRemoteDirOwnership(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "remote_dir" "resource_4" {
					conn {
						host = "remotehost"
						user = "root"
						sudo = true
						password = "password"
					}
					path = "/tmp/resource_4"
					permissions = "0777"
					owner = "1000"
					group = "1001"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_dir.resource_4", "owner_name", regexp.MustCompile("")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_4", "group_name", regexp.MustCompile("")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_4", "owner", regexp.MustCompile("1000")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_4", "group", regexp.MustCompile("1001")),
				),
			},
		},
	})
}

func TestAccResourceRemoteDirOwnershipWithDefaultConnection(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "remote_dir" "resource_5" {
					provider = remotehost
					path = "/tmp/resource_5"
					owner = "1000"
					group = "1001"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_dir.resource_5", "owner_name", regexp.MustCompile("")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_5", "group_name", regexp.MustCompile("")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_5", "owner", regexp.MustCompile("1000")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_5", "group", regexp.MustCompile("1001")),
				),
			},
		},
	})
}

func TestAccResourceRemoteDirOwnershipNames(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "remote_dir" "resource_6" {
					conn {
						host = "remotehost"
						user = "root"
						sudo = true
						password = "password"
					}
					path = "/tmp/resource_6"
					permissions = "0777"
					owner_name = "root"
					group_name = "root"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"remote_dir.resource_6", "owner_name", regexp.MustCompile("root")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_6", "group_name", regexp.MustCompile("root")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_6", "owner", regexp.MustCompile("")),
					resource.TestMatchResourceAttr(
						"remote_dir.resource_6", "group", regexp.MustCompile("")),
				),
			},
		},
	})
}
