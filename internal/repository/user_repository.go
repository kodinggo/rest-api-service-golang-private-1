package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type userRepository struct {
	db       *sql.DB
	esClient *elastic.Client
}

func NewUserRepository(db *sql.DB, esClient *elastic.Client) model.UserRepository {
	return &userRepository{db: db, esClient: esClient}
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (result *model.User, err error) {
	row := sq.Select("id", "username", "created_at").
		From("users").
		Where(sq.Eq{
			"id": id,
		}).
		RunWith(r.db).
		QueryRowContext(ctx)

	var user model.User
	if err = row.Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
		logrus.Errorf("failed when scanning data user, error: %v", err)
		return
	}
	result = &user

	return
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (result *model.User, err error) {
	row := sq.Select("id", "username", "password", "created_at").
		From("users").
		Where(sq.Eq{
			"username": username,
		}).
		RunWith(r.db).
		QueryRowContext(ctx)

	var user model.User
	if err = row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		logrus.Errorf("failed when scanning data user, error: %v", err)
		return
	}
	result = &user

	return
}
