package apimodel

type LoginParam struct {
	User     string `form:"user" json:"user" binding:"required" example:"admin"`
	Password string `form:"password" json:"password" binding:"required" example:"123456"`
}
