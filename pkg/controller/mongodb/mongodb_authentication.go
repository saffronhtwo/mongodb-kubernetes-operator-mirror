package mongodb

import (
	mdbv1 "github.com/mongodb/mongodb-kubernetes-operator/pkg/apis/mongodb/v1"
	"github.com/mongodb/mongodb-kubernetes-operator/pkg/authentication/scram"
	"github.com/mongodb/mongodb-kubernetes-operator/pkg/automationconfig"
	"github.com/mongodb/mongodb-kubernetes-operator/pkg/kube/secret"
)

const (
	scramShaOption = "SCRAM"
)

// noOpEnabler performs no changes, leaving authentication settings untouched
type noOpEnabler struct{}

func (n noOpEnabler) EnableAuth(auth automationconfig.Auth) automationconfig.Auth {
	return auth
}

// getAuthenticationEnabler returns a type that is able to configure the automation config's
// authentication settings
func getAuthenticationEnabler(getUpdateCreator secret.GetUpdateCreator, mdb mdbv1.MongoDB) (automationconfig.AuthEnabler, error) {
	if !mdb.Spec.Security.Authentication.Enabled {
		return noOpEnabler{}, nil
	}

	// currently, just enable auth if it's in the list as there is only one option
	if containsAuthMode(mdb.Spec.Security.Authentication.Modes, scramShaOption) {
		enabler, err := scram.EnsureAgentSecret(getUpdateCreator, mdb.ScramCredentialsNamespacedName())
		if err != nil {
			return noOpEnabler{}, err
		}
		return enabler, nil
	}
	return noOpEnabler{}, nil
}

func containsAuthMode(slice []mdbv1.AuthMode, s mdbv1.AuthMode) bool {
	for _, elem := range slice {
		if elem == s {
			return true
		}
	}
	return false
}
