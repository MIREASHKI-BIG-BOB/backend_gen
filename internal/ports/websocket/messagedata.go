package websocket

type MessageData struct {
	SensorID     string     `json:"sensorID"`
	SecFromStart float64    `json:"secFromStart"`
	Data         SensorData `json:"data"`
}

type SensorData struct {
	BPM    float64 `json:"bpm"`
	Uterus float64 `json:"uterus"`
	Spasms float64 `json:"spasms"`
}
