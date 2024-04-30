package services

import (
	"context"
	"github.com/tejiriaustin/narx_api/models"
	"github.com/tejiriaustin/narx_api/publisher"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/utils"
)

type (
	AccountsServiceInterface interface {
		CreateUser(ctx context.Context,
			input AddAccountInput,
			passwordGen utils.StrGenFunc,
			accountsRepo *repository.Repository[models.Account],
		) (*models.Account, error)

		EditAccount(ctx context.Context,
			input EditAccountInput,
			accountsRepo *repository.Repository[models.Account],
		) (*models.Account, error)

		LoginUser(ctx context.Context,
			input LoginUserInput,
			accountsRepo *repository.Repository[models.Account],
		) (*models.Account, error)

		ForgotPassword(ctx context.Context,
			input ForgotPasswordInput,
			accountsRepo *repository.Repository[models.Account],
			publisher publisher.PublishInterface,
		) (*models.Account, error)

		ResetPassword(ctx context.Context,
			input ResetPasswordInput,
			accountsRepo *repository.Repository[models.Account],
		) (*models.Account, error)

		ListAccounts(ctx context.Context,
			input ListAccountReportsInput,
			accountsRepo *repository.Repository[models.Account],
		) ([]models.Account, *repository.Paginator, error)
	}
)
