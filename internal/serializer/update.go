package serializer

type UpdateIncidentRequest struct {
	Message          string `json:"message"`
	Ack              string `json:"acknowledged"`
	ResolutionStatus string `json:"status"`
	Comment          string `json:"comment"`
}

type UpdateIncidentResponse struct {
	Status string `json:"status"`
}
