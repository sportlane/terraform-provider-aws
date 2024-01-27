// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package qbusiness

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/qbusiness"
	"github.com/aws/aws-sdk-go-v2/service/qbusiness/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_qbusiness_group", name="Group")
func ResourceGroup() *schema.Resource {
	return &schema.Resource{

		CreateWithoutTimeout: resourceGroupUpsert,
		ReadWithoutTimeout:   resourceGroupRead,
		UpdateWithoutTimeout: resourceGroupUpsert,
		DeleteWithoutTimeout: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the application in which the user and group mapping belongs.",
				ValidateFunc: validation.All(
					validation.StringMatch(regexache.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{35}$`), "must be a valid application ID"),
				),
			},
			"datasource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the data source for which you want to map users to their groups.",
				ValidateFunc: validation.All(
					validation.StringMatch(regexache.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{35}$`), "must be a valid datasource ID"),
				),
			},
			"member_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of sub-groups that belong to the group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the sub group.",
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 1024),
								validation.StringMatch(regexache.MustCompile(`^\P{C}*$`), "must not contain any control characters"),
							),
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Type of group.",
							ValidateDiagFunc: enum.Validate[types.MembershipType](),
						},
					},
				},
			},
			"member_users": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of sub-groups that belong to the group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Identifier of the user you want to map to a group.",
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 1024),
								validation.StringMatch(regexache.MustCompile(`^\P{C}*$`), "must not contain any control characters"),
							),
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Type of group.",
							ValidateDiagFunc: enum.Validate[types.MembershipType](),
						},
					},
				},
			},
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the group that is mapped to one or more users.",
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 1024),
					validation.StringMatch(regexache.MustCompile(`^\P{C}*$`), "must not contain any control characters"),
				),
			},
			"index_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the index in which you want to map users to their groups.",
				ValidateFunc: validation.All(
					validation.StringMatch(regexache.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{35}$`), "must be a valid index ID"),
				),
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Type of group.",
				ValidateDiagFunc: enum.Validate[types.MembershipType](),
			},
		},
	}
}

func resourceGroupUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).QBusinessClient(ctx)

	application_id := d.Get("application_id").(string)
	index_id := d.Get("index_id").(string)
	datasource_id := d.Get("datasource_id").(string)
	group_name := d.Get("group_name").(string)

	input := &qbusiness.PutGroupInput{
		ApplicationId: aws.String(application_id),
		IndexId:       aws.String(index_id),
		GroupName:     aws.String(group_name),
		DataSourceId:  aws.String(datasource_id),
		GroupMembers: &types.GroupMembers{
			MemberUsers:  expandMemberUsers(d.Get("member_users").([]interface{})),
			MemberGroups: expandMemberGroups(d.Get("member_groups").([]interface{})),
		},
		Type: types.MembershipType(d.Get("type").(string)),
	}

	_, err := conn.PutGroup(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Amazon Q group: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", application_id, index_id, group_name, datasource_id))

	return append(diags, resourceGroupRead(ctx, d, meta)...)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).QBusinessClient(ctx)

	application_id, index_id, group_name, datasource_id, err := parseGroupID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parse qbusiness group ID: %s", err)
	}

	input := &qbusiness.GetGroupInput{
		ApplicationId: aws.String(application_id),
		IndexId:       aws.String(index_id),
		GroupName:     aws.String(group_name),
		DataSourceId:  aws.String(datasource_id),
	}

	_, err = conn.GetGroup(ctx, input)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] qbusiness group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading qbusiness group (%s): %s", d.Id(), err)
	}

	d.Set("application_id", application_id)
	d.Set("index_id", index_id)
	d.Set("datasource_id", datasource_id)
	d.Set("group_name", group_name)

	// WE CANNOT SET MEMBER_USERS AND MEMBER_GROUPS AS THEY ARE NOT SUPPORTED BY THE SDK

	return nil
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).QBusinessClient(ctx)

	application_id, index_id, group_name, datasource_id, err := parseGroupID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parse qbusiness group ID: %s", err)
	}

	input := &qbusiness.DeleteGroupInput{
		ApplicationId: aws.String(application_id),
		IndexId:       aws.String(index_id),
		GroupName:     aws.String(group_name),
		DataSourceId:  aws.String(datasource_id),
	}

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil
	}

	_, err = conn.DeleteGroup(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Amazon Q group: %s", err)
	}

	return nil
}

func parseGroupID(id string) (string, string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 4 {
		return "", "", "", "", fmt.Errorf("unexpected format of group ID (%s), expected APP_ID/INDEX_ID/GROUP_NAME/DATASOURCE_ID", id)
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

func FindGroupByID(ctx context.Context, conn *qbusiness.Client, id string) (*qbusiness.GetGroupOutput, error) {
	application_id, index_id, group_name, datasource_id, err := parseGroupID(id)

	if err != nil {
		return nil, err
	}

	input := &qbusiness.GetGroupInput{
		ApplicationId: aws.String(application_id),
		IndexId:       aws.String(index_id),
		GroupName:     aws.String(group_name),
		DataSourceId:  aws.String(datasource_id),
	}

	output, err := conn.GetGroup(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func expandMemberUsers(memberUsers []interface{}) []types.MemberUser {
	if len(memberUsers) == 0 {
		return nil
	}
	result := make([]types.MemberUser, len(memberUsers))
	for i, memberUser := range memberUsers {
		memberUserMap := memberUser.(map[string]interface{})
		result[i] = types.MemberUser{
			UserId: aws.String(memberUserMap["user_id"].(string)),
			Type:   types.MembershipType(memberUserMap["type"].(string)),
		}
	}
	return result
}

func expandMemberGroups(memberGroups []interface{}) []types.MemberGroup {
	if len(memberGroups) == 0 {
		return nil
	}

	result := make([]types.MemberGroup, len(memberGroups))
	for i, memberGroup := range memberGroups {
		memberGroupMap := memberGroup.(map[string]interface{})
		result[i] = types.MemberGroup{
			GroupName: aws.String(memberGroupMap["group_name"].(string)),
			Type:      types.MembershipType(memberGroupMap["type"].(string)),
		}
	}
	return result
}

func flattenMemberUsers(memberUsers []types.MemberUser) []interface{} {
	if len(memberUsers) == 0 {
		return nil
	}

	result := make([]interface{}, len(memberUsers))
	for i, memberUser := range memberUsers {
		result[i] = map[string]interface{}{
			"user_id": memberUser.UserId,
			"type":    memberUser.Type,
		}
	}

	return result
}

func flattenMemberGroups(memberGroups []types.MemberGroup) []interface{} {
	if len(memberGroups) == 0 {
		return nil
	}

	result := make([]interface{}, len(memberGroups))
	for i, memberGroup := range memberGroups {
		result[i] = map[string]interface{}{
			"group_name": memberGroup.GroupName,
			"type":       memberGroup.Type,
		}
	}

	return result
}
