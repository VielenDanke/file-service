package validations

// Validator ...
type Validator interface {
	ValidateMap(body map[string]interface{}) error
	ValidateStruct(body interface{}) error
}
