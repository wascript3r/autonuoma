package review

// Submit

type SubmitReq struct {
	TicketID int     `json:"ticketID" validate:"required"`
	Stars    int     `json:"stars" validate:"required,r_stars"`
	Comment  *string `json:"comment" validate:"omitempty,required,r_comment"`
}

type ReviewInfo struct {
	Stars   int     `json:"stars"`
	Comment *string `json:"comment"`
}
