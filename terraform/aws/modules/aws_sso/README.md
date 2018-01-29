# AWS SSO using Google

A Terraform module to create an Identity Provider in IAM and some default roles to be used for SSO.

## Example

```hcl
module "google_sso" {
  source             = "github.com/deliveryhero/dhh-system-engineering/terraform/aws/modules/aws_sso"
  idp_data_file_path = "path/to/GoogleIDPMetadata-my-domain.com.xml"
}
```

## Setup

It is assumed that your G Suite Administrator has already setup the SAML provider and custom user attributes. If not, here is a guide:

https://medium.com/proud2becloud/single-sign-on-with-g-suite-on-the-amazon-web-services-console-d506fda88c90

You need to obtain the `GoogleIDPMetadata` XML file from your G Suite Administrator in order to use this module.

Once the module has been added to your terraform code base and terraform has been run, your G Suite Administrator will also need to add attributes to each Google user in order to allow them access to roles within the AWS account. Once this is all complete, a new application called `Amazon Web Services` will be visible in your Google Apps list, top right corner of a Google web page.

## AWS CLI authentication with SSO

This guide will show on how to get your credentials per account in your config file.

To use this tool you need to install it with pip

```
pip install aws-google-auth
```

After installation is possible to use:

```
aws-google-auth -u email@domain.com -S <GOOGLE_SP_ID>  -I <GOOGLE_IDP_ID>

Failed to import U2F libraries, U2F login unavailable. Other methods can still continue.
Google Password:
MFA token: 123456
[  1] arn:aws:iam::1234567890:role/sso/sso-administrator
[  2] arn:aws:iam::0987654321:role/sso/sso-administrator
Type the number (1 - 2) of the role to assume: 1


Assuming arn:aws:iam::1234567890:role/sso/sso-administrator
Credentials Expiration: 2018-01-25 16:11:43+01:00

export AWS_ACCESS_KEY_ID='ASIAJFPE557GYRMSNISA' AWS_SECRET_ACCESS_KEY='PZL0SgLCGKrjhmO/ZxVINRdDYPGT2EyZ8sMNyele' AWS_SESSION_TOKEN='eyJ3aGF0IjoiQXJlIHlvdSB0cnlpbmcgdG8gZG8_IiwiUmVhbGx5IjoiRGlkIHlvdSB0aG91Z2h0IiwidGhhdCI6Ikkgd291bGQgcHV0IGhlcmUgYSByZWFsIGNlcnRpZmljYXRlIiwiaXMiOiJiZXR0ZXIgeW91IHRyeSBpdCBpbiBvdGhlciBwbGFjZSIsImFsZyI6IkhTMjU2In0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.kbpcPLr822Uq2lXsrCClXpZ_CE4RbZlCurzkhF5r_d0' AWS_SESSION_EXPIRATION='2018-01-25T15:11:43+0000'
```

Parameters:
 - -S : Google SSO SP identifier  ($GOOGLE_SP_ID)
 - -I : Google SSO IDP identifier ($GOOGLE_IDP_ID)

Is possible to get these parameters from your G Suite Administrator or from `Amazon Web Services` application URL that was installed in your Google Apps list.
