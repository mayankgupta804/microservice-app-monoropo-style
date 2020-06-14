package domain

type Incident struct {
	ID      string
	Message string
	Status  string
	Ack     string
	Comment []string
}
