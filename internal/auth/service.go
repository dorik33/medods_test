package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"time"

	"github.com/dorik33/medods_test/internal/config"
	"github.com/dorik33/medods_test/internal/models"
	"github.com/dorik33/medods_test/internal/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRefreshToken = fmt.Errorf("invalid refresh token")
)

type Service struct {
	store *store.Store
	cfg   *config.Config
}

func NewService(store *store.Store, cfg *config.Config) *Service {
	return &Service{
		store: store,
		cfg:   cfg,
	}
}

func (s *Service) GenerateTokens(userID uuid.UUID, ipAddress string) (string, string, error) {
	jti := uuid.New()
	accessToken, err := s.generateAccessToken(userID, ipAddress, jti)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(userID, ipAddress, jti)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RefreshToken(accessToken string, refreshToken string, ipAddress string) (string, string, error) {
	claims, err := s.verifyAccessToken(accessToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid access token: %w", err)
	}

	decodedRefreshToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token format: %w", err)
	}
	originalRefreshToken := string(decodedRefreshToken)

	refreshTokenRecord, err := s.store.RefreshTokenRepository.FindByJTI(context.Background(), claims.JTI)
	if err != nil {
		return "", "", fmt.Errorf("refresh token not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(refreshTokenRecord.TokenHash), []byte(originalRefreshToken)); err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if time.Now().After(refreshTokenRecord.ExpiresAt) {
		return "", "", fmt.Errorf("refresh token has expired")
	}
	if claims.IPAddress != ipAddress {
		err := s.SendEmail()
		if err != nil {
			fmt.Println(err)
		}
	}

	newAccessToken, newRefreshToken, err := s.GenerateTokens(claims.UserID, ipAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %w", err)
	}

	if err := s.store.RefreshTokenRepository.Revoke(context.Background(), refreshTokenRecord.ID); err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *Service) generateAccessToken(userID uuid.UUID, ip string, jti uuid.UUID) (string, error) {
	claims := &models.Claims{
		UserID:    userID,
		IPAddress: ip,
		JTI:       jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *Service) generateRefreshToken(userID uuid.UUID, ipAddress string, jti uuid.UUID) (string, error) {
	refreshToken := uuid.New().String()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash refresh token: %w", err)
	}
	expiresAt := time.Now().Add(s.cfg.RefreshTokenTTL)

	refreshTokenRecord := &models.RefreshToken{
		UserID:         userID,
		TokenHash:      string(hashedToken),
		AccessTokenJTI: jti,
		ExpiresAt:      expiresAt,
		IPAddress:      ipAddress,
	}

	if err := s.store.RefreshTokenRepository.Create(context.Background(), refreshTokenRecord); err != nil {
		return "", fmt.Errorf("failed to save refresh token to database: %w", err)
	}

	encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	return encodedRefreshToken, nil
}

func (s *Service) verifyAccessToken(accessToken string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token claims")
	}

	return claims, nil
}

func (s *Service) SendEmail() error {
	subject := "Предупреждение"
	body := "Вы вошли с нового ip"
	message := []byte("Subject: " + subject + "\r\n" +
		"From: " + s.cfg.SMTPConfig.EmailFrom + "\r\n" +
		"To: " + s.cfg.SMTPConfig.EmailTo + "\r\n" +
		"\r\n" + body + "\r\n")
	auth := smtp.PlainAuth("", s.cfg.SMTPConfig.EmailFrom, s.cfg.SMTPConfig.EmailPassword, s.cfg.SMTPConfig.SMTPServer)

	err := smtp.SendMail(s.cfg.SMTPConfig.SMTPServer+":"+s.cfg.SMTPConfig.SMTPPort, auth, s.cfg.SMTPConfig.EmailFrom, []string{s.cfg.SMTPConfig.EmailTo}, message)
	if err != nil {
		return fmt.Errorf("ошибка при отправке email: %w", err)
	}

	fmt.Println("Уведомление об изменении IP отправлено на", s.cfg.SMTPConfig.EmailTo)
	return nil
}
