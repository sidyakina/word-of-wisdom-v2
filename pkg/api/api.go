package api

type Challenge struct {
	RandomString       string `json:"random_string"`
	NumberLeadingZeros int32  `json:"number_leading_zeros"`
	NumberSymbols      int32  `json:"number_symbols"`
}

type Solution struct {
	Solution string `json:"solution"`
}

type Quote struct {
	Quote string `json:"quote"`
}
