package response

import (
	"github.com/tejiriaustin/narx_api/models"
)

func SingleAccountResponse(account *models.Account) map[string]interface{} {
	return map[string]interface{}{
		"email":     account.Email,
		"password":  account.Password,
		"firstName": account.FirstName,
		"lastName":  account.LastName,
		"fullName":  account.FullName,
		"status":    account.Status,
		"token":     account.Token,
	}
}

func MultipleAccountResponse(accounts []models.Account) interface{} {
	m := make([]map[string]interface{}, 0, len(accounts))
	for _, a := range accounts {
		m = append(m, SingleAccountResponse(&a))
	}
	return m
}

func SingleSensorResponse(sensor *models.Sensor) map[string]interface{} {
	return map[string]interface{}{
		"_id":          sensor.ID.Hex(),
		"name":         sensor.Name,
		"ipAddress":    sensor.IpAddress,
		"status":       sensor.Status,
		"token":        sensor.Token,
		"account_info": sensor.AccountInfo,
	}
}

func MultipleSensorResponse(sensors []models.Sensor) interface{} {
	m := make([]map[string]interface{}, 0, len(sensors))
	for _, a := range sensors {
		m = append(m, SingleSensorResponse(&a))
	}
	return m
}
