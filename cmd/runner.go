package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/k8s-jwt-to-vault-token/internal/label"
	"github.com/giantswarm/k8s-jwt-to-vault-token/pkg/project"
)

type runner struct {
	flag   *flag
	logger micrologger.Logger
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) PersistentPreRun(cmd *cobra.Command, args []string) error {
	fmt.Printf("Version = %#q\n", project.Version())
	fmt.Printf("Git SHA = %#q\n", project.GitSHA())
	fmt.Printf("Command = %#q\n", cmd.Name())
	fmt.Println()

	return nil
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	jwt, err := readJWT()
	if err != nil {
		return microerror.Mask(err)
	}

	vaultToken, err := vaultLogin(jwt, r.flag.VaultRole, r.flag.VaultAddr)
	if err != nil {
		return microerror.Mask(err)
	}
	if vaultToken == "" {
		return microerror.Maskf(executionFailedError, "vault token must not be empty")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return microerror.Mask(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ensureVaultTokenSecretExist(clientSet, vaultToken, r.flag.VaultTokenSecretName, r.flag.VaultTokenSecretNamespace)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func ensureVaultTokenSecretExist(clientSet *kubernetes.Clientset, vaultToken, vaultTokenSecretName, vaultTokenSecretNamespace string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vaultTokenSecretName,
			Namespace: vaultTokenSecretNamespace,
			Labels: map[string]string{
				label.ManagedBy: name,
			},
		},
		Data: map[string][]byte{
			"token": []byte(vaultToken),
		},
	}

	_, err := clientSet.CoreV1().Secrets(vaultTokenSecretNamespace).Get(vaultTokenSecretName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Printf("Creating secret %#q for vault token...\n", vaultTokenSecretName)
		_, err = clientSet.CoreV1().Secrets(vaultTokenSecretNamespace).Create(secret)
		if err != nil {
			return microerror.Mask(err)
		}
		fmt.Printf("Secret %#q for vault token has been created.\n", vaultTokenSecretName)

	} else {
		fmt.Printf("Updating secret %#q for vault token...\n", vaultTokenSecretName)
		_, err = clientSet.CoreV1().Secrets(vaultTokenSecretNamespace).Update(secret)
		if err != nil {
			return microerror.Mask(err)
		}
		fmt.Printf("Secret %#q for vault token has been updated.\n", vaultTokenSecretName)
	}

	return nil
}

func readJWT() (string, error) {
	jwtFilePath := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	fmt.Printf("Reading service account JWT from file %#q ...\n", jwtFilePath)

	jwt, err := ioutil.ReadFile(jwtFilePath)

	if err != nil {
		return "", microerror.Mask(err)
	}

	return string(jwt), nil
}

func vaultLogin(jwt, role, vaultAddr string) (string, error) {
	fmt.Printf("Logging into vault Kubernetes backend...\n")

	vaultAuthEndpoint := "v1/auth/kubernetes/login"
	url := fmt.Sprintf("%s/%s", vaultAddr, vaultAuthEndpoint)
	values := map[string]string{"jwt": jwt, "role": role}

	jsonValues, err := json.Marshal(values)
	if err != nil {
		return "", microerror.Mask(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValues))
	if err != nil {
		return "", microerror.Mask(err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return "", microerror.Mask(err)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", microerror.Mask(err)
		}

		return "", microerror.Maskf(executionFailedError, "expected code %d got %d: %#q", http.StatusOK, resp.StatusCode, string(body))
	}
	defer resp.Body.Close()

	type responseData struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}

	var tokenData responseData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", microerror.Mask(err)
	}

	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return tokenData.Auth.ClientToken, nil
}
