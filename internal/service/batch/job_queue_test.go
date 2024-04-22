// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfbatch "github.com/hashicorp/terraform-provider-aws/internal/service/batch"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// ComputeEnvironments has been deprecated. The Import step of tests that use ComputeEnvironments
// need to ignore
var ignoreDeprecatedCEOForImports = []string{
	"compute_environment_order",
	"compute_environment_order.#",
	"compute_environment_order.0.%",
	"compute_environment_order.0.compute_environment",
	"compute_environment_order.0.order",
	"compute_environments.#",
	"compute_environments.0",
}

func TestAccBatchJobQueue_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_state(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "batch", fmt.Sprintf("job-queue/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "compute_environments.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environments.0", "aws_batch_compute_environment.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "state", string(awstypes.JQStateEnabled)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccBatchJobQueue_basicCEO(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_stateCEO(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "batch", fmt.Sprintf("job-queue/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "compute_environment_order.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environment_order.0.compute_environment", "aws_batch_compute_environment.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "state", string(awstypes.JQStateEnabled)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreDeprecatedCEOForImports,
			},
		},
	})
}

func TestAccBatchJobQueue_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLaunchTemplateDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_state(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfbatch.ResourceJobQueue, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBatchJobQueue_disappearsCEO(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLaunchTemplateDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_stateCEO(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfbatch.ResourceJobQueue, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBatchJobQueue_MigrateFromPluginSDK(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:   acctest.ErrorCheck(t, names.BatchServiceID),
		CheckDestroy: testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						Source:            "hashicorp/aws",
						VersionConstraint: "5.13.1",
					},
				},
				Config: testAccJobQueueConfig_state(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "batch", fmt.Sprintf("job-queue/%s", rName)),
				),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
				Config:                   testAccJobQueueConfig_state(rName, string(awstypes.JQStateEnabled)),
				PlanOnly:                 true,
			},
		},
	})
}

