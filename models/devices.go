package models

type Devices struct {
	Shared      `bson:",inline"`
	AccountInfo AccountInfo `json:"accountInfo" bson:"account_info"`
	DeviceToken string      `json:"deviceToken" bson:"device_token"`
}
