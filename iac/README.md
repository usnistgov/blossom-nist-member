# BloSS@M NIST Member IaC

## Member Configuration

This Terraform configuration uses [Terraform Workspaces](https://developer.hashicorp.com/terraform/language/state/workspaces) to segment member-specific configuration & state.

Each workspace should have a corresponding JSON configuration file in the [`configurations/`](./configurations/) directory.
As an example, the Terraform workspace `BLOSSON_NIST2` would expect the corresponding configuration file [`configurations/BLOSSON_NIST2.json`](./configurations/BLOSSON_NIST2.json).

Below are the keys expected in the configuration:

- `network_id`: The ID of the Hyperledger Fabric network managed by AMB.
- `network_name`: The name of the Hyperledger Fabric network.
- `member_id`: The ID of the Hyperledger Fabric member managed by AMB.
- `member_name`: The name of the Hyperledger Fabric member.
- `peer_node_id`: The ID of the Hyperledger Fabric peer node managed by AMB. 
- `channel_name`: The name of the Hyperledger Fabric channel that the BLOSS@M chaincode is deployed onto.
- `contract_name`: The name given to the BLOSS@M chaincode deployed on the channel.
- `identities_ssm_prefix`: The prefix the chaincode transaction lambda will use to extract user identities (e.x. `/identities`).

    *NOTE: the IAM role defined by `lambda_execution_iam_role_name` MUST be configured to allow access to SSM parameters within this prefix.*

- `cognito_user_pool_name`: The name of the Cognito User Pool to use for user identities.
- `apigw_s3_integration_iam_role_name`: The name of the IAM role attached to the API Gateway S3 integration.
- `lambda_execution_iam_role_name`:  The name of the IAM role attached to the Lambda.

### IAM Roles

Currently IAM roles must be managed externally to this configuration due to permission limits imposed onto our accounts.

#### API-GW S3 Integration Role

This IAM role must allow API Gateway to read and list S3 objects.
In order to limit the scope of the role, we suggest limiting the role's access to objects that have been tagged with `Purpose=blossom-frontend`.

Example role:

```json
{
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
                "Service": "apigateway.amazonaws.com"
            },
            "Sid": ""
        },
    ],
    "Version": "2012-10-17"
}
```

Policy attachment:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:Get*",
                "s3:List*",
                "s3-object-lambda:Get*",
                "s3-object-lambda:List*"
            ],
            "Resource": "*",
            "Condition": {
                "StringEquals": {
                    "s3:ExistingObjectTag/Purpose": "blossom-frontend"
                }
            }
        }
    ]
}
```

#### Lambda Execution Role

This IAM role must allow the identity Lambda to access AmazonManagedBlockchain, Get SSM parameters, and write to Cloudwatch logs.

*NOTE: this is a snapshot of the policy at its current state, please restrict this role before deployment (TODO)*

Example role:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```

Policy Attachments:

- AWS Built-in `AWSLambdaVPCAccessExecutionRole`
- *NOTE: TODO heavily restrict this role*

    ```json
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Action": [
                    "managedblockchain:*"
                ],
                "Resource": "*",
                "Effect": "Allow",
                "Sid": "ManagedBlockchainAllowAll"
            },
            {
                "Action": [
                    "s3:ListBucket",
                    "s3:GetObject*"
                ],
                "Resource": "*",
                "Effect": "Allow",
                "Sid": "S3ReadAll"
            },
            {
                "Action": [
                    "ssm:DescribeParameters"
                ],
                "Resource": "*",
                "Effect": "Allow",
                "Sid": "SsmParametersDescribeAll"
            },
            {
                "Action": [
                    "ssm:GetParameter",
                    "ssm:GetParameters",
                    "ssm:GetParametersByPath"
                ],
                "Resource": "arn:aws:ssm:<region>:<account id>:parameter/nist/blossom/*", // this should correspond to the configuration parameter `identities_ssm_prefix`
                "Effect": "Allow",
                "Sid": "SsmParametersGetBlossomOnly"
            },
            {
                "Action": [
                    "logs:CreateLogGroup",
                    "logs:CreateLogStream",
                    "logs:PutLogEvents"
                ],
                "Resource": "*",
                "Effect": "Allow",
                "Sid": "CloudwatchLogWrite"
            }
        ]
    }
    ```

## Usage

To provision or update the infrastructure, perform the following steps:

1. If you cloned the repository without initializing submodules, do so now:

    ```sh
    $ git submodule update --init --recursive
    ```

1. Select the correct workspace:

    Before provisioning the infrastructure for a member, select the member in Terraform.

    You can view the available workspaces using the command:

    ```sh
    $ ./ter.sh workspace list
    ```

    You can select the appropriate workspace using the command:

    ```sh
    $ ./ter.sh workspace select <workspace>
    ```

1. Provision the infrastructure:

    Provisioning in Terraform is wrapped by a Makefile.
    This allows build operations for the related projects (lambda and dashboard) to happen before provisioning.

    In order to provision the infrastructure, run:

    ```sh
    $ make apply
    ```

### Makefile targets

This project provides several utility Makefile targets.
For more details, run `make help`.

### Running other Terraform commands

Any terraform command should be run using the [`./ter.sh`](./ter.sh) wrapper.
This ensures that your credentials are valid before running a command.

As an example, if you want to run `terraform refresh`, instead run `./ter.sh refresh`.
