// Handles the user's authentication
package account

import (
	"github.com/blocklords/gosds/app/service"
	"github.com/blocklords/gosds/common/data_type/key_value"
)

// Requeter to the SDS Service. It's either a developer or another SDS service.
type Account struct {
	Id             uint64           `json:"id,omitempty"`    // Auto incremented for every new developer
	PublicKey      string           `json:"public_key"`      // Public Key for authentication.
	Organization   string           `json:"organization"`    // Organization
	NonceTimestamp uint64           `json:"nonce_timestamp"` // Nonce since the last usage. Only acceptable for developers
	service        *service.Service // If the account is another service, then this parameter keeps the data. Otherwise this parameter is a nil.
}

type Accounts []*Account

// Creates a new Account for a developer.
func NewDeveloper(id uint64, public_key string, nonce_timestamp uint64, organization string) *Account {
	return &Account{
		Id:             id,
		PublicKey:      public_key,
		NonceTimestamp: nonce_timestamp,
		Organization:   organization,
		service:        nil,
	}
}

// Creates a new Account for a service
func NewService(service *service.Service) *Account {
	return &Account{
		Id:             0,
		NonceTimestamp: 0,
		PublicKey:      service.PublicKey,
		Organization:   "",
		service:        service,
	}
}

// Creates an account for a service
func NewServices(services []*service.Service) []*Account {
	accounts := make([]*Account, 0, len(services))
	for i, s := range services {
		accounts[i] = NewService(s)
	}

	return accounts
}

func (account *Account) IsDeveloper() bool {
	return account.service == nil
}

func (account *Account) IsService() bool {
	return account.service != nil
}

func ParseJson(raw key_value.KeyValue) (*Account, error) {
	public_key, err := raw.GetString("public_key")
	if err != nil {
		return nil, err
	}
	service, err := service.GetByPublicKey(public_key)
	if err != nil {
		id, err := raw.GetUint64("id")
		if err != nil {
			return nil, err
		}
		nonce_timestamp, err := raw.GetUint64("nonce_timestamp")
		if err != nil {
			return nil, err
		}

		organization, err := raw.GetString("organization")
		if err != nil {
			return nil, err
		}
		return NewDeveloper(id, public_key, nonce_timestamp, organization), nil
	} else {
		return NewService(service), nil
	}
}

///////////////////////////////////////////////////////////
//
// Group operations
//
///////////////////////////////////////////////////////////

func NewAccounts(new_accounts ...*Account) Accounts {
	accounts := make(Accounts, len(new_accounts))
	copy(accounts, new_accounts)

	return accounts
}

func NewAccountsFromJson(raw_accounts []key_value.KeyValue) (Accounts, error) {
	accounts := make(Accounts, len(raw_accounts))

	for i, raw := range raw_accounts {
		account, err := ParseJson(raw)
		if err != nil {
			return nil, err
		}

		accounts[i] = account
	}

	return accounts, nil
}

func (accounts Accounts) Add(new_accounts ...*Account) Accounts {
	for _, account := range new_accounts {
		accounts = append(accounts, account)
	}

	return accounts
}

func (accounts Accounts) Remove(new_accounts ...*Account) Accounts {
	for _, account := range new_accounts {
		for i := range accounts {
			if account.PublicKey == accounts[i].PublicKey {
				accounts = append(accounts[:i], accounts[i+1:]...)
				return accounts
			}
		}
	}

	return accounts
}

func (accounts Accounts) PublicKeys() []string {
	public_keys := make([]string, len(accounts))

	for i := range accounts {
		public_keys[i] = accounts[i].PublicKey
	}

	return public_keys
}
