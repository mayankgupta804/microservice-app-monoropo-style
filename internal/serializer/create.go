package serializer

type CreateIncidentRequest struct {
	Message string `json:"message"`
}

type CreateIncidentResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}
