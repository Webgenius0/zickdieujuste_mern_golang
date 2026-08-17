package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gotickets/internal/auth"
	"gotickets/internal/domain/user/dto"
	"gotickets/internal/upload"

	"github.com/google/uuid"
)

var (
	ErrDuplicateEmail      = errors.New("email is already registered")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrInvalidOTP          = errors.New("OTP is invalid, expired, or already used")
	ErrAccountDeleted      = errors.New("account has been deleted")
)

// Service defines the business logic contract for the user domain.
// It depends only on the Repository interface and never imports echo or gorm.
type Service interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(rawToken string) (*dto.AuthResponse, error)
	Logout(rawToken string) error
	ForgotPassword(req dto.ForgotPasswordRequest) error
	ResetPassword(req dto.ResetPasswordRequest) error

	GetProfileByEmail(email string) (*dto.ProfileResponse, error)
	UpdateProfileByEmail(email string, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error)
	ChangePasswordByEmail(email string, req dto.ChangePasswordRequest) error
	DeleteAccountByEmail(email string) error
	GetUserIDByEmail(email string) (uuid.UUID, error)

	RegisterDevice(userID uuid.UUID, req dto.RegisterDeviceRequest) error
}

type service struct {
	repo     Repository
	jwt      auth.JWTService
	uploader upload.Uploader // may be nil if upload not configured
}

// NewService creates a new user Service. uploader may be nil — avatar upload will return an error if so.
func NewService(repo Repository, jwt auth.JWTService, uploader upload.Uploader) Service {
	return &service{repo: repo, jwt: jwt, uploader: uploader}
}

// ──────────────────────────────────────────────────────────────────────────────
// Auth operations
// ──────────────────────────────────────────────────────────────────────────────

func (s *service) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check for duplicate email
	existing, err := s.repo.GetUserByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, ErrDuplicateEmail
	}

	u := &User{
		Name:         req.Name,
		Email:        req.Email,
		AuthProvider: AuthProviderEmail,
	}
	if err := u.HashPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.CreateUser(u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.issueTokenPair(u)
}

func (s *service) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// TODO: add rate limiting here (e.g. Redis-backed per-email limiter)
	u, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := u.CheckPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokenPair(u)
}

