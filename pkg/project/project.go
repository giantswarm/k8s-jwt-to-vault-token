package project

var (
	description = "CLI tool to request Vault token via Kubernetes service account JWT."
	gitSHA      = "n/a"
	name        = "k8s-jwt-to-vault-token"
	source      = "https://github.com/giantswarm/k8s-jwt-to-vault-token"
	version     = "0.1.1-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
