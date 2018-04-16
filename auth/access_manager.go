package auth

type AccessManager interface {
	GrantAccess(accountID string) error
	RevokeAccess(accountID string) error
	Accounts() ([]string, error)
}
