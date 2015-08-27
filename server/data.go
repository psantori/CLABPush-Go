package server

// CLabInfo represents the CLab subdictionary of the JSON body.
type CLabInfo struct {
	AppID    string // Application package.
	Language string // Application language.
}

// Data represents the JSON data uploaded in the request body.
type Data struct {
	CLab     CLabInfo // Application info dictionary.
	UserInfo string   // Optional user info, we threat this as string.
}

// NewData returns a new Data.
func NewData() *Data {
	return &Data{}
}
