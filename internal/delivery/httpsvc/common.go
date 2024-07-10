package httpsvc

type response struct {
	Status   string         `json:"status"` // success/failed
	Metadata map[string]any `json:"metadata"`
	Data     any            `json:"data"`
}
