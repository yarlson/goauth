package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/model"
)

type UserRepository struct {
	db *gorm.DB
}

type TokenRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *UserRepository) GetOrCreateUser(ctx context.Context, email, firstName, lastName, pictureURL string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where(model.User{Email: email}).FirstOrCreate(&user, model.User{
			Email:      email,
			FirstName:  firstName,
			LastName:   lastName,
			PictureURL: pictureURL,
		})
		return result.Error
	})

	return user, err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).First(&user, "id = ?", id.String())

	return user, result.Error
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	token := model.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days expiration
	}

	return r.db.WithContext(ctx).Create(&token).Error
}

func (r *TokenRepository) VerifyRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	var token model.RefreshToken
	result := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", refreshToken, time.Now()).First(&token)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return token.UserID, nil
}

func (r *TokenRepository) UpdateRefreshToken(ctx context.Context, oldRefreshToken, newRefreshToken string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var token model.RefreshToken
		if err := tx.Where("token = ?", oldRefreshToken).First(&token).Error; err != nil {
			return err
		}

		token.Token = newRefreshToken
		token.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)

		return tx.Save(&token).Error
	})
}

func (r *TokenRepository) InvalidateRefreshToken(ctx context.Context, refreshToken string) error {
	return r.db.WithContext(ctx).Where("token = ?", refreshToken).Delete(&model.RefreshToken{}).Error
}
