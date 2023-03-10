# BloSS@M NIST Member IaC

## Usage

To provision or update the infrastructure, perform the following steps:

1. If you cloned the repository without initializing submodules, do so now:

    ```
    $ git submodule update --init --recursive
    ```

1. Select the correct workspace:

    This project uses [Terraform Workspaces](https://developer.hashicorp.com/terraform/language/state/workspaces) to segment member-specific configuration & state.
    Before provisioning the infrastructure for a member, select the member in Terraform.

    You can view the available workspaces using the command:

    ```
    $ ./ter.sh workspace list
    ```

    You can select the appropriate workspace using the command:

    ```
    $ ./ter.sh workspace select <workspace>
    ```

1. Provision the infrastructure:

    Provisioning in Terraform is wrapped by a Makefile.
    This allows build operations for the related projects (lambda and dashboard) to happen before provisioning.

    In order to provision the infrastructure, run:

    ```
    $ make apply
    ```

### Running other Terraform commands

Any terraform command should be run using the [`./ter.sh`](./ter.sh) wrapper.
This ensures that your credentials are valid before running a command.

As an example, if you want to run `terraform refresh`, instead run `./ter.sh refresh`.
