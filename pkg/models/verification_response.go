package models

// VerificationResponse represents the response from a verification process.
type VerificationResponse struct {
	// ErrorMessage can be used to provide more information about the error
	// encountered during the verification process.
	ErrorMessage string

	// Success represents a boolean value indicating whether the verification
	// process was successful or not.
	Success bool
}
