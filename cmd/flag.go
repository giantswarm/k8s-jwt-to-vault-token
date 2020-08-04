package cmd

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagVaultAddr        = "vault-address"
	flagVaultPolicy      = "vault-policy"
	flagVaultTokenSecret = "vault-token-secret"
)

type flag struct {
	VaultAddr        string
	VaultPolicy      string
	VaultTokenSecret string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.VaultAddr, flagVaultAddr, "", `Vault address to request token. E.g.: "https://127.0.0.1".`)
	cmd.Flags().StringVar(&f.VaultPolicy, flagVaultPolicy, "", `Existing vault policy for requested token.`)
	cmd.Flags().StringVar(&f.VaultTokenSecret, flagVaultTokenSecret, "", `Kubernetes secret name, where vault token is stored.`)
}

func (f *flag) Validate() error {
	if f.VaultAddr == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagVaultAddr)
	}
	if f.VaultPolicy == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagVaultPolicy)
	}
	if f.VaultTokenSecret == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagVaultTokenSecret)
	}
	return nil
}
