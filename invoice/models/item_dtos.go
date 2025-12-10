package invoice

type SaveItemDTO struct {
	Id       string  `json:"id,omitempty" dynamodbav:"id,omitempty" validate:"required,uuid"`
	Name     string  `json:"name,omitempty" dynamodbav:"name,omitempty" validate:"required,min=1,max=100"`
	Quantity int     `json:"quantity,omitempty" dynamodbav:"quantity,omitempty" validate:"required,min=1"`
	Price    float64 `json:"price,omitempty" dynamodbav:"price,omitempty" validate:"required,gt=0"`
}
