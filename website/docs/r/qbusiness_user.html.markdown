---
subcategory: "Amazon Q Business"
layout: "aws"
page_title: "AWS: aws_qbusiness_user"
description: |-
  Provides a Q Business User resource.
---

# Resource: aws_qbusiness_user

Provides a Q Business User resource.

## Example Usage

```terraform
resource "aws_qbusiness_user" "example" {
  application_id = aws_qbusiness_app.test.application_id
  user_id        = "user@exmaple.com"

  user_aliases {
    alias {
      user_id       = "user@example.com"
      datasource_id = aws_qbusiness_datasource.test.datasource_id
      index_id      = aws_qbusiness_index.test.index_id
    }
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `application_id` - (Required) Id of the Q Business application.
* `user_id` - (Required) User email attached to a user mapping.
* `user_aliases` - (Optional) List of user aliases attached to a user mapping.

`user_aliases` supports the following:

* `user_id` - (Required) Identifier of the user id associated with the user aliases.
* `datasource_id` - (Required) Identifier of the data source that the user aliases are associated with.
* `index_id` - (Required) Identifier of the index that the user aliases are associated with.
