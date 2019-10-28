package dcc

// Config allows to store configuration settings to initialize go-dcc.
type Config struct {
	Locomotives []Locomotive `json:"locomotives"`
}