func TestAccBatchJobQueue_ComputeEnvironments_multiple(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_ComputeEnvironments_multiple(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					resource.TestCheckResourceAttr(resourceName, "compute_environments.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environments.0", "aws_batch_compute_environment.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environments.1", "aws_batch_compute_environment.more.0", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environments.2", "aws_batch_compute_environment.more.1", "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJobQueueConfig_ComputeEnvironments_multipleReorder(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					resource.TestCheckResourceAttr(resourceName, "compute_environments.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environments.0", "aws_batch_compute_environment.more.0", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environments.1", "aws_batch_compute_environment.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environments.2", "aws_batch_compute_environment.more.1", "arn"),
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

func TestAccBatchJobQueue_ComputeEnvironmentOrder_multiple(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_ComputeEnvironmentOrder_multiple(rName, string(awstypes.JQStateEnabled), 1, 2, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					resource.TestCheckResourceAttr(resourceName, "compute_environment_order.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environment_order.0.compute_environment", "aws_batch_compute_environment.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environment_order.1.compute_environment", "aws_batch_compute_environment.more.0", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environment_order.2.compute_environment", "aws_batch_compute_environment.more.1", "arn"),
				),
			},
			{
				Config: testAccJobQueueConfig_ComputeEnvironmentOrder_multiple(rName, string(awstypes.JQStateEnabled), 2, 1, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					resource.TestCheckResourceAttr(resourceName, "compute_environment_order.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "compute_environment_order.0.order", "2"),
					resource.TestCheckResourceAttr(resourceName, "compute_environment_order.1.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "compute_environment_order.2.order", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environment_order.0.compute_environment", "aws_batch_compute_environment.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environment_order.1.compute_environment", "aws_batch_compute_environment.more.0", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "compute_environment_order.2.compute_environment", "aws_batch_compute_environment.more.1", "arn"),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/8083
func TestAccBatchJobQueue_ComputeEnvironments_externalOrderUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_state(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					testAccCheckJobQueueComputeEnvironmentOrderUpdate(ctx, &jobQueue1),
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

func TestAccBatchJobQueue_priority(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1, jobQueue2 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_priority(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
				),
			},
			{
				Config: testAccJobQueueConfig_priority(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue2),
					resource.TestCheckResourceAttr(resourceName, "priority", "2"),
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

func TestAccBatchJobQueue_schedulingPolicy(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1, jobQueue2 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	schedulingPolicyName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	schedulingPolicyName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				// last variable selects the scheduling policy's arn. In this case, the first scheduling policy's arn.
				Config: testAccJobQueueConfig_schedulingPolicy(rName, schedulingPolicyName1, schedulingPolicyName2, "first"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					resource.TestCheckResourceAttrSet(resourceName, "scheduling_policy_arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// test switching the scheduling_policy_arn by changing the last variable to select the second scheduling policy's arn.
				Config: testAccJobQueueConfig_schedulingPolicy(rName, schedulingPolicyName1, schedulingPolicyName2, "second"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue2),
					resource.TestCheckResourceAttrSet(resourceName, "scheduling_policy_arn"),
				),
			},
		},
	})
}

func TestAccBatchJobQueue_state(t *testing.T) {
	ctx := acctest.Context(t)
	var jobQueue1, jobQueue2 awstypes.JobQueueDetail
	resourceName := "aws_batch_job_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckJobQueueDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccJobQueueConfig_state(rName, string(awstypes.JQStateDisabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue1),
					resource.TestCheckResourceAttr(resourceName, "state", string(awstypes.JQStateDisabled)),
				),
			},
			{
				Config: testAccJobQueueConfig_state(rName, string(awstypes.JQStateEnabled)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobQueueExists(ctx, resourceName, &jobQueue2),
					resource.TestCheckResourceAttr(resourceName, "state", string(awstypes.JQStateEnabled)),
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

func testAccCheckJobQueueExists(ctx context.Context, n string, jq *awstypes.JobQueueDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		log.Printf("State: %#v", s.RootModule().Resources)
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Batch Job Queue ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).BatchClient(ctx)
		name := rs.Primary.Attributes["name"]
		queue, err := tfbatch.FindJobQueueByName(ctx, conn, name)
		if err != nil {
			return err
		}
		if queue == nil {
			return fmt.Errorf("Not found: %s", n)
		}
		*jq = *queue

		return nil
	}
}

func testAccCheckJobQueueDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_batch_job_queue" {
				continue
			}
			conn := acctest.Provider.Meta().(*conns.AWSClient).BatchClient(ctx)
			jq, err := tfbatch.FindJobQueueByName(ctx, conn, rs.Primary.Attributes["name"])
			if err == nil {
				if jq != nil {
					return fmt.Errorf("Error: Job Queue still exists")
				}
			}
			return nil
		}
		return nil
	}
}

// testAccCheckJobQueueComputeEnvironmentOrderUpdate simulates the change of a Compute Environment Order
// An external update to the Batch Job Queue (e.g. console) may trigger changes to the value of the Order
// parameter that do not affect the operation of the queue itself, but the resource logic needs to handle.
// For example, Terraform may set a single Compute Environment with Order 0, but the console updates it to 1.
func testAccCheckJobQueueComputeEnvironmentOrderUpdate(ctx context.Context, jobQueue *awstypes.JobQueueDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).BatchClient(ctx)

		input := &batch.UpdateJobQueueInput{
			ComputeEnvironmentOrder: jobQueue.ComputeEnvironmentOrder,
			JobQueue:                jobQueue.JobQueueName,
		}
		name := aws.ToString(jobQueue.JobQueueName)

		if len(input.ComputeEnvironmentOrder) != 1 {
			return fmt.Errorf("expected one ComputeEnvironmentOrder in Batch Job Queue (%s)", name)
		}

		if aws.ToInt32(input.ComputeEnvironmentOrder[0].Order) != 0 {
			return fmt.Errorf("expected first ComputeEnvironmentOrder in Batch Job Queue (%s) to have existing Order value of 0", name)
		}

		input.ComputeEnvironmentOrder[0].Order = aws.Int32(1)

		_, err := conn.UpdateJobQueue(ctx, input)

		if err != nil {
			return fmt.Errorf("error updating Batch Job Queue (%s): %s", name, err)
		}

		return nil
	}
}

func testAccJobQueueConfigBase(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name               = %[1]q
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
    {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
        "Service": "batch.${data.aws_partition.current.dns_suffix}"
        }
    }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "test" {
  role       = aws_iam_role.test.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSBatchServiceRole"
}

resource "aws_iam_role" "ecs_instance_role" {
  name = "%[1]s-ecs"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
    {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
        "Service": "ec2.${data.aws_partition.current.dns_suffix}"
        }
    }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs_instance_role" {
  role       = aws_iam_role.ecs_instance_role.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs_instance_role" {
  name = aws_iam_role.ecs_instance_role.name
  role = aws_iam_role_policy_attachment.ecs_instance_role.role
}

resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "terraform-testacc-batch-job-queue"
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id
}

resource "aws_subnet" "test" {
  cidr_block = "10.1.1.0/24"
  vpc_id     = aws_vpc.test.id

  tags = {
    Name = "tf-acc-batch-job-queue"
  }
}

resource "aws_batch_compute_environment" "test" {
  compute_environment_name = %[1]q
  service_role             = aws_iam_role.test.arn
  type                     = "MANAGED"

  compute_resources {
    instance_role      = aws_iam_instance_profile.ecs_instance_role.arn
    instance_type      = ["c5", "m5", "r5"]
    max_vcpus          = 1
    min_vcpus          = 0
    security_group_ids = [aws_security_group.test.id]
    subnets            = [aws_subnet.test.id]
    type               = "EC2"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName)
}

func testAccJobQueueConfig_priority(rName string, priority int) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = [aws_batch_compute_environment.test.arn]
  name                 = %[1]q
  priority             = %[2]d
  state                = "ENABLED"
}
`, rName, priority))
}

func testAccJobQueueSchedulingPolicy(rName string, rName2 string) string {
	return fmt.Sprintf(`
resource "aws_batch_scheduling_policy" "test1" {
  name = %[1]q

  fair_share_policy {
    compute_reservation = 1
    share_decay_seconds = 3600

    share_distribution {
      share_identifier = "A1*"
      weight_factor    = 0.1
    }
  }
}

resource "aws_batch_scheduling_policy" "test2" {
  name = %[2]q

  fair_share_policy {
    compute_reservation = 1
    share_decay_seconds = 3600

    share_distribution {
      share_identifier = "A2"
      weight_factor    = 0.2
    }
  }
}
`, rName, rName2)
}

func testAccJobQueueConfig_schedulingPolicy(rName string, schedulingPolicyName1 string, schedulingPolicyName2 string, selectSchedulingPolicy string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		testAccJobQueueSchedulingPolicy(schedulingPolicyName1, schedulingPolicyName2),
		fmt.Sprintf(`
locals {
  select_scheduling_policy = %[2]q
}

resource "aws_batch_job_queue" "test" {
  compute_environments  = [aws_batch_compute_environment.test.arn]
  name                  = %[1]q
  priority              = 1
  scheduling_policy_arn = local.select_scheduling_policy == "first" ? aws_batch_scheduling_policy.test1.arn : aws_batch_scheduling_policy.test2.arn
  state                 = "ENABLED"
}
`, rName, selectSchedulingPolicy))
}

func testAccJobQueueConfig_state(rName string, state string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = [aws_batch_compute_environment.test.arn]

  name     = %[1]q
  priority = 1
  state    = %[2]q
}
`, rName, state))
}

