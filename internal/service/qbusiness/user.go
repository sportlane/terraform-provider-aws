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
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_qbusiness_user", name="User")
// @Tags(identifierAttribute="arn")
func ResourceUser() *schema.Resource {
	return &schema.Resource{

		CreateWithoutTimeout: resourceUserCreate,
		ReadWithoutTimeout:   resourceUserRead,
		UpdateWithoutTimeout: resourceUserUpdate,
		DeleteWithoutTimeout: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the Amazon Q application associated with the user.",
				ValidateFunc: validation.All(
					validation.StringMatch(regexache.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{35}$`), "must be a valid application ID"),
				),
			},
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "User email attached to a user mapping.",
				ValidateFunc: validation.StringLenBetween(0, 2048),
			},
			"user_aliases": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "List of user aliases attached to a user mapping.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 100,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "Identifier of the user id associated with the user aliases",
										ValidateFunc: validation.StringLenBetween(1, 2048),
									},
									"datasource_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Identifier of the data source that the user aliases are associated with.",
										ValidateFunc: validation.All(
											validation.StringMatch(regexache.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{35}$`), "must be a valid datasource ID"),
										),
									},
									"index_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Identifier of the index that the user aliases are associated with.",
										ValidateFunc: validation.All(
											validation.StringMatch(regexache.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{35}$`), "must be a valid index ID"),
										),
									},
								},
							},
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).QBusinessClient(ctx)

	application_id := d.Get("application_id").(string)
	user_id := d.Get("user_id").(string)

	input := &qbusiness.CreateUserInput{
		ApplicationId: aws.String(application_id),
		UserId:        aws.String(user_id),
	}

	if v, ok := d.GetOk("user_aliases"); ok {
		input.UserAliases = expandUserAliases(v.(*schema.Set).List())
	}

	_, err := conn.CreateUser(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Amazon Q user: %s", err)
	}

	d.SetId(application_id + "/" + user_id)

	return append(diags, resourceUserRead(ctx, d, meta)...)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).QBusinessClient(ctx)

	output, err := FindUserByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] qbusiness user (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading qbusiness datasource (%s): %s", d.Id(), err)
	}

	application_id, user_id, _ := parseUserID(d.Id())

	d.Set("application_id", application_id)
	d.Set("user_id", user_id)
	d.Set("user_aliases", flattenUserAliases(output.UserAliases))

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func parseUserID(id string) (string, string, error) {
	parts := strings.Split(id, "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid user ID: %s", id)
	}

	return parts[0], parts[1], nil
}

func FindUserByID(ctx context.Context, conn *qbusiness.Client, id string) (*qbusiness.GetUserOutput, error) {
	application_id, user_id, err := parseUserID(id)
	if err != nil {
		return nil, err
	}

	input := &qbusiness.GetUserInput{
		ApplicationId: aws.String(application_id),
		UserId:        aws.String(user_id),
	}

	output, err := conn.GetUser(ctx, input)

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

func flattenUserAliases(v []types.UserAlias) []interface{} {
	if v == nil {
		return nil
	}

	res := make([]interface{}, 0, len(v))
	for _, r := range v {
		res = append(res, map[string]interface{}{
			"user_id":       aws.ToString(r.UserId),
			"datasource_id": aws.ToString(r.DataSourceId),
			"index_id":      aws.ToString(r.IndexId),
		})
	}
	return res
}

func expandUserAliases(v []interface{}) []types.UserAlias {
	if len(v) == 0 {
		return nil
	}

	res := make([]types.UserAlias, 0, len(v))
	for _, r := range v {
		m := r.(map[string]interface{})
		res = append(res, types.UserAlias{
			UserId:       aws.String(m["user_id"].(string)),
			DataSourceId: aws.String(m["datasource_id"].(string)),
			IndexId:      aws.String(m["index_id"].(string)),
		})
	}
	return res
}
