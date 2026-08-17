package dto

// RegisterRequest is the body for POST /api/v1/auth/register.
type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest is the body for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest is the optional body for POST /api/v1/auth/refresh.
// The refresh token can also be read from the "refresh_token" cookie.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// ForgotPasswordRequest is the body for POST /api/v1/auth/forgot-password.
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest is the body for POST /api/v1/auth/reset-password.
type ResetPasswordRequest struct {
	Email       string `json:"email"        validate:"required,email"`
	OTP         string `json:"otp"          validate:"required,len=5"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UpdateProfileRequest is the body for PUT /api/v1/users/me.
type UpdateProfileRequest struct {
	Name               *string `json:"name"                omitempty:"true" validate:"omitempty,min=2,max=100"`
	Location           *string `json:"location"            omitempty:"true"`
	ThemePreference    *string `json:"theme_preference"    omitempty:"true" validate:"omitempty,oneof=IVORY NAVY"`
	LanguagePreference *string `json:"language_preference" omitempty:"true"`
}

// ChangePasswordRequest is the body for PUT /api/v1/users/me/password.
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"     validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// RegisterDeviceRequest is the body for POST /api/v1/devices.
type RegisterDeviceRequest struct {
	Token    string `json:"token"    validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=IOS ANDROID"`
}
