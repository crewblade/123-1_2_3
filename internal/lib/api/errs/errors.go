package errs

import "fmt"

var (
	ErrUserIsNotAuthorized   = fmt.Errorf("user is not authorized")
	ErrUserDoesNotHaveAccess = fmt.Errorf("user does not have access")
	ErrBannerNotFound        = fmt.Errorf("banner not found")
	ErrBannerForTagNotFound  = fmt.Errorf("banner for tag not found")
)
