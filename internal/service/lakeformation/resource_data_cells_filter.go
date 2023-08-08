// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lakeformation

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lakeformation"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	ResNameDataCellsFilter = "Resource Data Cells Filter"
)

// @SDKResource("aws_lakeformation_resource_lf_tags")
func ResourceDataCellsFilter() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceResourceDataCellsFilterCreate,
		ReadWithoutTimeout:   resourceResourceDataCellsFilterRead,
		DeleteWithoutTimeout: resourceResourceDataCellsFilterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"catalog_id": {
				Type:         schema.TypeString,
				Computed:     true,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"data_cells_filter": {
				Type:     schema.TypeList,
				ForceNew: true,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"catalog_id": {
							Type:         schema.TypeString,
							Computed:     true,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: verify.ValidAccountID,
						},
						"database_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"filter_expression": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"excluded_column_names": {
							Type:     schema.TypeSet,
							ForceNew: true,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
						},
					},
				},
			},
		},
	}
}

func resourceResourceDataCellsFilterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).LakeFormationConn(ctx)

	input := &lakeformation.CreateDataCellsFilterInput{}

	if v, ok := d.GetOk("data_cells_filter"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.TableData = expandDataCellsFilter(v.([]interface{})[0].(map[string]interface{}))
	}

	err := retry.RetryContext(ctx, IAMPropagationTimeout, func() *retry.RetryError {
		var err error
		_, err = conn.CreateDataCellsFilterWithContext(ctx, input)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, lakeformation.ErrCodeConcurrentModificationException) {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "AccessDeniedException", "is not authorized") {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.CreateDataCellsFilterWithContext(ctx, input)
	}

	if err != nil {
		return create.AddError(diags, names.LakeFormation, create.ErrActionCreating, ResNameDataCellsFilter, input.String(), err)
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(fmt.Sprintf("%d", create.StringHashcode(input.String())))

	return append(diags, resourceResourceLFTagsRead(ctx, d, meta)...)
}

func resourceResourceDataCellsFilterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).LakeFormationConn(ctx)

	var apiObject *lakeformation.DataCellsFilter

	if v, ok := d.GetOk("data_cells_filter"); ok {
		apiObject = expandDataCellsFilter(v.([]interface{})[0].(map[string]interface{}))
	}

	//apiObject := expandDataCellsFilter(d)
	input := &lakeformation.GetDataCellsFilterInput{
		Name:           apiObject.Name,
		DatabaseName:   apiObject.DatabaseName,
		TableName:      apiObject.TableName,
		TableCatalogId: apiObject.TableCatalogId,
	}

	_, err := conn.GetDataCellsFilterWithContext(ctx, input)

	if err != nil {
		return create.AddError(diags, names.LakeFormation, create.ErrActionReading, ResNameDataCellsFilter, d.Id(), err)
	}

	return diags
}

func resourceResourceDataCellsFilterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).LakeFormationConn(ctx)

	var apiObject *lakeformation.DataCellsFilter

	if v, ok := d.GetOk("data_cells_filter"); ok {
		apiObject = expandDataCellsFilter(v.([]interface{})[0].(map[string]interface{}))
	}

	input := &lakeformation.DeleteDataCellsFilterInput{
		Name:           apiObject.Name,
		DatabaseName:   apiObject.DatabaseName,
		TableName:      apiObject.TableName,
		TableCatalogId: apiObject.TableCatalogId,
	}

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		var err error
		_, err = conn.DeleteDataCellsFilterWithContext(ctx, input)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, lakeformation.ErrCodeConcurrentModificationException) {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "AccessDeniedException", "is not authorized") {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(fmt.Errorf("removing Lake Formation Data Cells Filter: %w", err))
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteDataCellsFilterWithContext(ctx, input)
	}

	if err != nil {
		return create.AddError(diags, names.LakeFormation, create.ErrActionDeleting, ResNameDataCellsFilter, d.Id(), err)
	}

	return diags
}

func expandDataCellsFilter(tfMap map[string]interface{}) *lakeformation.DataCellsFilter {

	input := &lakeformation.DataCellsFilter{}

	if v, ok := tfMap["catalog_id"].(string); ok && v != "" {
		input.TableCatalogId = aws.String(v)
	}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		input.Name = aws.String(v)
	}

	if v, ok := tfMap["database_name"].(string); ok && v != "" {
		input.DatabaseName = aws.String(v)
	}

	if v, ok := tfMap["table_name"].(string); ok && v != "" {
		input.TableName = aws.String(v)
	}

	if v, ok := tfMap["filter_expression"].(string); ok && v != "" {
		input.RowFilter = &lakeformation.RowFilter{
			FilterExpression: aws.String(v),
		}
	}

	if v, ok := tfMap["excluded_column_names"]; ok {
		if v, ok := v.(*schema.Set); ok && v.Len() > 0 {
			input.ColumnWildcard = &lakeformation.ColumnWildcard{
				ExcludedColumnNames: flex.ExpandStringSet(v),
			}
		}
	}

	return input
}
