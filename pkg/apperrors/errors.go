package apperrors

import "errors"

var (
	// Report state machine errors
	ErrNotPermitted        = errors.New("operation not permitted in current state")
	ErrMethodNotPermitted  = errors.New("structurally blocked in current state")
	ErrSuperApprovalNeeded = errors.New("super approval required before transition")
	ErrConflictExists      = errors.New("unresolved conflict exists")
	ErrSchemaFrozen        = errors.New("schema version is frozen")

	// Sync errors
	ErrOutboxFull = errors.New("outbox full - connect to server before continuing")

	// Auth and permission errors
	ErrPermissionDenied = errors.New("permission denied")
	ErrSectionNotFound  = errors.New("section not found in schema")
	ErrNotFound         = errors.New("not found")
	ErrUnauthorised     = errors.New("unauthorised")

	// License errors - these have specific UX flows
	ErrSeatLimitReached         = errors.New("seat limit reached")
	ErrLicenseRevoked           = errors.New("license revoked")
	ErrDeploymentRevoked        = errors.New("deployment license revoked")
	ErrGracePeriodExpired       = errors.New("offline grace period expired")
	ErrMachineMismatch          = errors.New("machine fingerprint mismatch")
	ErrInvalidSignature         = errors.New("invalid license signature")
	ErrNoLocalLicense           = errors.New("no local license file found")
	ErrFeatureNotLicensed       = errors.New("this feature is not included in your license")
	ErrDeploymentLicenseInvalid = errors.New("deployment license is invalid")
)
