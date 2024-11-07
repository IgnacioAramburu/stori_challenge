package validation

const (
	ErrFieldRequired = "%s must be provided"
	ErrAgeTooLow     = "age must be at least 18, instead given: %d"
	ErrEmailFormat   = "email must be given in mail format <local>@<domain>.<top-level-domain>, instead given: %s"
)
