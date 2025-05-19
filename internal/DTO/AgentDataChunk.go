package DTO

type AgentDataChunk struct {
	AgentName string                 `json:"agent_name"`
	Data      map[string]interface{} `json:"data"`
	Duration  float64                `json:"duration"`
}
