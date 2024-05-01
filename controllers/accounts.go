package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/models"
	"github.com/tejiriaustin/narx_api/publisher"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/requests"
	"github.com/tejiriaustin/narx_api/response"
	"github.com/tejiriaustin/narx_api/services"
	"github.com/tejiriaustin/narx_api/utils"
)

type AccountsController struct {
	conf *env.Environment
}

func NewAccountController(conf *env.Environment) *AccountsController {
	return &AccountsController{
		conf: conf,
	}
}

func (c *AccountsController) SignUp(
	passwordGen utils.StrGenFunc,
	acctService services.AccountsServiceInterface,
	accountsRepo *repository.Repository[models.Account],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.CreateUserRequest

		err := ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, "Bad Request", nil)
			return
		}

		input := services.AddAccountInput{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		}

		user, err := acctService.CreateUser(ctx, input, passwordGen, accountsRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", response.SingleAccountResponse(user))
	}
}

func (c *AccountsController) Login(
	acctService services.AccountsServiceInterface,
	accountsRepo *repository.Repository[models.Account],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.LoginUserRequest

		err := ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		input := services.LoginUserInput{
			Email:    req.Email,
			Password: req.Password,
		}

		user, err := acctService.LoginUser(ctx, input, accountsRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", response.SingleAccountResponse(user))
	}
}

func (c *AccountsController) EditAccount(
	acctService services.AccountsServiceInterface,
	accountsRepo *repository.Repository[models.Account],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.CreateUserRequest

		err := ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, "Bad Request", nil)
			return
		}

		input := services.EditAccountInput{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		}

		user, err := acctService.EditAccount(ctx, input, accountsRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", response.SingleAccountResponse(user))
	}
}

func (c *AccountsController) LogOut() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie("auth", "", -1, "/", c.conf.GetAsString(env.FrontendUrl), false, true)
		ctx.String(http.StatusOK, "Cookie has been deleted")
	}
}

func (c *AccountsController) ForgotPassword(
	acctService services.AccountsServiceInterface,
	accountsRepo *repository.Repository[models.Account],
	publisher publisher.PublishInterface,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.ForgotPasswordRequest

		err := ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		input := services.ForgotPasswordInput{
			Email: req.Email,
		}

		user, err := acctService.ForgotPassword(ctx, input, accountsRepo, publisher)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", response.SingleAccountResponse(user))
	}
}

func (c *AccountsController) ResetPassword(
	acctService services.AccountsServiceInterface,
	accountsRepo *repository.Repository[models.Account],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.ResetPasswordRequest

		err := ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		input := services.ResetPasswordInput{
			NewPassword: req.NewPassword,
			ResetToken:  req.ResetCode,
		}

		user, err := acctService.ResetPassword(ctx, input, accountsRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", response.SingleAccountResponse(user))
	}
}

// ---------------------------------------------------------------- Irrelevant ---------------------------------------------------------------- //

func (c *AccountsController) ListAccounts(
	acctService services.AccountsServiceInterface,
	accountsRepo *repository.Repository[models.Account],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		_, err := GetAccountInfo(ctx, c.conf.GetAsBytes(env.JwtSecret))
		if err != nil {
			response.FormatResponse(ctx, http.StatusUnauthorized, "Unauthorized access", nil)
			return
		}

		var req requests.ListAccountsRequest

		err = ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		input := services.ListAccountReportsInput{
			Pager: services.Pager{
				Page:    services.GetPageNumberFromContext(ctx),
				PerPage: services.GetPerPageLimitFromContext(ctx),
			},
			Filters: services.AccountListFilters{
				Query: ctx.Param("query"),
			},
		}
		accounts, paginator, err := acctService.ListAccounts(ctx, input, accountsRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		payload := map[string]interface{}{
			"records": response.MultipleAccountResponse(accounts),
			"meta":    paginator,
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", payload)
	}
}
