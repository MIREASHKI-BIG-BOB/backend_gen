package websocket

type MessageData struct {
	Timestamp float64    `json:"timestamp"`
	Data      SensorData `json:"data"`
}

// показания сенсоров
type SensorData struct {
	BPM    float64 `json:"bpm"`
	Uterus float64 `json:"uterus"`
}
