package postgres

import (
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/entity/user"
	"github.com/Masterminds/squirrel"
)

type (
	UserQueryFactory interface {
		BuildQuery(criteria *user.Criteria) (string, []interface{}, error)
	}
	userQueryFactory struct {
		psql squirrel.StatementBuilderType
	}
)

func NewUserQueryFactory() UserQueryFactory {
	return &userQueryFactory{
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (f *userQueryFactory) BuildQuery(criteria *user.Criteria) (string, []interface{}, error) {
	query := f.psql.
		Select(
			"id",
			"email",
			"public_name",
			"tag",
			"password",
			"last_active_at",
			"updated_at",
			"created_at",
		).
		From("users")

	query = f.applyFilters(query, criteria)

	sql, args, err := query.ToSql()

	return sql, args, err
}

func (f *userQueryFactory) applyFilters(query squirrel.SelectBuilder, criteria *user.Criteria) squirrel.SelectBuilder {
	if criteria == nil {
		return query
	}

	if criteria.Email != nil {
		query = query.Where(squirrel.Eq{"email": *criteria.Email})
	}

	if criteria.Tag != nil {
		query = query.Where(squirrel.Eq{"tag": *criteria.Tag})
	}

	if criteria.PublicName != nil {
		query = query.Where(squirrel.Eq{"public_name": *criteria.PublicName})
	}

	if criteria.LastActiveBefore != nil {
		query = query.Where(squirrel.LtOrEq{"last_active_at": *criteria.LastActiveBefore})
	}

	if criteria.LastActiveAfter != nil {
		query = query.Where(squirrel.GtOrEq{"last_active_at": *criteria.LastActiveAfter})
	}

	if criteria.UpdatedBefore != nil {
		query = query.Where(squirrel.LtOrEq{"updated_at": *criteria.UpdatedBefore})
	}

	if criteria.UpdatedAfter != nil {
		query = query.Where(squirrel.GtOrEq{"updated_at": *criteria.UpdatedAfter})
	}

	if criteria.CreatedBefore != nil {
		query = query.Where(squirrel.LtOrEq{"created_at": *criteria.CreatedBefore})
	}

	if criteria.CreatedAfter != nil {
		query = query.Where(squirrel.GtOrEq{"created_at": *criteria.CreatedAfter})
	}

	return query
}
