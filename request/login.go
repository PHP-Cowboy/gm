package request

type Login struct {
	Id       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha" binding:"required"`
}

type LoginOut struct {
	Id int `json:"id"`
}

type CheckGoogleCaptcha struct {
	Uid  int    `json:"uid"`
	Code string `json:"code"`
}
