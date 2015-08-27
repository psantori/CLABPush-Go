package server

import (
	"encoding/json"
	"strconv"
)

// CLabInfo represents the CLab subdictionary of the JSON body.
type CLabInfo struct {
	AppID    string // Application package.
	Language string // Application language.
}

// Data represents the JSON data uploaded in the request body.
type Data struct {
	CLab     CLabInfo    // Application info dictionary.
	UserInfo interface{} // Optional user info, we threat this as string.
}

// UserInfoAsString attemp to return a string rapresentation of UserInfo.
func (d *Data) UserInfoAsString() (string, error) {
	if v, ok := d.UserInfo.(string); ok {
		return v, nil
	} else if v, ok := d.UserInfo.(int); ok {
		return strconv.Itoa(v), nil
	} else {
		bytes, err := json.Marshal(d.UserInfo)
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}
}

// NewData returns a new Data.
func NewData() *Data {
	return &Data{}
}
