package models

// Country :
type Country struct {
	ID   int64             `json:"id"`
	Name NullToEmptyString `json:"name"`
	Iso  NullToEmptyString `json:"iso"`
}

// Countries :
type Countries []Country

// CountryAlias :
type CountryAlias Country
