package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	callbackURL  = os.Getenv("CALLBACK_URL")
	clusterName  = os.Getenv("CLUSTER_NAME")
	caData       = os.Getenv("CA_DATA")
	expectedHostedDomain = os.Getenv("ALLOWED_DOMAIN")
)

var expectedEmails = map[string]bool{}

func init() {
	emailsStr := os.Getenv("ALLOWED_EMAIL_ADDRESSES")
	if emailsStr != "" {
		emails := strings.Split(emailsStr, ",")
		for _, addr := range emails {
			expectedEmails[addr] = true
		}
	}
}

const oauthURL = "https://accounts.google.com/o/oauth2/auth?redirect_uri=%s&response_type=code&client_id=%s&scope=openid+email+profile&approval_prompt=force&access_type=offline"
const tokenURL = "https://www.googleapis.com/oauth2/v3/token"
const userInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo"
const idpIssuerURL = "https://accounts.google.com"


type GoogleConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type UserInfo struct {
	Email string `json:"email"`
}

type HostedDomain struct {
	HostedDomain string `json:"hd"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

// Get the id_token and refresh_token from google
func getTokens(code string) (*TokenResponse, error) {
	val := url.Values{}
	val.Add("grant_type", "authorization_code")
	val.Add("redirect_uri", callbackURL)
	val.Add("client_id", clientID)
	val.Add("client_secret", clientSecret)
	val.Add("code", code)

	resp, err := http.PostForm(tokenURL, val)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got: %d calling %s", resp.StatusCode, tokenURL)
	}
	if err != nil {
		return nil, err
	}
	tr := &TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(tr)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func getUserEmail(accessToken string) (string, error) {
	uri, _ := url.Parse(userInfoURL)
	q := uri.Query()
	q.Set("alt", "json")
	q.Set("access_token", accessToken)
	uri.RawQuery = q.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Got: %d calling %s", resp.StatusCode, tokenURL)
	}
	if err != nil {
		return "", err
	}
	ui := &UserInfo{}
	err = json.NewDecoder(resp.Body).Decode(ui)
	if err != nil {
		return "", err
	}
	return ui.Email, nil
}

func getHostedDomain(accessToken string) (string, error) {
	uri, _ := url.Parse(userInfoURL)
	q := uri.Query()
	q.Set("alt", "json")
	q.Set("access_token", accessToken)
	uri.RawQuery = q.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Got: %d calling %s", resp.StatusCode, tokenURL)
	}
	if err != nil {
		return "", err
	}
	hd := &HostedDomain{}
	err = json.NewDecoder(resp.Body).Decode(hd)
	if err != nil {
		return "", err
	}
	return hd.HostedDomain, nil
}

func googleRedirect() http.Handler {
	redirectURL := fmt.Sprintf(oauthURL, callbackURL, clientID)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectURL, http.StatusFound)
	})
}

func googleCallback() http.Handler {
	outputTemplate := template.Must(template.New("shell").Parse(`
<html><head>
<style type="text/css">
.shell-wrap {
  width: auto;
  margin: 20px 50px 50px 50px;
  box-shadow: 0 0 30px rgba(0,0,0,0.4);
  border-radius: 3px;
}

.shell-top-bar {
  text-align: center;
  color: #525252;
  padding: 5px 0;
  margin: 0;
  text-shadow: 1px 1px 0 rgba(255,255,255,0.5);
  font-size: 0.85em;
  border: 1px solid #CCCCCC;
  border-bottom: none;

  border-top-left-radius: 3px;
  border-top-right-radius: 3px;

  background: #f7f7f7; /* Old browsers */
  background: linear-gradient(to bottom,  #f7f7f7 0%,#B8B8B8 100%); /* W3C */
}

.shell-body {
  margin: 0;
  padding: 5px;
  list-style: none;
  background: #141414;
  color: #45D40C;
  font: 0.8em 'Andale Mono', Consolas, 'Courier New';
  line-height: 1.6em;
  overflow: scroll;

  border-bottom-right-radius: 3px;
  border-bottom-left-radius: 3px;
}

.shell-body li.cmd:before {
  content: '$';
  position: absolute;
  left: 0;
  top: 0;
}

.shell-body li.comment:before {
  content: '#';
  position: absolute;
  left: 0;
  top: 0;
}

.shell-body li.comment {
  color: #450CD4;
}

.shell-body li.comment, li.space, .shell-top-bar, #btn-copy {
  user-select: none;
  -moz-user-select: none;
  -webkit-user-select: none;
  -ms-user-select: none;
}

.shell-body li {
  word-wrap: break-word;
  position: relative;
  padding: 0 0 0 15px;
}
#clipboard {
  width: 0px;
  height: 0px;
  overflow: hidden;
  border: none;
  position: fixed;
  top: -20px;
}
#btn-copy {
  color: #fff;
  background-color: #007bff;
  border: 1px solid #007bff;
  border-radius: .3rem;
  text-align: center;
  vertical-align: middle;
  font-family: sans-serif;
  margin: 0px auto;
  width: 10em;
  padding: 4px 20px;
  transition: color .15s ease-in-out, background-color .15s ease-in-out, border-color .15s ease-in-out, box-shadow .15s ease-in-out;
  box-shadow: 2px 2px 2px 1px rgba(100,100,100,.60);
}
#btn-copy:hover{
  color: #fff;
  background-color: #0069d9;
  border-color: #0062cc;
}
</style>
<script type="text/javascript">
function copy_to_clipboard() {
  let clipboard = document.getElementById("clipboard");
  let commands = Array.from(document.querySelectorAll('.cmd')).map((n) => n.innerText).join('\n');
  clipboard.value = commands.replace(/\n+/g, "\n");
  clipboard.select();

  if(!document.execCommand("copy")) {
    alert("Sorry, your browser does not allow direct access to clipboard.\n" +
      "Please select all (no worries) and copy manually.");
  }

  document.getElementById("btn-copy").style.boxShadow = 'none';
}
</script>
</head><body>
<textarea id="clipboard"></textarea>
<div id="btn-copy" onclick="copy_to_clipboard();">Copy to clipboard</div>
<div class="shell-wrap">
  <p class="shell-top-bar">~/work/docker-setup</p>
  <ul class="shell-body" id="commands">
