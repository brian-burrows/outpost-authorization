package auth

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID         string     `bson:"_id"`
	Email      string     `bson:"email"`
	Identities []Identity `bson:"identities"`
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

func (u *User) MarshalBSON() ([]byte, error) {
	// We create an anonymous struct with public fields just for the marshaler
	return bson.Marshal(struct {
		ID         string     `bson:"_id"`
		Email      string     `bson:"email"`
		Identities []Identity `bson:"identities"`
	}{
		ID:         u.ID,
		Email:      u.Email,
		Identities: u.Identities,
	})
}

func (u *User) UnmarshalBSON(data []byte) error {
	type rawIdentity struct {
		Type string `bson:"provider_type"`
		Key  string `bson:"provider_key"`
	}

	raw := struct {
		ID         string        `bson:"_id"`
		Email      string        `bson:"email"`
		Identities []rawIdentity `bson:"identities"`
	}{}

	if err := bson.Unmarshal(data, &raw); err != nil {
		return err
	}

	u.ID = raw.ID
	u.Email = raw.Email
	u.Identities = make([]Identity, len(raw.Identities))
	for i, rid := range raw.Identities {
		u.Identities[i] = NewIdentity(rid.Type, rid.Key)
	}

	return nil
}
