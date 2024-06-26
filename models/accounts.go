package models

import "strings"

type (
	Status string // Status of an account type

	Kind string // Kind of account
)

const (
	ActiveStatus    Status = "ACTIVE"
	SuspendedStatus Status = "SUSPENDED"

	AdminAccountKind    Kind = "SLABMARK.ACCOUNT.KIND.ADMIN"
	EmployeeAccountKind Kind = "SLABMARK.ACCOUNT.KIND.EMPLOYEE"
)

var (
	FieldId                = "_id"
	FieldAccountUsername   = "username"
	FieldAccountPassword   = "password"
	FieldAccountPhone      = "phone"
	FieldAccountEmail      = "email"
	FieldAccountFirstName  = "first_name"
	FieldAccountLastName   = "last_name"
	FieldAccountDepartment = "department"
	FieldAccountInfoId     = "account_info.id"
)

type (
	Role struct {
		Shared
		RoleName string `json:"role_name" bson:"role_name"`
	}

	Account struct {
		Shared    `bson:",inline"`
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
		FullName  string `json:"full_name" bson:"full_name"`
		Email     string `json:"email" bson:"email"`
		Status    Status `json:"status" bson:"status"`
		Password  string `json:"password" bson:"password"`
		Token     string `json:"token" bson:"-"`
	}
)

type AccountInfo struct {
	Id        string `json:"id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	FullName  string `json:"full_name"`
	Email     string `json:"email" bson:"email"`
}
type AccountInterface interface {
	GetFullName() string
}

func NewAccount() AccountInterface {
	return &Account{}
}

func (a Account) GetFullName() string {
	return a.FirstName + " " + a.LastName
}

func (a Account) GetUsername() string {
	return string(a.FirstName[0]) + strings.ToLower(a.LastName)
}