func testAccJobQueueConfig_stateCEO(rName string, state string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environment_order {
    compute_environment = aws_batch_compute_environment.test.arn
    order               = 1
  }

  name     = %[1]q
  priority = 1
  state    = %[2]q
}
`, rName, state))
}

func testAccJobQueueConfig_ComputeEnvironments_multiple(rName string, state string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = concat(
    [aws_batch_compute_environment.test.arn],
    aws_batch_compute_environment.more[*].arn,
  )
  name     = %[1]q
  priority = 1
  state    = %[2]q
}

resource "aws_batch_compute_environment" "more" {
  count = 2

  compute_environment_name = "%[1]s-${count.index + 1}"
  service_role             = aws_iam_role.test.arn
  type                     = "MANAGED"

  compute_resources {
    instance_role      = aws_iam_instance_profile.ecs_instance_role.arn
    instance_type      = ["c5", "m5", "r5"]
    max_vcpus          = 1
    min_vcpus          = 0
    security_group_ids = [aws_security_group.test.id]
    subnets            = [aws_subnet.test.id]
    type               = "EC2"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, state))
}

func testAccJobQueueConfig_ComputeEnvironmentOrder_multiple(rName string, state string, o1 int, o2 int, o3 int) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environment_order {
    order               = %[3]d
    compute_environment = aws_batch_compute_environment.test.arn
  }

  compute_environment_order {
    order               = %[4]d
    compute_environment = aws_batch_compute_environment.more[0].arn
  }

  compute_environment_order {
    order               = %[5]d
    compute_environment = aws_batch_compute_environment.more[1].arn
  }

  name     = %[1]q
  priority = 1
  state    = %[2]q
}

resource "aws_batch_compute_environment" "more" {
  count = 2

  compute_environment_name = "%[1]s-${count.index + 1}"
  service_role             = aws_iam_role.test.arn
  type                     = "MANAGED"

  compute_resources {
    instance_role      = aws_iam_instance_profile.ecs_instance_role.arn
    instance_type      = ["c5", "m5", "r5"]
    max_vcpus          = 1
    min_vcpus          = 0
    security_group_ids = [aws_security_group.test.id]
    subnets            = [aws_subnet.test.id]
    type               = "EC2"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, state, o1, o2, o3))
}

func testAccJobQueueConfig_ComputeEnvironments_multipleReorder(rName string, state string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = [
    aws_batch_compute_environment.more[0].arn,
    aws_batch_compute_environment.test.arn,
    aws_batch_compute_environment.more[1].arn,
  ]
  name     = %[1]q
  priority = 1
  state    = %[2]q
}

resource "aws_batch_compute_environment" "more" {
  count = 2

  compute_environment_name = "%[1]s-${count.index + 1}"
  service_role             = aws_iam_role.test.arn
  type                     = "MANAGED"

  compute_resources {
    instance_role      = aws_iam_instance_profile.ecs_instance_role.arn
    instance_type      = ["c5", "m5", "r5"]
    max_vcpus          = 1
    min_vcpus          = 0
    security_group_ids = [aws_security_group.test.id]
    subnets            = [aws_subnet.test.id]
    type               = "EC2"
  }

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, state))
}

func testAccJobQueueConfig_tags0(rName string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = [aws_batch_compute_environment.test.arn]
  name                 = %[1]q
  priority             = 1
  state                = "DISABLED"
}
`, rName))
}

func testAccJobQueueConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = [aws_batch_compute_environment.test.arn]
  name                 = %[1]q
  priority             = 1
  state                = "DISABLED"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccJobQueueConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = [aws_batch_compute_environment.test.arn]
  name                 = %[1]q
  priority             = 1
  state                = "DISABLED"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccJobQueueConfig_tagsNull(rName, tagKey1 string) string {
	return acctest.ConfigCompose(
		testAccJobQueueConfigBase(rName),
		fmt.Sprintf(`
resource "aws_batch_job_queue" "test" {
  compute_environments = [aws_batch_compute_environment.test.arn]
  name                 = %[1]q
  priority             = 1
  state                = "DISABLED"

  tags = {
    %[2]q = null
  }
}
`, rName, tagKey1))
}

func testAccCheckLaunchTemplateDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_launch_template" {
				continue
			}

			resp, err := conn.DescribeLaunchTemplatesWithContext(ctx, &ec2.DescribeLaunchTemplatesInput{
				LaunchTemplateIds: []*string{aws.String(rs.Primary.ID)},
			})

			if err == nil {
				if len(resp.LaunchTemplates) != 0 && *resp.LaunchTemplates[0].LaunchTemplateId == rs.Primary.ID {
					return fmt.Errorf("Launch Template still exists")
				}
			}

			if tfawserr.ErrCodeEquals(err, "InvalidLaunchTemplateId.NotFound") {
				log.Printf("[WARN] launch template (%s) not found.", rs.Primary.ID)
				continue
			}
			return err
		}

		return nil
	}
}
