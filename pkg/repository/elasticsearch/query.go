package elasticsearch

type ReadLogsQuery struct {
	Query Query `json:"query"`
}

type Query struct {
	Match map[string]interface{} `json:"match"`
}

type SearchResponse struct {
	Timestamp string `json:"@timestamp"`
	Message   string `json:"message"`
}
