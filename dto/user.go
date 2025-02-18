package dto

type UserSignInPayload struct {
	Email    string `valid:"required~Email can't be empty, email" example:"example@email.com" json:"email"`
	Password string `valid:"required~Password can't be empty" example:"secret" json:"password"`
}

type UserSignUpPayload struct {
	FullName string `valid:"required~Full name can't be empty" example:"Jhon Doe" json:"full_name"`
	Email    string `valid:"required~Email can't be empty, email" example:"example@email.com" json:"email"`
	Password string `valid:"required~Password can't be empty" example:"secret" json:"password"`
}

type UserModifyPayload struct {
	FullName string `valid:"required~Full name can't be empty" example:"Jhon Doe" json:"full_name"`
	Email    string `valid:"required~Email can't be empty, email" example:"example@email.com" json:"email"`
	Address  string `json:"address"`
}

type UserChangePasswordPayload struct {
	OldPassword        string `valid:"required~Old password can't be empty" example:"oldpassword" json:"old_password"`
	NewPassword        string `valid:"required~New password can't be empty" example:"newpassword" json:"new_password"`
	ConfirmNewPassword string `valid:"required~Confirm new password can't be empty" example:"newpassword" json:"confirm_new_password"`
}

type TokenString struct {
	Token string `json:"token"`
}
