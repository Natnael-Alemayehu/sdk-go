package models

// Parameter represents a key-value parameter
type Parameter struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Reference represents reference data
type Reference struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Initiator represents the transaction initiator
type Initiator struct {
	IdentifierType     int    `json:"IdentifierType"`
	Identifier         string `json:"Identifier"`
	SecurityCredential string `json:"SecurityCredential"`
	SecretKey          string `json:"SecretKey"`
}

// Party represents a transaction party
type Party struct {
	IdentifierType int    `json:"IdentifierType"`
	Identifier     string `json:"Identifier"`
}

// ReceiverParty represents the receiving party
type ReceiverParty struct {
	IdentifierType int    `json:"IdentifierType"`
	Identifier     string `json:"Identifier"`
	ShortCode      string `json:"ShortCode"`
}
