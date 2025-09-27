package generator

import "backend_gen/internal/ports/websocket"

// DataGenerator интерфейс для генерации медицинских данных
type DataGenerator interface {
	// GenerateNext генерирует следующую точку данных на основе времени
	GenerateNext(timestamp float64) websocket.SensorData

	// Reset сбрасывает генератор в начальное состояние
	Reset()

	// SetParameters устанавливает параметры генерации
	SetParameters(params GenerationParameters)
}

// GenerationParameters параметры для генерации данных
type GenerationParameters struct {
	// BPM параметры
	BPMBase      float64 // Базовое значение BPM
	BPMAmplitude float64 // Амплитуда колебаний BPM
	BPMFrequency float64 // Частота колебаний BPM

	// Uterus параметры
	UterusBase      float64 // Базовое значение Uterus
	UterusAmplitude float64 // Амплитуда колебаний Uterus
	UterusFrequency float64 // Частота колебаний Uterus

	// Общие параметры
	NoiseLevel float64 // Уровень шума (0-1)
}
