package account

import "github.com/abdelrahman146/kyora/internal/platform/types/problem"

func ErrInvalidCredentials(err error) *problem.Problem {
	return problem.Unauthorized("invalid email or password").WithError(err).WithCode("account.invalid_credentials")
}

func ErrInvalidOrExpiredToken(err error) *problem.Problem {
	return problem.Unauthorized("invalid or expired token").WithError(err).WithCode("account.invalid_token")
}

func ErrUserAlreadyExists(err error) *problem.Problem {
	return problem.Conflict("user with this email already exists").WithError(err).WithCode("account.user_already_exists")
}

func ErrInvitationAlreadyExists(err error) *problem.Problem {
	return problem.Conflict("active invitation already exists for this email").WithError(err).WithCode("account.invitation_already_exists")
}

func ErrInvitationNotFound(err error) *problem.Problem {
	return problem.NotFound("invitation not found").WithError(err).WithCode("account.invitation_not_found")
}

func ErrInvitationExpired(err error) *problem.Problem {
	return problem.Forbidden("invitation has expired").WithError(err).WithCode("account.invitation_expired")
}

func ErrInvitationAlreadyAccepted(err error) *problem.Problem {
	return problem.Conflict("invitation has already been accepted").WithError(err).WithCode("account.invitation_already_accepted")
}

func ErrInvitationCannotBeRevoked(err error) *problem.Problem {
	return problem.Conflict("only pending invitations can be revoked").WithError(err).WithCode("account.invitation_cannot_be_revoked")
}

func ErrCannotUpdateOwnRole(err error) *problem.Problem {
	return problem.Forbidden("you cannot update your own role").WithError(err).WithCode("account.cannot_update_own_role")
}

func ErrCannotUpdateOwnerRole(err error) *problem.Problem {
	return problem.Forbidden("you cannot update the workspace owner's role").WithError(err).WithCode("account.cannot_update_owner_role")
}

func ErrUserNotInWorkspace(err error) *problem.Problem {
	return problem.NotFound("user is not a member of this workspace").WithError(err).WithCode("account.user_not_member")
}

func ErrAuthRateLimited(err error) *problem.Problem {
	return problem.TooManyRequests("too many requests").WithError(err).WithCode("account.rate_limited")
}

func ErrInvalidInvitationToken(err error) *problem.Problem {
	return problem.Unauthorized("invalid or expired token").WithError(err).WithCode("account.invalid_invitation_token")
}

func ErrAccountOperationFailed(err error) *problem.Problem {
	return problem.InternalError().WithError(err).WithCode("account.operation_failed")
}
