package response

import "github.com/tejiriaustin/narx_api/models"

func SingleAccountResponse(account *models.Account) map[string]interface{} {
	return map[string]interface{}{
		"email":      account.Email,
		"password":   account.Password,
		"firstName":  account.FirstName,
		"lastName":   account.LastName,
		"username":   account.Username,
		"fullName":   account.FullName,
		"phone":      account.Phone,
		"department": account.Department,
		"status":     account.Status,
		"token":      account.Token,
	}
}

func MultipleAccountResponse(accounts []models.Account) interface{} {
	m := make([]map[string]interface{}, 0, len(accounts))
	for _, a := range accounts {
		m = append(m, SingleAccountResponse(&a))
	}
	return m
}
