# Delivery Hero System Engineering

This repository contains [Terraform](https://www.terraform.io/) modules, Kubernetes configuration and other best practice examples to give a head start when building a new environment on AWS or Kubernetes.

The intended audience is teams or projects that are starting with AWS or Kubernetes and want to begin with a good design.

Pull requests welcome!

## Guides

- [terraform-kubernetes-getting-started](docs/terraform-kubernetes-getting-started.md): A guide to getting started with Kubernetes on AWS using Terraform. Includes AWS and Kubernetes configuration and running an example application.

## Terraform modules

- [aws_sso](terraform/aws/modules/aws_sso): Creates an Identity Provider in IAM and some default roles to be used for SSO.
- [autoscaling_group_elb_app](terraform/aws/modules/autoscaling_group_elb_app): An example of EC2 autoscaling group and ELB.
- [kubernetes_kops_cluster](terraform/aws/modules/kubernetes_kops_cluster): A Kubernetes cluster with multi-AZ master based on [kops](https://github.com/kubernetes/kops).
- [lambda_kubernetes_deployer](terraform/aws/modules/lambda_kubernetes_deployer): A lambda function that deploys to a Kubernetes cluster when a container image is pushed to an ECR repository.
- [lambda_splunk_forwarder](terraform/aws/modules/lambda_splunk_forwarder): A lambda function for forwarding Cloudwatch logs to a Splunk HTTP events collector.
- [nat_gateway](terraform/aws/modules/nat_gateway): Creates multi-AZ NAT gateways, associated private subnets and route tables.
- [ssh_bastion](terraform/aws/modules/ssh_bastion): SSH bastion host in an ASG with a fixed EIP.
- [security_monkey](terraform/aws/modules/security_monkey): IAM role and associated policy to run Netflix's [Security Monkey](https://github.com/Netflix/security_monkey).
- [vpc](terraform/aws/modules/vpc): A VPC setup that includes public and private subnets for each AZ, route tables.

## Examples

- [base](terraform/aws/examples/base): A good starting point with a VPC and related resources, SSH bastion, some default security groups and S3 bucket for Terraform state.
- [kubernetes](terraform/aws/examples/kubernetes): Same as `base` example but with shared kubernetes resources and 2 kubernetes clusters added.

## Kubernetes Examples

- [Deployment, Service and Autoscaling](kubernetes/examples/deployment_service): A simple example of a Kubernetes deployment, different service types and a horizontal pod autoscaler.
- [Helm templating](kubernetes/examples/helm): A simple example using Helm to template Kubernetes resources.
- [Ingress Controller](kubernetes/examples/ingress): Using an Ingress controller and resource to split traffic across multiple applications.
- [ExternalDNS](kubernetes/examples/external-dns): A tool that automatically creates DNS records for Kubernetes resources.
- [Kubernetes Best Practices](https://speakerdeck.com/thesandlord/kubernetes-best-practices): A great presentation from a Google engineer.
- [Kubernetes Autoscaling](kubernetes/examples/autoscaling): How to set up autoscaling for both deployment and cluster nodes.
- [Kubernetes OpenID Connect authenticator](kubernetes/examples/k8s-oidc-authenticator): API Authentication using Google OIDC
- [Kubernetes dashboard with OAuth2](kubernetes/examples/dashboard-oauth): Run the Kubernetes dashboard behind OAuth2

## Other useful tools

#### Keymaker

Lightweight SSH key management on AWS EC2. Add public SSH keys to IAM users and then they can log into EC2 hosts.

https://github.com/kislyuk/keymaker

#### Elasticsearch on Kubernetes

https://github.com/pires/kubernetes-elasticsearch-cluster

#### Invokust

Run [Locust](http://locust.io/) load tests on AWS Lambda.

https://github.com/FutureSharks/invokust

#### Helm

[Helm](https://github.com/kubernetes/helm) is a powerful tool for creating templates for Kubernetes resources, creating reproducible builds or for packaging and installing predefined configurations for services.

## Dos and Don'ts

- Do make pull requests to this repository.
- Don't bother using a NAT gateway unless you specifically need a fixed source IP address for outgoing traffic.
- Do enable [MFA](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html) for your IAM accounts.
- Do store your Terraform state in a bucket.
- Do use an internal Route53 zone to hold records for RDS endpoints, ES endpoints, Elasticache endpoints etc.
- Do terminate SSL on ELBs and forward as HTTP in VPC. This means you never need to deal with SSL or certificates on instances.
- Do use a SSH bastion for all SSH connections and restrict SSH access by IP ranges.
- Consider registering an external domain in Route53. It only costs a few dollars and you can have a free SSL certificate. Then use this domain and certificate for all external ELBs.
- Don't have instances that are not part of an Austoscaling Group.
- Do write [Terraform modules](https://www.terraform.io/docs/configuration/modules.html) to reduce duplicated code.
