package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/events/notifications"
	"github.com/tejiriaustin/narx_api/models"
	"github.com/tejiriaustin/narx_api/publisher"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/utils"
)

type AccountsService struct {
	conf *env.Environment
}

func NewAccountsService(conf *env.Environment) *AccountsService {
	return &AccountsService{
		conf: conf,
	}
}

type (
	AddAccountInput struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}
	EditAccountInput struct {
		Id        string
		FirstName string
		LastName  string
		Email     string
	}

	LoginUserInput struct {
		Username string
		Password string
	}
	ForgotPasswordInput struct {
		Email string
	}
	ResetPasswordInput struct {
		ResetToken  string
		NewPassword string
	}
	Claims struct {
		Exp           time.Time
		Authorization bool
		jwt.StandardClaims
		Content any
	}
	AccountListFilters struct {
		Query string // for partial free hand lookups
	}

	ListAccountReportsInput struct {
		Pager
		Projection *repository.QueryProjection
		Sort       *repository.QuerySort
		Filters    AccountListFilters
	}
)

func (s *AccountsService) CreateUser(ctx context.Context,
	input AddAccountInput,
	passwordGen utils.StrGenFunc,
	accountsRepo *repository.Repository[models.Account],
) (*models.Account, error) {
	if input.Email == "" {
		return nil, errors.New("email is required")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}
	if input.FirstName == "" {
		return nil, errors.New("first name is required")
	}

	qf := repository.NewQueryFilter().AddFilter("email", input.Email)
	matchedUser, err := accountsRepo.FindOne(ctx, qf, nil, nil)
	if err != nil && err != repository.NoDocumentsFound {
		return nil, err
	}

	if matchedUser.Email == input.Email {
		return nil, errors.New("user with this email already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, errors.New("password hashing failed")
	}

	now := time.Now()

	account := models.Account{
		Shared: models.Shared{
			ID:        primitive.NewObjectID(),
			CreatedAt: &now,
		},
		Username:  input.FirstName,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Status:    models.ActiveStatus,
		Password:  string(passwordHash),
	}

	account.FullName = account.GetFullName()
	account.Username = account.GetUsername()
	account.FullName = account.GetFullName()

	acct, err := accountsRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return &acct, nil
}

func (s *AccountsService) EditAccount(ctx context.Context,
	input EditAccountInput,
	accountsRepo *repository.Repository[models.Account],
) (*models.Account, error) {

	fields := map[string]interface{}{}

	if input.FirstName != "" {
		fields[models.FieldAccountFirstName] = input.FirstName
	}
	if input.LastName != "" {
		fields[models.FieldAccountLastName] = input.LastName
	}
	if input.Email != "" {
		fields[models.FieldAccountEmail] = input.Email
	}
	
	updates := map[string]interface{}{
		"$set": fields,
	}

	filter := repository.NewQueryFilter().AddFilter(models.FieldId, input.Id)
	err := accountsRepo.UpdateMany(ctx, filter, updates)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *AccountsService) LoginUser(ctx context.Context,
	input LoginUserInput,
	accountsRepo *repository.Repository[models.Account],
) (*models.Account, error) {

	filter := repository.NewQueryFilter().AddFilter(models.FieldAccountUsername, input.Username)

	account, err := accountsRepo.FindOne(ctx, filter, nil, nil)
	if err != nil {
		if err == repository.NoDocumentsFound {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	token, err := s.generateSignedToken(ctx, models.AccountInfo{
		Id:         account.ID.Hex(),
		FirstName:  account.FirstName,
		LastName:   account.LastName,
		FullName:   account.FullName,
		Email:      account.Email,
		Department: account.Department,
	})
	if err != nil {
		return nil, errors.New("an error occurred: " + err.Error())
	}

	account.Token = token

	return &account, nil
}

func (s *AccountsService) ForgotPassword(ctx context.Context,
	input ForgotPasswordInput,
	accountsRepo *repository.Repository[models.Account],
	publisher publisher.PublishInterface,
) (*models.Account, error) {

	filter := repository.NewQueryFilter().AddFilter(models.FieldAccountEmail, input.Email)

	account, err := accountsRepo.FindOne(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	token, err := s.generateSignedToken(ctx, models.AccountInfo{
		Id:         account.ID.Hex(),
		FirstName:  account.FirstName,
		LastName:   account.LastName,
		FullName:   account.FullName,
		Email:      account.Email,
		Department: account.Department,
	})
	if err != nil {
		return nil, err
	}

	event := map[string]interface{}{
		"id":         account.ID.Hex(),
		"first_name": account.FirstName,
		"last_name":  account.LastName,
		"full_name":  account.FullName,
		"department": account.Department,
		"email":      account.Email,
		"link":       token,
	}
	err = publisher.Publish(ctx, notifications.ForgotPasswordNotification, event)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *AccountsService) ResetPassword(ctx context.Context,
	input ResetPasswordInput,
	accountsRepo *repository.Repository[models.Account],
) (*models.Account, error) {

	accountInfo := &models.AccountInfo{}
	err := s.verifySignedToken(ctx, input.ResetToken, accountInfo)
	if err != nil {
		return nil, err
	}

	filter := repository.NewQueryFilter().AddFilter(models.FieldAccountEmail, accountInfo.Email)

	account, err := accountsRepo.FindOne(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 8)
	if err != nil {
		return nil, errors.New("couldn't generate password")
	}

	account.Password = string(passwordHash)

	updatedAccount, err := accountsRepo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	return &updatedAccount, nil
}

func (s *AccountsService) ListAccounts(ctx context.Context,
	input ListAccountReportsInput,
	accountsRepo *repository.Repository[models.Account],
) ([]models.Account, *repository.Paginator, error) {

	filter := repository.NewQueryFilter()

	if input.Filters.Query != "" {
		freeHandFilters := []map[string]interface{}{
			{"first_name": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"last_name": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"full_name": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"phone": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"email": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"username": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
		}
		filter.AddFilter("$or", freeHandFilters)
	}

	account, paginator, err := accountsRepo.Paginate(ctx, filter, input.PerPage, input.Page, input.Projection, input.Sort)
	if err != nil {
		return nil, nil, err
	}

	return account, paginator, nil
}

func (s *AccountsService) generateSignedToken(ctx context.Context, content any) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Exp:           time.Now().Add(3600 * time.Minute),
		Authorization: true,
		Content:       content,
	})

	pkey := s.conf.GetAsBytes(env.JwtSecret)
	tokenString, err := token.SignedString(pkey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AccountsService) verifySignedToken(ctx context.Context, token string, target any) error {

	if token != "" {
		return errors.New("token not set")
	}
	claims := &Claims{
		Content: target,
	}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt was signed with an unknown signature")
		}
		return s.conf.GetAsBytes(env.JwtSecret), nil
	})
	if err != nil {
		return err
	}

	return nil
}
