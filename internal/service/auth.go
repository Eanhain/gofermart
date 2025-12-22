package service

import (
	"context"
	"errors"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	hash "github.com/Eanhain/gofermart/internal/hash"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) AuthUser(ctx context.Context, user dto.UserInput) (bool, error) {
	tUser, err := s.auth.CheckUser(ctx, user)
	if err != nil {
		return false, err
	}
	ok := hash.VerifyUserHash(s.log, user, tUser)
	return ok, nil
}

func (s *Service) RegUser(ctx context.Context, user dto.UserInput) error {
	var pgErr *pgconn.PgError
	hashedUser := hash.CreateUserHash(s.log, user)
	err := s.auth.RegisterUser(ctx, hashedUser)

	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = domain.ErrConflict
		return err
	}
	id, err := s.auth.GetUserID(ctx, user.Login)
	if err != nil {
		return err
	}

	s.balance.InsertNewUserBalance(ctx, id, 0)

	return err
}
