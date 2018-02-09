# Kubernetes OpenID Connect authenticator

Kubernetes provides various methods to [authenticate against the API](https://kubernetes.io/docs/admin/authentication/), here we will use OpenID Connect (OIDC) authentication with RBAC enabled.

Credit to [cu12](https://github.com/cu12/k8s-oidc-helper) for this adaptation of the original k8s-oidc-helper from micahhausler: https://github.com/micahhausler/k8s-oidc-helper

## Security Properties

**You should never use OIDC authentication without RBAC enabled**

The setup described in this README checks authentication and authorization as follows:
  1. The k8s-oidc-authenticator will use google to authenticate a user (and the associated email address)
     as a valid and semi-current google user.
  2. The k8s-oidc-authenticator then checks the google-vetted identity against an expected domain and
     a whitelist of specific email addresses. If it matches it provides the user with *all* data
     needed to access the k8s cluster as this user or any other google user under his control.
     **This must be coupled with email-specific access controls on the cluster to actually restrict access 
     to google users under the same organizational control as the cluster.**
  3. The cluster's role-based access control deployed via helm authorizes certain email addresses
     to perform administrative roles.

## Installation

#### Create OAuth credentials

You need a Google Cloud project.

  1. Go to https://console.cloud.google.com/apis/credentials
  2. Click `Create credentials` > `OAuth client ID` > `Web application`
  3. Choose a name. e.g. `k8s-oidc-authenticator`
  4. Enter a URL in `Authorized redirect URIs`. e.g. https://oidc-auth.my-domain.com/callback
  5. Record the `Client ID` and `Client secret` for later

#### Configure Kubernetes API for OIDC

To configure the Kubernetes API server you need to add some parameters. This is detailed here: https://kubernetes.io/docs/admin/authentication/#configuring-the-api-server

If you used [kops](https://github.com/kubernetes/kops) to build your Kubernetes cluster, you can edit it like this:

```
kops edit cluster --name cluster.my-domain.com
```

And insert these lines:

```yaml
spec:
  api:
...
  kubeAPIServer:
    authorizationMode: RBAC
    oidcClientID:  <Client ID from GCP>
    oidcIssuerURL: https://accounts.google.com
    oidcUsernameClaim: email
```

If you don't have the `kubeAPIServer` section, create it.

#### Create Kubernetes resources

There is a Helm chart in `kubernetes/examples/k8s-oidc-authenticator/helm-chart` that will create:

  - A deployment and service for the k8s-oidc-authenticator container
  - Ingress resource
  - RoleBinding and ClusterRoleBinding to allow admin access or access to specific namespaces

Here is an example values file to be passed to the Helm chart:

```yaml
ingress:
  # Hostname used for the ingress resource
  hostname: k8s-oauth.my-domain.com
  # Optionally white list IPs on the Ingress resource
  whitelist_ips:
    - 201.1.2.3
    - 184.1.2.3

oidc:
  # These values are from the API & Services page of GCP
  client_id: "2222222-xxxxxx.apps.googleusercontent.com"
  client_secret: "xxxxxxx"
  allowed_domain: my-domain.com

role_bindings:
  # Email addresses given admin access
  cluster-admin:
    - user1@my-domain.com
    - user2@my-domain.com
  namespaces:
    # Email addresses given access to specific namespaces
    default:
      - user3@my-domain.com
      - user4@my-domain.com

# The name of your cluster used in KUBECONFIG file
cluster_name: cluster1.my-domain.com

# Optional CA data to verify cluster certificate. Otherwise `insecure-skip-tls-verify` option will be used
cluster_ca_data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMwekNDQWJ1Z0F3SUJBZ0lNRlFoOTdjTmdTNmpYMjBOa01BMEdDU3FHU0liM0RRRUJDd1VBTUJVeEV6QVIKYmxhaApibGFoCmJsYWgKTUlJQzB6Q0NBYnVnQXdJQkFnSU1GUWg5N2NOZ1M2algyME5rTUEwR0NTcUdTSWIzRFFFQkN3VUFNQlV4RXpBUgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
```

Install the chart:

```
helm install kubernetes/examples/k8s-oidc-authenticator/helm-chart --name k8s-oidc-authenticator -f my-k8s-oidc-authenticator-values.yaml
```

## Use

Simply open your configured ingress hostname in a browser, `k8s-oauth.my-domain.com` from above, and you will be redirected to Google for authentication. If your email address matches one of from `role_bindings` section of the configuration, you will see the required `kubectl` commands in order to connect to the cluster.
