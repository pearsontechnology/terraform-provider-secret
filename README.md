# Terraform Provider

## Requirements

  * [Terraform](https://www.terraform.io/downloads.html) 0.10.x
  * Go 1.8 (to build the provider plugin)

## Building The Provider

Clone repository to: $GOPATH/src/github.com/pearsontechnology/terraform-provider-secret

$ mkdir -p $GOPATH/src/github.com/pearsontechnology; cd $GOPATH/src/github.com/pearsontechnology
$ git clone git@github.com:pearsontechnology/terraform-provider-secret

Enter the provider directory and build the provider

$ cd $GOPATH/src/github.com/hashicorp/terraform-provider-secret
$ hack/build/build.sh

## Using the provider

### Using KMS backend

Define your provider's config:

```
# The following example shows provider using kms backend

provider "secret" {
  backend = "kms"
  config = {
    # You can use AWS_* environment variables, or specify the following options here:
    shared_credentials_file = "~/.aws/credentials"
    profile = "custom_profile"

    # Alternatively, configure AWS settings directly:
    # aws_access_key = "your_access_key"
    # aws_secret_key = "your_secret_key"
    # region = "your_region"
  }
}
```

Then, you will need to encrypt your secrets using KMS and set variables to encrypted values:

```
aws kms encrypt --key-id <your_key_id> --region <aws_region> --plaintext "<secret_password_here>" | jq -r ".CiphertextBlob"
```

Set variable to the encrypted value in `terraform.tfvars` :

```
my_secret="AQICAHgtTQGsSDH8txmi3mOt4SDnq6Nb8/3yzY8w/EIHs4S6PAEWV/V6FR5m9DPo02vkTd53AAAAZTBjBgkqhkiG9w0BBwagVjBUAgEAME8GCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMCcL8O2e4qn2m57gsAgEQgCJg7l1u5O0jUudz99t1bLnfV/YOvmg+C5ekB968Egs2FGZB"
```

Use this together with data source to decrypt your value. In your `data_sources.tf`

```
data "secret" "my_secret" {
  encrypted_value = "${var.my_secret}"
}

# ...
# Access value of this secret anywhere referring to secret's "value" attribute

resource ... {
  param  = "${data.secret.my_secret.value}"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.8+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ hack/build/build.sh
...
$ $GOPATH/bin/terraform-provider-secret
...

```

In order to test the provider, you can simply run make test.

$ make test

In order to run the full suite of Acceptance tests, run make testacc.

Note: Acceptance tests create real resources, and often cost money to run.

$ make testacc
