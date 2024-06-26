// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rds

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_db_snapshot", name="DB Snapshot")
// @Tags(identifierAttribute="db_snapshot_arn")
// @Testing(tagsTest=false)
func ResourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSnapshotCreate,
		ReadWithoutTimeout:   resourceSnapshotRead,
		UpdateWithoutTimeout: resourceSnapshotUpdate,
		DeleteWithoutTimeout: resourceSnapshotDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			names.AttrAllocatedStorage: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			names.AttrAvailabilityZone: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_snapshot_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_snapshot_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrEncrypted: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			names.AttrEngine: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrEngineVersion: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrIOPS: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			names.AttrKMSKeyID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"license_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"option_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrPort: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shared_accounts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_db_snapshot_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrStorageType: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			names.AttrVPCID: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSConn(ctx)

	dbSnapshotID := d.Get("db_snapshot_identifier").(string)
	input := &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(d.Get("db_instance_identifier").(string)),
		DBSnapshotIdentifier: aws.String(dbSnapshotID),
		Tags:                 getTagsIn(ctx),
	}

	output, err := conn.CreateDBSnapshotWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating RDS DB Snapshot (%s): %s", dbSnapshotID, err)
	}

	d.SetId(aws.StringValue(output.DBSnapshot.DBSnapshotIdentifier))

	if err := waitDBSnapshotCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for RDS DB Snapshot (%s) create: %s", d.Id(), err)
	}

	if v, ok := d.GetOk("shared_accounts"); ok && v.(*schema.Set).Len() > 0 {
		input := &rds.ModifyDBSnapshotAttributeInput{
			AttributeName:        aws.String("restore"),
			DBSnapshotIdentifier: aws.String(d.Id()),
			ValuesToAdd:          flex.ExpandStringSet(v.(*schema.Set)),
		}

		_, err := conn.ModifyDBSnapshotAttributeWithContext(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying RDS DB Snapshot (%s) attribute: %s", d.Id(), err)
		}
	}

	return append(diags, resourceSnapshotRead(ctx, d, meta)...)
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSConn(ctx)

	snapshot, err := FindDBSnapshotByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] RDS DB Snapshot (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading RDS DB Snapshot (%s): %s", d.Id(), err)
	}

	arn := aws.StringValue(snapshot.DBSnapshotArn)
	d.Set(names.AttrAllocatedStorage, snapshot.AllocatedStorage)
	d.Set(names.AttrAvailabilityZone, snapshot.AvailabilityZone)
	d.Set("db_instance_identifier", snapshot.DBInstanceIdentifier)
	d.Set("db_snapshot_arn", arn)
	d.Set("db_snapshot_identifier", snapshot.DBSnapshotIdentifier)
	d.Set(names.AttrEncrypted, snapshot.Encrypted)
	d.Set(names.AttrEngine, snapshot.Engine)
	d.Set(names.AttrEngineVersion, snapshot.EngineVersion)
	d.Set(names.AttrIOPS, snapshot.Iops)
	d.Set(names.AttrKMSKeyID, snapshot.KmsKeyId)
	d.Set("license_model", snapshot.LicenseModel)
	d.Set("option_group_name", snapshot.OptionGroupName)
	d.Set(names.AttrPort, snapshot.Port)
	d.Set("source_db_snapshot_identifier", snapshot.SourceDBSnapshotIdentifier)
	d.Set("source_region", snapshot.SourceRegion)
	d.Set("snapshot_type", snapshot.SnapshotType)
	d.Set(names.AttrStatus, snapshot.Status)
	d.Set(names.AttrVPCID, snapshot.VpcId)

	input := &rds.DescribeDBSnapshotAttributesInput{
		DBSnapshotIdentifier: aws.String(d.Id()),
	}

	output, err := conn.DescribeDBSnapshotAttributesWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading RDS DB Snapshot (%s) attributes: %s", d.Id(), err)
	}

	d.Set("shared_accounts", flex.FlattenStringSet(output.DBSnapshotAttributesResult.DBSnapshotAttributes[0].AttributeValues))

	setTagsOut(ctx, snapshot.TagList)

	return diags
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSConn(ctx)

	if d.HasChange("shared_accounts") {
		o, n := d.GetChange("shared_accounts")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		additionList := ns.Difference(os)
		removalList := os.Difference(ns)

		input := &rds.ModifyDBSnapshotAttributeInput{
			AttributeName:        aws.String("restore"),
			DBSnapshotIdentifier: aws.String(d.Id()),
			ValuesToAdd:          flex.ExpandStringSet(additionList),
			ValuesToRemove:       flex.ExpandStringSet(removalList),
		}

		_, err := conn.ModifyDBSnapshotAttributeWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying RDS DB Snapshot (%s) attributes: %s", d.Id(), err)
		}
	}

	return diags
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSConn(ctx)

	log.Printf("[DEBUG] Deleting RDS DB Snapshot: %s", d.Id())
	_, err := conn.DeleteDBSnapshotWithContext(ctx, &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBSnapshotNotFoundFault) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting RDS DB Snapshot (%s): %s", d.Id(), err)
	}

	return diags
}

func FindDBSnapshotByID(ctx context.Context, conn *rds.RDS, id string) (*rds.DBSnapshot, error) {
	input := &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: aws.String(id),
	}
	output, err := findDBSnapshot(ctx, conn, input, tfslices.PredicateTrue[*rds.DBSnapshot]())

	if err != nil {
		return nil, err
	}

	// Eventual consistency check.
	if aws.StringValue(output.DBSnapshotIdentifier) != id {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return output, nil
}

func findDBSnapshot(ctx context.Context, conn *rds.RDS, input *rds.DescribeDBSnapshotsInput, filter tfslices.Predicate[*rds.DBSnapshot]) (*rds.DBSnapshot, error) {
	output, err := findDBSnapshots(ctx, conn, input, filter)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSinglePtrResult(output)
}

func findDBSnapshots(ctx context.Context, conn *rds.RDS, input *rds.DescribeDBSnapshotsInput, filter tfslices.Predicate[*rds.DBSnapshot]) ([]*rds.DBSnapshot, error) {
	var output []*rds.DBSnapshot

	err := conn.DescribeDBSnapshotsPagesWithContext(ctx, input, func(page *rds.DescribeDBSnapshotsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.DBSnapshots {
			if v != nil && filter(v) {
				output = append(output, v)
			}
		}

		return !lastPage
	})

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBSnapshotNotFoundFault) {
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
