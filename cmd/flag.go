package cmd

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagVaultAddr                 = "vault-address"
	flagVaultPolicy               = "vault-policy"
	flagVaultTokenSecretName      = "vault-token-secret-name"      // nolint
	flagVaultTokenSecretNamespace = "vault-token-secret-namespace" // nolint
)

type flag struct {
	VaultAddr                 string
	VaultPolicy               string
	VaultTokenSecretName      string
	VaultTokenSecretNamespace string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.VaultAddr, flagVaultAddr, "", `Vault address to request token. E.g.: "https://127.0.0.1".`)
	cmd.Flags().StringVar(&f.VaultPolicy, flagVaultPolicy, "", `Existing vault policy for requested token.`)
	cmd.Flags().StringVar(&f.VaultTokenSecretName, flagVaultTokenSecretName, "", `Kubernetes secret name, where vault token is stored.`)
	cmd.Flags().StringVar(&f.VaultTokenSecretNamespace, flagVaultTokenSecretNamespace, "", `Kubernetes secret namespace, where vault token secret is stored.`)
}

func (f *flag) Validate() error {
	if f.VaultAddr == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagVaultAddr)
	}
	if f.VaultPolicy == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagVaultPolicy)
	}
	if f.VaultTokenSecretName == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagVaultTokenSecretName)
	}
	if f.VaultTokenSecretNamespace == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagVaultTokenSecretNamespace)
	}
	return nil
}
