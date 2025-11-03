package account

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrInvalidCredentials(err error) *problem.Problem {
	return problem.Unauthorized("invalid email or password").WithError(err)
}

func ErrInvalidOrExpiredToken(err error) *problem.Problem {
	return problem.Unauthorized("invalid or expired token").WithError(err)
}
