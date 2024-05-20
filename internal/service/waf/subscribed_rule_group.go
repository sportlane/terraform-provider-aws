// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package waf

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	DSNameSubscribedRuleGroup = "Subscribed Rule Group Data Source"
)

// @SDKDataSource("aws_waf_subscribed_rule_group")
func DataSourceSubscribedRuleGroup() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSubscribedRuleGroupRead,

		Schema: map[string]*schema.Schema{
			names.AttrName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			names.AttrMetricName: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceSubscribedRuleGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).WAFConn(ctx)
	name, nameOk := d.Get(names.AttrName).(string)
	metricName, metricNameOk := d.Get(names.AttrMetricName).(string)

	// Error out if string-assertion fails for either name or metricName
	if !nameOk || !metricNameOk {
		if !nameOk {
			name = DSNameSubscribedRuleGroup
		}

		err := errors.New("unable to read attributes")
		return create.DiagError(names.WAF, create.ErrActionReading, DSNameSubscribedRuleGroup, name, err)
	}

	output, err := FindSubscribedRuleGroupByNameOrMetricName(ctx, conn, name, metricName)

	if err != nil {
		return create.DiagError(names.WAF, create.ErrActionReading, DSNameSubscribedRuleGroup, name, err)
	}

	d.SetId(aws.StringValue(output.RuleGroupId))
	d.Set(names.AttrMetricName, output.MetricName)
	d.Set(names.AttrName, output.Name)

	return nil
}