package serializer

type DeleteIncidentRequest struct {
	ID string `json:"id"`
}

type DeleteIncidentResponse struct {
	Status string `json:"status"`
}
