package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	// Using WithContext to propagate context to the database layer
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Let service handle the domain error
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	return r.db.WithContext(ctx).Create(user).Error
}
