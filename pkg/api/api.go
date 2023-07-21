package api

type Challenge struct {
	RandomString       string `json:"random_string"`
	NumberLeadingZeros int32  `json:"number_leading_zeros"`
	NumberSymbols      int32  `json:"number_symbols"`
}

type Solution struct {
	RandomString string `json:"random_string"`
}

type Quote struct {
	Quote string `json:"quote"`
}
