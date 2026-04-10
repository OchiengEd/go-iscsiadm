package command

import "errors"

var (
	ErrSessionNotFound = errors.New("iSCSI session not found")
	// ErrLogoutFailed is returned when iSCSI logout issues are encountered
	ErrLoginFailed = errors.New("iSCSI target login failed")
	// ErrLogoutFailed is returned when iSCSI logout issues are encountered
	ErrLogoutFailed = errors.New("iSCSI target logout failed")
	// ErrPermissionDenied is returned when limited OS permissions are insufficient
	ErrPermissionDenied = errors.New("insufficient permission")
	// ErrResourceNotFound returned when nodes/targets/portal are not found
	ErrResourceNotFound = errors.New("iSCSI resource(s) not found")
	// ErrSessionExists is returned when  you attempt to login when an active
	// iSCSI session already exists
	ErrSessionExists = errors.New("iSCSI session already exists")
)

// Common iscsiadm exit codes
const (
	// iSCSI session not found
	ExitCodeSessionNotFound int = 2
	// generic iSCSI login failure
	ExitCodeLoginFailure int = 5
	// iSCSI logout failed
	ExitCodeLogoutFailure int = 10
	// insufficient OS permissions to access
	// iscsid or execute iscsiadm commands
	ExitCodeAccessDenied int = 13
	// iSCSI session alread exists
	ExitCodeSessionExists int = 15
	// no records/targets/sessions/portal
	// found to execute operation on
	ExitCodeObjectsNotFound int = 21
	// login failed due to authz failure
	ExitCodeAuthorization int = 24
)
