package auth

type User struct {
	ID         string
	Email      string
	Identities []Identity
}

func (u *User) Authenticate(providerType, providerKey, attempt string) bool {
	for _, userAlias := range u.Identities {
		if userAlias.Matches(providerType, providerKey) {
			return userAlias.Validate(attempt)
		}
	}
	return false
}

func (u *User) AddIdentity(identity Identity) *User {
	newIdentities := make([]Identity, 0, len(u.Identities)+1)
	found := false
	for _, userAlias := range u.Identities {
		if userAlias.Matches(identity.ProviderType(), identity.ProviderKey()) {
			newIdentities = append(newIdentities, identity)
			found = true
		} else {
			newIdentities = append(newIdentities, userAlias)
		}
	}
	if !found {
		newIdentities = append(newIdentities, identity)
	}
	return &User{ID: u.ID, Email: u.Email, Identities: newIdentities}
}
