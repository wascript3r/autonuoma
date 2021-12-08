package faq

// GetAll

type FAQListInfo struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"categoryID"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
}

type GetAllRes struct {
	FAQ []*FAQListInfo `json:"faq"`
}
