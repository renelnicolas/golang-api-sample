package models

// Company :
type Company struct {
	ID           int64             `json:"id"`
	Name         NullToEmptyString `json:"name"`
	Enabled      bool              `json:"enabled"`
	ExternalID   NullToEmptyString `json:"externalId"`
	ContactEmail NullToEmptyString `json:"contact_email,omitempty"`
	Address      NullToEmptyString `json:"address,omitempty"`
	ZipCode      NullToEmptyString `json:"zip_code,omitempty"`
	City         NullToEmptyString `json:"city,omitempty"`
	Phone        NullToEmptyString `json:"phone,omitempty"`
	VAT          NullToEmptyString `json:"vat,omitempty"`
	RCS          NullToEmptyString `json:"rcs,omitempty"`
	Website      NullToEmptyString `json:"website,omitempty"`
	Country      *Country          `json:"country,omitempty"`
}

// Companies :
type Companies []Company

// CompanyAlias :
type CompanyAlias Company
