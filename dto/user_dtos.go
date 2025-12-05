package dto

type User struct {
	//@Description User First Name. Max length 50, min length 2
	Firstname string `json:"firstName" validate:"required,max=50,min=2"`
	//@Description User Last Name. Max length 50, min length 2
	Lastname string `json:"lastName" validate:"required,max=50,min=2"`
	//@Description User email
	Email string `json:"email" validate:"required,email"`
	//@Description User phone. Optional
	Phone string `json:"phone" validate:"e164"`
	//@Description User age. Optional
	Age int32 `json:"age" validate:"numeric,gt=0"`
	//@Description User status. Optional
	Status string `json:"status"`
}
