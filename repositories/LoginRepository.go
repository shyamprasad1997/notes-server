package repositories

import (
	"context"
	"errors"
	"notes-server/db"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"notes-server/utils"
)

type loginRepository struct {
	db     db.DB
	logger *loggers.Logger
}

func NewLoginRepository(db db.DB, logger *loggers.Logger) interfaces.ILoginRepository {
	return &loginRepository{db: db, logger: logger}
}

// Login - Checks if user exists in the db and checks if the password matches based on email
func (r *loginRepository) Login(ctx context.Context, request models.LoginRequest) (models.LoginRepoResponse, error) {
	r.logger.Info(ctx, "Entering loginRepository.Login()")
	defer r.logger.Info(ctx, "Exiting loginRepository.Login()")
	// Query DB to validate email and password
	txn := r.db.Txn(ctx, false)
	row, err := txn.First("user", "email", request.Email)
	if err != nil {
		txn.Abort()
		r.logger.Warn(ctx, "error in loginRepository.Login(), error from txn.First()", err)
		return models.LoginRepoResponse{}, err
	}
	response, ok := row.(*models.User)
	if !ok {
		txn.Abort()
		r.logger.Warn(ctx, "error in loginRepository.Login(), failed to fetch user data", err)
		return models.LoginRepoResponse{}, errors.New("failed to fetch user data, data not existing")
	}
	if response.Password != request.Password {
		txn.Abort()
		r.logger.Warn(ctx, "error in loginRepository.Login(), invalid credentials", err)
		return models.LoginRepoResponse{}, errors.New("invalid credentials")
	}
	txn.Commit()
	return models.LoginRepoResponse{
		Email: response.Email,
		Name:  response.Name,
	}, nil
}

// SignUp - Creates a new user in the db
func (r *loginRepository) SignUp(ctx context.Context, request models.SignUpRequest) error {
	r.logger.Info(ctx, "Entering loginRepository.SignUp()")
	defer r.logger.Info(ctx, "Exiting loginRepository.SignUp()")
	txn := r.db.Txn(ctx, true)

	user := models.User{Name: request.Name, Email: request.Email, Password: request.Password, Id: utils.NewID()}
	err := txn.Insert("user", &user)
	if err != nil {
		txn.Abort()
		r.logger.Warn(ctx, "error in loginRepository.SignUp(), error from txn.Insert()", err)
		return err
	}
	txn.Commit()
	return nil
}

// CheckIfUserExists - Checks if a user exists
func (r *loginRepository) CheckIfUserExists(ctx context.Context, email string) (bool, error) {
	r.logger.Info(ctx, "Entering loginRepository.CheckIfUserExists()")
	defer r.logger.Info(ctx, "Exiting loginRepository.CheckIfUserExists()")
	// Query DB to validate email
	txn := r.db.Txn(ctx, false)
	row, err := txn.First("user", "email", email)
	if err != nil {
		txn.Abort()
		r.logger.Warn(ctx, "error in loginRepository.CheckIfUserExists(), error from txn.First()", err)
		return false, err
	}
	txn.Commit()
	_, ok := row.(*models.User)
	if ok {
		return true, nil
	}
	return false, nil
}

// ValidateUser - validate creds
func (r *loginRepository) ValidateUser(ctx context.Context, email, name string) error {
	r.logger.Info(ctx, "Entering loginRepository.ValidateUser()")
	defer r.logger.Info(ctx, "Exiting loginRepository.ValidateUser()")
	// Query DB to validate email and name
	txn := r.db.Txn(ctx, false)
	row, err := txn.First("user", "email", email)
	if err != nil {
		txn.Abort()
		r.logger.Warn(ctx, "error in loginRepository.ValidateUser(), error from txn.First()", err)
		return err
	}
	txn.Commit()
	user, ok := row.(*models.User)
	if !ok {
		return errors.New("invalid user")
	}
	if user.Name != name {
		return errors.New("invalid user")
	}
	return nil
}
