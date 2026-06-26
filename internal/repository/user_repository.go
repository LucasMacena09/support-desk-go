package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"support-desk-go/internal/models"
)

var ErrUserNotFound = errors.New("usuário não cadastrado")
var ErrUserAlreadyExists = errors.New("usuário já cadastrado com esse email")

type UserRepository struct {
	rdb *redis.Client
}

func NewUserRepository(rdb *redis.Client) *UserRepository {
	return &UserRepository{rdb: rdb}
}

func userKey(email string) string {
	return fmt.Sprintf("user:%s", email)
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	exists, err := r.rdb.Exists(ctx, userKey(user.Email)).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		return ErrUserAlreadyExists
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, userKey(user.Email), data, 0).Err()
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	data, err := r.rdb.Get(ctx, userKey(email)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}