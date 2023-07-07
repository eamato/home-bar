package domain

type PaginationRequest struct {
	Take int `form:"take" binding:"max=500"`
	Skip int `form:"skip"`
}

type PaginationResponse struct {
	Total int64 `json:"total"`
}
