package models

type RequestListsGeneral struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Total    int    `json:"total"`
	Order    string `json:"order"`
	OrderBy  string `json:"order_by"`
	Search   string `json:"search"`
	SearchBy string `json:"search_by"`
}

type ResponseGeneral struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
