package models

import "encoding/json"

// SomeData key + value
type SomeData struct {
	Key string
	SomeDataValue
}

// SomeDataValue value define
type SomeDataValue struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
}

// MarshalBinary struct to json string
func (data *SomeDataValue) MarshalBinary() ([]byte, error) {
	return json.Marshal(data)
}

// UnmarshalBinary json string to struct
func (data *SomeDataValue) UnmarshalBinary(in []byte) error {
	// convert data to yours, let's assume its json data
	return json.Unmarshal(in, data)
}
