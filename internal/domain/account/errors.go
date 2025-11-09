package account

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrInvalidCredentials(err error) *problem.Problem {
	return problem.Unauthorized("invalid email or password").WithError(err)
}

func ErrInvalidOrExpiredToken(err error) *problem.Problem {
	return problem.Unauthorized("invalid or expired token").WithError(err)
}

func ErrUserAlreadyExists(err error) *problem.Problem {
	return problem.Conflict("user with this email already exists").WithError(err)
}

func ErrInvitationAlreadyExists(err error) *problem.Problem {
	return problem.Conflict("active invitation already exists for this email").WithError(err)
}

func ErrInvitationNotFound(err error) *problem.Problem {
	return problem.NotFound("invitation not found").WithError(err)
}

func ErrInvitationExpired(err error) *problem.Problem {
	return problem.Forbidden("invitation has expired").WithError(err)
}

func ErrInvitationAlreadyAccepted(err error) *problem.Problem {
	return problem.Conflict("invitation has already been accepted").WithError(err)
}

func ErrCannotUpdateOwnRole(err error) *problem.Problem {
	return problem.Forbidden("you cannot update your own role").WithError(err)
}

func ErrCannotUpdateOwnerRole(err error) *problem.Problem {
	return problem.Forbidden("you cannot update the workspace owner's role").WithError(err)
}

func ErrUserNotInWorkspace(err error) *problem.Problem {
	return problem.NotFound("user is not a member of this workspace").WithError(err)
}
