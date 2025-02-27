---
subcategory: "CloudFront"
layout: "aws"
page_title: "AWS: aws_cloudfront_key_group"
description: |-
  Provides a CloudFront key group.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_cloudfront_key_group

## Example Usage

The following example below creates a CloudFront key group.

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import Fn, Token, TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.cloudfront_key_group import CloudfrontKeyGroup
from imports.aws.cloudfront_public_key import CloudfrontPublicKey
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        example = CloudfrontPublicKey(self, "example",
            comment="example public key",
            encoded_key=Token.as_string(Fn.file("public_key.pem")),
            name="example-key"
        )
        aws_cloudfront_key_group_example = CloudfrontKeyGroup(self, "example_1",
            comment="example key group",
            items=[example.id],
            name="example-key-group"
        )
        # This allows the Terraform resource name to match the original name. You can remove the call if you don't need them to match.
        aws_cloudfront_key_group_example.override_logical_id("example")
```

## Argument Reference

This resource supports the following arguments:

* `comment` - (Optional) A comment to describe the key group..
* `items` - (Required) A list of the identifiers of the public keys in the key group.
* `name` - (Required) A name to identify the key group.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `etag` - The identifier for this version of the key group.
* `id` - The identifier for the key group.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import CloudFront Key Group using the `id`. For example:

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.cloudfront_key_group import CloudfrontKeyGroup
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        CloudfrontKeyGroup.generate_config_for_import(self, "example", "4b4f2r1c-315d-5c2e-f093-216t50jed10f")
```

Using `terraform import`, import CloudFront Key Group using the `id`. For example:

```console
% terraform import aws_cloudfront_key_group.example 4b4f2r1c-315d-5c2e-f093-216t50jed10f
```

<!-- cache-key: cdktf-0.20.8 input-0db0e51d4954b3fa9a11b3c80e65d4eec2e3dd9da0abaf0109fbb3ae52ac9582 -->