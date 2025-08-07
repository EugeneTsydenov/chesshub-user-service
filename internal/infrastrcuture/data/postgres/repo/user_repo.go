package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/entity/user"
	domainerrors "github.com/EugeneTsydenov/chesshub-user-service/internal/domain/errors"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/interfaces"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/email"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/password"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/publicname"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/tag"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/infrastrcuture/data/postgres"
	postgreserrors "github.com/EugeneTsydenov/chesshub-user-service/internal/infrastrcuture/data/postgres/errors"
)

type PostgresSessionRepo struct {
	database     *postgres.Database
	queryFactory postgres.UserQueryFactory
}

var _ interfaces.UserRepo = new(PostgresSessionRepo)

func NewPostgresSessionRepository(db *postgres.Database, factory postgres.UserQueryFactory) *PostgresSessionRepo {
	return &PostgresSessionRepo{
		database:     db,
		queryFactory: factory,
	}
}

func (r *PostgresSessionRepo) Create(ctx context.Context, u *user.User) (*user.User, error) {
	query := `
		INSERT INTO users (
			email, public_name, tag, password, last_active_at
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING 
			id, email, public_name, tag, password, last_active_at, updated_at, created_at
	`

	row := r.database.Pool().QueryRow(ctx, query,
		u.Email().Value(),
		u.PublicName().Value(),
		u.Tag().Value(),
		u.Password().Value(),
		u.LastActiveAt(),
	)

	user, err := scanUser(row)
	if err != nil {
		return nil, postgreserrors.WrapWithMapper("PostgresUserRepo.Create", err, nil)
	}

	return user, nil
}

func (r *PostgresSessionRepo) GetByID(ctx context.Context, userID int64) (*user.User, error) {
	query := `
		SELECT 
			id, email, public_name, tag, password, last_active_at, updated_at, created_at
		FROM users
		WHERE id = $1
	`

	row := r.database.Pool().QueryRow(ctx, query, userID)

	user, err := scanUser(row)
	if err != nil {
		return nil, postgreserrors.WrapWithMapper("PostgresUserRepo.GetByID", err, func(e error) error {
			if errors.Is(e, pgx.ErrNoRows) {
				return domainerrors.ErrUserNotFound
			}
			return fmt.Errorf("scan error: %w", e)
		})
	}

	return user, nil
}

func (r *PostgresSessionRepo) Update(ctx context.Context, u *user.User) (*user.User, error) {
	query := `
		UPDATE users SET
			email = $1,
			public_name = $2,
			tag = $3,
			password = $4,
			last_active_at = $5,
			updated_at = $6,
			created_at = $7
		WHERE id = $8
		RETURNING 
			id, email, public_name, tag, password, last_active_at, updated_at, created_at
	`

	row := r.database.Pool().QueryRow(ctx, query,
		u.Email().Value(),
		u.PublicName().Value(),
		u.Tag().Value(),
		u.Password().Value(),
		u.LastActiveAt(),
		u.UpdatedAt(),
		u.CreatedAt(),
		u.ID(),
	)

	user, err := scanUser(row)
	if err != nil {
		return nil, postgreserrors.WrapWithMapper("PostgresUserRepo.Update", err, nil)
	}

	return user, nil
}

func (r *PostgresSessionRepo) Find(ctx context.Context, criteria *user.Criteria) ([]*user.User, error) {
	query, args, err := r.queryFactory.BuildQuery(criteria)
	if err != nil {
		return nil, postgreserrors.WrapWithMapper("PostgresUserRepo.Find query builder", err, nil)
	}

	rows, err := r.database.Pool().Query(ctx, query, args...)
	if err != nil {
		return nil, postgreserrors.WrapWithMapper("PostgresUserRepo.Find query", err, nil)
	}
	defer rows.Close()

	var users []*user.User

	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, postgreserrors.WrapWithMapper("PostgresUserRepo.Find scan row", err, nil)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, postgreserrors.WrapWithMapper("PostgresUserRepo.Find rows iteration", err, nil)
	}

	return users, nil
}

func scanUser(row pgx.Row) (*user.User, error) {
	var (
		id           int64
		emailStr     string
		publicName   string
		tagStr       string
		passwordStr  string
		lastActiveAt time.Time
		updatedAt    time.Time
		createdAt    time.Time
	)

	err := row.Scan(
		&id,
		&emailStr,
		&publicName,
		&tagStr,
		&passwordStr,
		&lastActiveAt,
		&updatedAt,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	emailVO, _ := email.New(emailStr)
	publicNameVO, _ := publicname.New(publicName)
	tagVO, _ := tag.New(tagStr)
	passwordVO := password.NewHashedPassword(passwordStr)

	builder := user.NewBuilder()

	return builder.
		WithID(id).
		WithEmail(emailVO).
		WithPublicName(publicNameVO).
		WithTag(tagVO).
		WithPassword(passwordVO).
		WithLastActiveAt(lastActiveAt).
		WithUpdatedAt(updatedAt).
		WithCreatedAt(createdAt).
		Build(), nil
}
