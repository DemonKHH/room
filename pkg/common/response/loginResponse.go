package response

import "time"

// type LoginResponse struct {
// 	FirstName    *string   `json:"firstName" validate:"required,min=2,max=100"`
// 	LastName     *string   `json:"lastName" validate:"required"`
// 	Email        *string   `json:"email" validate:"email,required"`
// 	PhoneNumber  *string   `json:"phoneNumber" validate:"required"`
// 	Avator       *string   `json:"avator"`
// 	AccessToken  *string   `json:"accessToken"`
// 	RefreshToken *string   `json:"refreshToken"`
// 	CreatedAt    time.Time `json:"createdAt"`
// 	UpdatedAt    time.Time `json:"updatedAt"`
// 	UserId       string    `json:"userId"`
// }

type LoginResponse struct {
	FirstName    *string   `json:"firstName" validate:"required,min=2,max=100"`
	Email        *string   `json:"email" validate:"email,required"`
	Avator       *string   `json:"avator"`
	AccessToken  *string   `json:"accessToken"`
	RefreshToken *string   `json:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UserId       string    `json:"userId"`
}
