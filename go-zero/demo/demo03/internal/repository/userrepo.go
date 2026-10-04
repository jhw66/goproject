package repository

import (
	"context"
	"demo03/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Find(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	user, err := r.FindByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).
		Save(&user).Error
}
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&model.User{}, id)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
