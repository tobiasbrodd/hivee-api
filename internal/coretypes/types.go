package coretypes

type Measure struct {
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}