func (s *service) RefreshToken(rawToken string) (*dto.AuthResponse, error) {
	// Validate JWT signature/expiry first
	_, err := s.jwt.ValidateToken(rawToken, true)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	hash := hashToken(rawToken)
	rt, err := s.repo.GetRefreshToken(hash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Rotate: revoke the old token
	if err := s.repo.RevokeRefreshToken(rt.TokenHash); err != nil {
		return nil, fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	u, err := s.repo.GetUserByID(rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.issueTokenPair(u)
}

func (s *service) Logout(rawToken string) error {
	if rawToken == "" {
		return nil // idempotent
	}
	hash := hashToken(rawToken)
	return s.repo.RevokeRefreshToken(hash)
}

func (s *service) ForgotPassword(req dto.ForgotPasswordRequest) error {
	// TODO: add rate limiting here (e.g. Redis-backed per-email limiter)

	// Always return 200 regardless of whether the email exists (no user enumeration)
	u, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		// Email not found — return silently
		return nil
	}

	// Invalidate any prior unused OTPs
	_ = s.repo.InvalidatePendingOTPs(req.Email)

	code, err := generateOTPCode()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	otp := &OTP{
		Email:     req.Email,
		Code:      code,
		Purpose:   "PASSWORD_RESET",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if u != nil {
		otp.UserID = &u.ID
	}

	if err := s.repo.CreateOTP(otp); err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	// TODO: dispatch email with OTP code (email provider integration goes here)
	fmt.Printf("[OTP] Email: %s Code: %s (expires in 10m)\n", req.Email, code)

	return nil
}

func (s *service) ResetPassword(req dto.ResetPasswordRequest) error {
	// TODO: add rate limiting here (e.g. Redis-backed per-email limiter)

	otp, err := s.repo.GetValidOTP(req.Email, req.OTP)
	if err != nil {
		return ErrInvalidOTP
	}

	u, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if err := u.HashPassword(req.NewPassword); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	if err := s.repo.UpdateUser(u); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return s.repo.MarkOTPUsed(otp.ID)
}

// ──────────────────────────────────────────────────────────────────────────────
// Profile operations
// ──────────────────────────────────────────────────────────────────────────────

func (s *service) GetProfileByEmail(email string) (*dto.ProfileResponse, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return toProfileResponse(u), nil
}

func (s *service) GetUserIDByEmail(email string) (uuid.UUID, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return uuid.Nil, err
	}
	return u.ID, nil
}

func (s *service) UpdateProfileByEmail(email string, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Location != nil {
		u.Location = req.Location
	}
	if req.ThemePreference != nil {
		u.ThemePreference = ThemePreference(*req.ThemePreference)
	}
	if req.LanguagePreference != nil {
		u.LanguagePreference = *req.LanguagePreference
	}
	if err := s.repo.UpdateUser(u); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return toProfileResponse(u), nil
}

func (s *service) UpdateAvatarURLByEmail(email string, avatarURL string) (*dto.ProfileResponse, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	u.AvatarURL = &avatarURL
	if err := s.repo.UpdateUser(u); err != nil {
		return nil, fmt.Errorf("failed to update avatar: %w", err)
	}
	return toProfileResponse(u), nil
}

func (s *service) ChangePasswordByEmail(email string, req dto.ChangePasswordRequest) error {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if err := u.CheckPassword(req.OldPassword); err != nil {
		return errors.New("current password is incorrect")
	}
	if err := u.HashPassword(req.NewPassword); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return s.repo.UpdateUser(u)
}

func (s *service) DeleteAccountByEmail(email string) error {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	return s.repo.SoftDeleteUser(u.ID)
}

func (s *service) ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error {
	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if err := u.CheckPassword(req.OldPassword); err != nil {
		return errors.New("current password is incorrect")
	}
	if err := u.HashPassword(req.NewPassword); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return s.repo.UpdateUser(u)
}

func (s *service) DeleteAccount(userID uuid.UUID) error {
	return s.repo.SoftDeleteUser(userID)
}

func (s *service) UpdateAvatar(ctx context.Context, userID uuid.UUID, file interface{}, _ interface{}) (*dto.AvatarResponse, error) {
	if s.uploader == nil {
		return nil, errors.New("upload service not configured")
	}
	return nil, errors.New("use handler-level UploadAvatar — this method is a handler-level concern")
}

// ──────────────────────────────────────────────────────────────────────────────
// Device operations
// ──────────────────────────────────────────────────────────────────────────────

func (s *service) RegisterDevice(userID uuid.UUID, req dto.RegisterDeviceRequest) error {
	dt := &DeviceToken{
		UserID:   userID,
		Token:    req.Token,
		Platform: Platform(req.Platform),
	}
	return s.repo.UpsertDeviceToken(dt)
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

// issueTokenPair generates a new JWT access+refresh token pair, persists the refresh token, and returns both.
func (s *service) issueTokenPair(u *User) (*dto.AuthResponse, error) {
	// jwt.GenerateToken takes (userId uint, name, email) — existing jwt.go contract.
	// We pass 0 as the uint placeholder since we now use UUID as the primary key.
	// The user's UUID is embedded in the refresh token record for lookup.
	accessToken, refreshToken, err := s.jwt.GenerateToken(0, u.Name, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	rt := &RefreshToken{
		UserID:    u.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.repo.CreateRefreshToken(rt); err != nil {
		return nil, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// hashToken returns a SHA-256 hex hash of the token string for safe DB storage.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

// generateOTPCode generates a random 5-digit numeric string.
func generateOTPCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 5)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}

// toProfileResponse maps a User entity to a ProfileResponse DTO.
func toProfileResponse(u *User) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		ID:                 u.ID,
		Name:               u.Name,
		Email:              u.Email,
		AuthProvider:       string(u.AuthProvider),
		Location:           u.Location,
		AvatarURL:          u.AvatarURL,
		ThemePreference:    string(u.ThemePreference),
		LanguagePreference: u.LanguagePreference,
		IsPremium:          u.IsPremium,
		TermsAcceptedAt:    u.TermsAcceptedAt,
		CreatedAt:          u.CreatedAt,
	}
}