<li class='comment'>Run the following command to configure a kubernetes user for use with 'kubectl'
<br><br>
</li>

<li class='cmd'>
kubectl config set-credentials {{ .email }} \<br>
--auth-provider=oidc 
--auth-provider-arg=client-id={{ .clientID }} \<br>
--auth-provider-arg=client-secret={{ .clientSecret }}
--auth-provider-arg=id-token={{ .idToken }} \<br>
--auth-provider-arg=idp-issuer-url={{ .issuerURL }}
--auth-provider-arg=refresh-token={{ .refreshToken }}; \<br>
</li>

<li class='space'>&nbsp;</li>
<li class='comment'>Configure your context to use your Google user account</li>
<li class='comment'>Go to to your kube config and update user under context part to:</li>
<li class='comment'>'''</li>
<li class='comment'>user: {{ .email }}</li>
<li class='comment'>'''<br><br></li>

<li class='cmd'>
kubectl config
set-context {{ .context }}
--cluster {{ .cluster }}
--user {{ .email }}; \
</li>
{{ if .caData }}
<li class='cmd'>echo '{{ .caData }}' | \<br>
tr -cd 'a-zA-Z0-9=/' | base64 --decode > ca-data; \</li>
<li class='cmd'>
kubectl config
set-cluster {{ .cluster }}
--embed-certs=true
--certificate-authority=ca-data
--server=https://api.{{ .cluster }}; \</li>
<li class='cmd'>rm ca-data; \</li>
{{ else }}
<li class='cmd'>
kubectl config
set-cluster {{ .cluster }}
--insecure-skip-tls-verify=true
--server=https://api.{{ .cluster }}; \</li>
{{ end }}
<li class='cmd'>kubectl config use-context {{ .context }}; \</li>
<li class='space'>&nbsp;</li>
<li class='comment'>Test connection by running (you should see list of nodes in cluster)<br><br></li>
<li class='cmd'>kubectl get nodes;</li>
  </ul>
</div>
</body></html>
	`))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		tokResponse, err := getTokens(code)
		if err != nil {
			log.Printf("Error getting tokens: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		email, err := getUserEmail(tokResponse.AccessToken)
		if err != nil {
			log.Printf("Error getting user email: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hostedDomain, err := getHostedDomain(tokResponse.AccessToken)
		if err != nil {
			log.Printf("Error getting user hosted domain: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, emailAllowed := expectedEmails[email]
		if hostedDomain != expectedHostedDomain && !emailAllowed {
			log.Printf("Error hosted domain does not match (expected domain %s, but email was %s)\n", expectedHostedDomain, email)
			http.Error(w, "Forbidden", 403)
			return
		}


		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		outputCluster := "my-cluster.my-domain.com"
		if clusterName != "" {
			outputCluster = clusterName
		}

		err = outputTemplate.Execute(w, map[string]string {
			"email": email,
			"clientID": clientID,
			"clientSecret": clientSecret,
			"idToken": tokResponse.IdToken,
			"issuerURL": idpIssuerURL,
			"refreshToken": tokResponse.RefreshToken,
			"cluster": outputCluster,
			"context": outputCluster,
			"caData": caData,
		});

		if err != nil {
			log.Println("failed to write about response: %s", err)
		}
	})
}

func main() {
	m := http.NewServeMux()

	m.Handle("/", googleRedirect())
	m.Handle("/callback", googleCallback())

	http.Handle("/", m)
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
