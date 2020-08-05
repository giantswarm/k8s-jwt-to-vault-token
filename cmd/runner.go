package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/giantswarm/k8s-jwt-to-vault-token/pkg/project"
)

const (
	jwtFilePath       = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	vaultAuthEndpoint = "v1/auth/kubernetes/login"
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
	jwt, err := readJWTFromFile(jwtFilePath)
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Printf("Token: %s", string(jwt))

	return nil
}

func readJWTFromFile(filepath string) ([]byte, error) {
	jwt, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, microerror.Mask(err)
	}

	return jwt, nil
}
