package serializer

type Error struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
