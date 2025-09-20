package model

type Request struct {
	ListId   string `form:"listId"`
	Interval int    `form:"interval"`
}

type Response struct {
	Message string `json:"message"`
}

type DelRequest struct {
	ListId string `form:"listId"`
}
