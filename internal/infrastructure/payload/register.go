package payload

type RegisterRequest struct {
	Email     string  `json:"email" binding:"required"`
	Phone     *string `json:"phone"`
	FullName  string  `json:"full_name" binding:"required"`
	Role      int     `json:"role" binding:"required,oneof=1 2 3"`
	Password  string  `json:"password" binding:"required"`
	Avatar    *string `json:"avatar"`
	Birthdate *string `json:"birthdate"`
}
