package models

type (
	Sensor struct {
		Shared      `bson:",inline"`
		AccountInfo AccountInfo `json:"accountInfo" bson:"account_info"`
		Name        string      `json:"name" bson:"name"`
		IpAddress   string      `json:"ip_address" bson:"ip_address"`
		Status      string      `json:"status" bson:"status"`
		Token       string      `json:"token" bson:"-"`
	}
)
