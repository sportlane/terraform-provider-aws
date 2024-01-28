// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package qbusiness_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/qbusiness"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfqbusiness "github.com/hashicorp/terraform-provider-aws/internal/service/qbusiness"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccQBusinessUser_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var user qbusiness.GetUserOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_qbusiness_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckUser(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, "qbusiness"),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccQBusinessUser_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var user qbusiness.GetUserOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_qbusiness_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckUser(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, "qbusiness"),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfqbusiness.ResourceUser(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccQBusinessUser_alias(t *testing.T) {
	ctx := acctest.Context(t)
	var user qbusiness.GetUserOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_qbusiness_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckUser(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, "qbusiness"),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_alias(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_aliases.0.alias.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUserConfig_aliasRemove(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_aliases.0.alias.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_aliases.0.alias.0.user_id", "alias2"),
				),
			},
		},
	})
}

func testAccPreCheckUser(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).QBusinessClient(ctx)

	input := &qbusiness.ListApplicationsInput{}

	_, err := conn.ListApplications(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccCheckUserExists(ctx context.Context, n string, v *qbusiness.GetUserOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).QBusinessClient(ctx)

		output, err := tfqbusiness.FindUserByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckUserDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).QBusinessClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_qbusiness_user" {
				continue
			}

			_, err := tfqbusiness.FindUserByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Amazon Q User %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccUserConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_qbusiness_app" "test" {
  display_name         = %[1]q
  iam_service_role_arn = aws_iam_role.test.arn
}

resource "aws_qbusiness_user" "test" {
  application_id = aws_qbusiness_app.test.application_id
  user_id        = "user@example.com"
}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
"Version": "2012-10-17",
"Statement": [
    {
    "Action": "sts:AssumeRole",
    "Principal": {
        "Service": "qbusiness.${data.aws_partition.current.dns_suffix}"
    },
    "Effect": "Allow",
    "Sid": ""
    }
  ]
}
EOF
}
`, rName)
}

func testAccUserConfig_alias(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_qbusiness_app" "test" {
  display_name         = %[1]q
  iam_service_role_arn = aws_iam_role.test.arn
}

resource "aws_qbusiness_user" "test" {
  application_id = aws_qbusiness_app.test.application_id
  user_id        = "user@example.com"

  user_aliases {
    alias {
      user_id       = "alias1"
      index_id      = aws_qbusiness_index.test.index_id
      datasource_id = aws_qbusiness_datasource.test.datasource_id
    }
    alias {
      user_id       = "alias2"
      index_id      = aws_qbusiness_index.test.index_id
      datasource_id = aws_qbusiness_datasource.test1.datasource_id
    }
  }
}

resource "aws_qbusiness_index" "test" {
  application_id       = aws_qbusiness_app.test.application_id
  display_name         = %[1]q
  capacity_configuration {
    units = 1
  }
  description          = %[1]q
}

resource "aws_s3_bucket" "test" {
  bucket = "%[1]q-1"
  force_destroy = true
}

resource "aws_s3_bucket" "test1" {
  bucket = "%[1]q-2"
  force_destroy = true
}

resource "aws_qbusiness_datasource" "test" {
  application_id       = aws_qbusiness_app.test.application_id
  index_id             = aws_qbusiness_index.test.index_id
  display_name         = %[1]q
  iam_service_role_arn = aws_iam_role.test.arn
  configuration        = jsonencode({
    type                     = "S3"
    connectionConfiguration  = {
      repositoryEndpointMetadata = {
        BucketName = aws_s3_bucket.test.bucket
      }
    }
    syncMode                 = "FULL_CRAWL"
      repositoryConfigurations = {
        document = {
          fieldMappings = []
        }
      }
  })
}

resource "aws_qbusiness_datasource" "test1" {
  application_id       = aws_qbusiness_app.test.application_id
  index_id             = aws_qbusiness_index.test.index_id
  display_name         = %[1]q
  iam_service_role_arn = aws_iam_role.test.arn
  configuration        = jsonencode({
    type                     = "S3"
    connectionConfiguration  = {
      repositoryEndpointMetadata = {
        BucketName = aws_s3_bucket.test1.bucket
      }
    }
    syncMode                 = "FULL_CRAWL"
      repositoryConfigurations = {
        document = {
          fieldMappings = []
        }
      }
  })
}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
"Version": "2012-10-17",
"Statement": [
    {
    "Action": "sts:AssumeRole",
    "Principal": {
        "Service": "qbusiness.${data.aws_partition.current.dns_suffix}"
    },
    "Effect": "Allow",
    "Sid": ""
    }
  ]
}
EOF
}
`, rName)
}

func testAccUserConfig_aliasRemove(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_qbusiness_app" "test" {
  display_name         = %[1]q
  iam_service_role_arn = aws_iam_role.test.arn
}

resource "aws_qbusiness_user" "test" {
  application_id = aws_qbusiness_app.test.application_id
  user_id        = "user@example.com"

  user_aliases {
    alias {
      user_id       = "alias2"
      index_id      = aws_qbusiness_index.test.index_id
      datasource_id = aws_data_qbusinessource.test.datasource_id
    }
  }
}

resource "aws_qbusiness_index" "test" {
  application_id       = aws_qbusiness_app.test.application_id
  display_name         = %[1]q
  capacity_configuration {
    units = 1
  }
  description          = %[1]q
}

resource "aws_s3_bucket" "test" {
  bucket = "%[1]q-1"
  force_destroy = true
}

resource "aws_s3_bucket" "test1" {
  bucket = "%[1]q-2"
  force_destroy = true
}

resource "aws_qbusiness_datasource" "test" {
  application_id       = aws_qbusiness_app.test.application_id
  index_id             = aws_qbusiness_index.test.index_id
  display_name         = %[1]q
  iam_service_role_arn = aws_iam_role.test.arn
  configuration        = jsonencode({
    type                     = "S3"
    connectionConfiguration  = {
      repositoryEndpointMetadata = {
        BucketName = aws_s3_bucket.test.bucket
      }
    }
    syncMode                 = "FULL_CRAWL"
      repositoryConfigurations = {
        document = {
          fieldMappings = []
        }
      }
  })
}

resource "aws_qbusiness_datasource" "test1" {
  application_id       = aws_qbusiness_app.test.application_id
  index_id             = aws_qbusiness_index.test.index_id
  display_name         = %[1]q
  iam_service_role_arn = aws_iam_role.test.arn
  configuration        = jsonencode({
    type                     = "S3"
    connectionConfiguration  = {
      repositoryEndpointMetadata = {
        BucketName = aws_s3_bucket.test1.bucket
      }
    }
    syncMode                 = "FULL_CRAWL"
      repositoryConfigurations = {
        document = {
          fieldMappings = []
        }
      }
  })
}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
"Version": "2012-10-17",
"Statement": [
    {
    "Action": "sts:AssumeRole",
    "Principal": {
        "Service": "qbusiness.${data.aws_partition.current.dns_suffix}"
    },
    "Effect": "Allow",
    "Sid": ""
    }
  ]
}
EOF
}
`, rName)
}
