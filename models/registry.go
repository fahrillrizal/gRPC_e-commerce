package models

var RegisteredModels []interface{}

func RegisterModel(m interface{}) {
	RegisteredModels = append(RegisteredModels, m)
}