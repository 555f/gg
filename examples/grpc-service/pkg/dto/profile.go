package dto

// Address Residential address
// Address Contract
type Address struct {
	// Street Street
	Street string `json:"street"`
	// Apt Apartment
	Apt  int   `json:"apt"`
	Apt2 int8  `json:"apt2"`
	Apt3 uint8 `json:"apt3"`
}

// Work Place of work
type Work struct {
	// Name
	Name string `json:"name"`
	// Age
	Age int `json:"age"`
	// Ration
	Ration float64 `json:"ration"`
}

// Profile User Profile
// Contract with profile data
type Profile struct {
	// ID profile
	ID int
	// FistName Surname
	FistName string `json:"fistName"`
	// LastName Name
	LastName string `json:"lastName"`
	// Address Residential address
	Address Address `json:"address"`
}
