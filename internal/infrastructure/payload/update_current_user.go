package payload

type UpdateCurrentUserRequest struct {
	FullName  *string `json:"full_name"`
	Phone     *string `json:"phone"`
	Avatar    *string `json:"avatar"`
	Birthdate *string `json:"birthdate"`
}
