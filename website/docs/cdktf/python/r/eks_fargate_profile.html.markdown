---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "AWS: aws_eks_fargate_profile"
description: |-
  Manages an EKS Fargate Profile
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_eks_fargate_profile

Manages an EKS Fargate Profile.

## Example Usage

```python
# Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import Token, property_access, TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.eks_fargate_profile import EksFargateProfile
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        EksFargateProfile(self, "example",
            cluster_name=Token.as_string(aws_eks_cluster_example.name),
            fargate_profile_name="example",
            pod_execution_role_arn=Token.as_string(aws_iam_role_example.arn),
            selector=[EksFargateProfileSelector(
                namespace="example"
            )
            ],
            subnet_ids=Token.as_list(property_access(aws_subnet_example, ["*", "id"]))
        )
```

### Example IAM Role for EKS Fargate Profile

```python
# Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import Fn, Token, TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.iam_role import IamRole
from imports.aws.iam_role_policy_attachment import IamRolePolicyAttachment
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        example = IamRole(self, "example",
            assume_role_policy=Token.as_string(
                Fn.jsonencode({
                    "Statement": [{
                        "Action": "sts:AssumeRole",
                        "Effect": "Allow",
                        "Principal": {
                            "Service": "eks-fargate-pods.amazonaws.com"
                        }
                    }
                    ],
                    "Version": "2012-10-17"
                })),
            name="eks-fargate-profile-example"
        )
        IamRolePolicyAttachment(self, "example-AmazonEKSFargatePodExecutionRolePolicy",
            policy_arn="arn:aws:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy",
            role=example.name
        )
```

## Argument Reference

The following arguments are required:

* `cluster_name` – (Required) Name of the EKS Cluster. Must be between 1-100 characters in length. Must begin with an alphanumeric character, and must only contain alphanumeric characters, dashes and underscores (`^[0-9A-Za-z][A-Za-z0-9\-_]+$`).
* `fargate_profile_name` – (Required) Name of the EKS Fargate Profile.
* `pod_execution_role_arn` – (Required) Amazon Resource Name (ARN) of the IAM Role that provides permissions for the EKS Fargate Profile.
* `selector` - (Required) Configuration block(s) for selecting Kubernetes Pods to execute with this EKS Fargate Profile. Detailed below.
* `subnet_ids` – (Required) Identifiers of private EC2 Subnets to associate with the EKS Fargate Profile. These subnets must have the following resource tag: `kubernetes.io/cluster/CLUSTER_NAME` (where `CLUSTER_NAME` is replaced with the name of the EKS Cluster).

The following arguments are optional:

* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

### selector Configuration Block

The following arguments are required:

* `namespace` - (Required) Kubernetes namespace for selection.

The following arguments are optional:

* `labels` - (Optional) Key-value map of Kubernetes labels for selection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the EKS Fargate Profile.
* `id` - EKS Cluster name and EKS Fargate Profile name separated by a colon (`:`).
* `status` - Status of the EKS Fargate Profile.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `10m`)
* `delete` - (Default `10m`)

## Import

EKS Fargate Profiles can be imported using the `cluster_name` and `fargate_profile_name` separated by a colon (`:`), e.g.,

```
$ terraform import aws_eks_fargate_profile.my_fargate_profile my_cluster:my_fargate_profile
```

<!-- cache-key: cdktf-0.17.1 input-d30b94566032f99ee0124f08aafaebec6c68a111159669a79f2ec46b163570e5 -->