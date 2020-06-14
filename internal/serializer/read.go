package serializer

type ReadIncidentRequest struct {
	ID string `json:"id"`
}

type ReadIncidentResponse struct {
	Status   string   `json:"status"`
	Ack      string   `json:"acknowledged"`
	Message  string   `json:"message"`
	Comments []string `json:"comments,omitempty"`
}
