package auth0

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Authenticate(user, pass string) (bool, error) {
	return false, nil
}
