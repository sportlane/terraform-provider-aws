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

//resource "aws_qbusiness_app" "test" {
//  display_name         = %[1]q
//  iam_service_role_arn = aws_iam_role.test.arn
//}

resource "aws_qbusiness_user" "test" {
  //application_id = aws_qbusiness_app.test.application_id
  application_id = "9e10b0a5-ae16-4109-bfff-ef81894c498d"
  user_id        = "user@example.com"

  user_aliases {
    alias {
      user_id       = "alias1"
	  index_id      = "6ed5dbab-6a79-407c-8d9a-3ac55ae3d67d"
	  datasource_id = "e7eadcde-3c98-4929-a907-d16721fc756b"
    }
	alias {
      user_id       = "alias2"
      index_id      = "6ed5dbab-6a79-407c-8d9a-3ac55ae3d67d"
      datasource_id = "0693542b-0b1c-4188-bbf8-a1e8a070dd4f"
    }
  }
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

//resource "aws_qbusiness_app" "test" {
//  display_name         = %[1]q
//  iam_service_role_arn = aws_iam_role.test.arn
//}

resource "aws_qbusiness_user" "test" {
  //application_id = aws_qbusiness_app.test.application_id
  application_id = "9e10b0a5-ae16-4109-bfff-ef81894c498d"
  user_id        = "user@example.com"

  user_aliases {
    alias {
      user_id       = "alias2"
	  index_id      = "6ed5dbab-6a79-407c-8d9a-3ac55ae3d67d"
	  datasource_id = "e7eadcde-3c98-4929-a907-d16721fc756b"
    }
  }
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
