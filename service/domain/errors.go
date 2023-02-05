package domain

import "fmt"

var (
	ErrTokenCredentialsExpired = fmt.Errorf("token is expired")
)
