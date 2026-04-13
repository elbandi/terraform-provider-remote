package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRemoteDir() *schema.Resource {
	return &schema.Resource{
		Description: "Directory on remote host.",

		ReadContext: dataSourceRemoteDirRead,

		Schema: map[string]*schema.Schema{
			"conn": {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "Connection to host where files are located.",
				Elem:        connectionSchemaResource,
			},
			"path": {
				Description: "Path to directory on remote host.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"check_only": {
				Description: "Only check if directory exists, do not read metadata.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"permissions": {
				Description: "Permissions of directory (in octal form).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"group": {
				Description: "Group ID (GID) of directory owner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"group_name": {
				Description: "Group name of directory owner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner": {
				Description: "User ID (UID) of directory owner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner_name": {
				Description: "User name of directory owner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceRemoteDirRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (error diag.Diagnostics) {
	conn, err := meta.(*apiClient).getConnWithDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setResourceID(d, conn); err != nil {
		return diag.FromErr(err)
	}

	client, err := meta.(*apiClient).getRemoteClient(ctx, conn)
	if err != nil {
		return diag.Errorf("unable to open remote client: %s", err.Error())
	}
	defer func() {
		if err := meta.(*apiClient).closeRemoteClient(conn); err != nil {
			error = append(error, diag.Errorf("unable to close remote client: %s", err.Error())...)
		}
	}()

	sudo, _, err := GetOk[bool](conn, "conn.0.sudo")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	path, err := Get[string](d, "path")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}
	checkOnly, err := Get[bool](d, "check_only")
	if err != nil {
		return diag.Diagnostics{{Severity: diag.Error, Summary: err.Error()}}
	}

	exists, err := client.DirExists(path, sudo)
	if err != nil {
		return diag.Errorf("unable to check if remote file exists: %s", err.Error())
	}
	if !exists {
		if !checkOnly {
			return diag.Errorf("remote directory does not exist")
		}
		return diag.Diagnostics{}
	}

	permissions, err := client.ReadFilePermissions(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file permissions: %s", err.Error())
	}
	if err := d.Set("permissions", permissions); err != nil {
		return diag.FromErr(err)
	}

	owner, err := client.ReadFileOwner(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file owner: %s", err.Error())
	}
	if err := d.Set("owner", owner); err != nil {
		return diag.FromErr(err)
	}

	ownerName, err := client.ReadFileOwnerName(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file owner_name: %s", err.Error())
	}
	if err := d.Set("owner_name", ownerName); err != nil {
		return diag.FromErr(err)
	}

	group, err := client.ReadFileGroup(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file group: %s", err.Error())
	}
	if err := d.Set("group", group); err != nil {
		return diag.FromErr(err)
	}

	groupName, err := client.ReadFileGroupName(path, sudo)
	if err != nil {
		return diag.Errorf("unable to read remote file group_name: %s", err.Error())
	}
	if err := d.Set("group_name", groupName); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
