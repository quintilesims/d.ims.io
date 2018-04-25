package auth

type AccountManager interface {
	GrantAccess(accountID string) error
	RevokeAccess(accountID string) error
	Accounts() ([]string, error)
}
