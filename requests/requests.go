package requests

type (
	CreateUserRequest struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	EditAccountRequest struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Email      string `json:"email"`
		Department string `json:"department"`
	}

	ListAccountFilters struct {
		AccountID string `json:"account_id"`
		Name      string `json:"name"`
	}
	ListAccountsRequest struct {
		Filters ListAccountFilters
	}

	LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email"`
	}

	ResetPasswordRequest struct {
		ResetCode   string
		NewPassword string
	}
)

type (
	CreateSensorRequest struct {
		Name      string `json:"name" bson:"name"`
		IpAddress string `json:"ipAddress" bson:"ip_address"`
	}

	UpdateSensorRequest struct {
		ID        string `json:"id" bson:"id"`
		Name      string `json:"name" bson:"name"`
		IpAddress string `json:"ipAddress" bson:"ip_address"`
	}
)
