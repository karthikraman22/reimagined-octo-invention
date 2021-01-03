package model

// Model represents ...
type Model struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

// GetModelByID ...
func GetModelByID(id string) (Model, error) {
	return Model{}, nil
}
