package model

type Request struct {
	UserName string `form:"userName"`
	UserId   string `form:"userId"`
	Interval int    `form:"interval"`
}

type Response struct {
	Message string `json:"message"`
}

type DelRequest struct {
	UserName string `form:"userName"`
}
