// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package wafv2

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	wafv2_sdkv2 "github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceIPSet,
			TypeName: "aws_wafv2_ip_set",
			Name:     "IP Set",
		},
		{
			Factory:  dataSourceRegexPatternSet,
			TypeName: "aws_wafv2_regex_pattern_set",
			Name:     "Regex Pattern Set",
		},
		{
			Factory:  dataSourceRuleGroup,
			TypeName: "aws_wafv2_rule_group",
			Name:     "Rule Group",
		},
		{
			Factory:  dataSourceWebACL,
			TypeName: "aws_wafv2_web_acl",
			Name:     "Web ACL",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceIPSet,
			TypeName: "aws_wafv2_ip_set",
			Name:     "IP Set",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceRegexPatternSet,
			TypeName: "aws_wafv2_regex_pattern_set",
			Name:     "Regex Pattern Set",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceRuleGroup,
			TypeName: "aws_wafv2_rule_group",
			Name:     "Rule Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceWebACL,
			TypeName: "aws_wafv2_web_acl",
			Name:     "Web ACL",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceWebACLAssociation,
			TypeName: "aws_wafv2_web_acl_association",
			Name:     "Web ACL Association",
		},
		{
			Factory:  resourceWebACLLoggingConfiguration,
			TypeName: "aws_wafv2_web_acl_logging_configuration",
			Name:     "Web ACL Logging Configuration",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.WAFV2
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*wafv2_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return wafv2_sdkv2.NewFromConfig(cfg, func(o *wafv2_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
