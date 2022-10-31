package token

import "time"

/**Makes is an interface for managing tokens */
type Maker interface {
	/**Create a new token for a specific username and duration */
	CreateToken(user_id int64, duration time.Duration) (string, error)

	/**check if generated token is valid */
	VerifyToken(token string) (*Payload, error)
}
