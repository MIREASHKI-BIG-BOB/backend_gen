package generator

import (
	"backend_gen/internal/ports/generator"
	"backend_gen/internal/ports/websocket"
	"math"
	"math/rand"
	"time"
)

// TODO: проверить на соответствие реальным данным (тут все сгенерирровано нейронкой)

// sinusoidalGenerator реализует генератор на основе синусоидальных функций
type sinusoidalGenerator struct {
	params generator.GenerationParameters
	rng    *rand.Rand
}

// NewSinusoidalGenerator создает новый генератор на основе синусов
func NewSinusoidalGenerator() generator.DataGenerator {
	return &sinusoidalGenerator{
		params: generator.GenerationParameters{
			// Параметры по умолчанию для BPM (сердцебиение) - реалистичные значения
			BPMBase:      80.0, // Нормальный пульс в покое
			BPMAmplitude: 20.0, // Колебания ±20 (60-100 BPM)
			BPMFrequency: 0.5,  // Hz

			// Параметры по умолчанию для Uterus
			UterusBase:      67.0, // Увеличено с 7.0 до 67.0, чтобы избежать отрицательных значений
			UterusAmplitude: 60.0,
			UterusFrequency: 0.3, // Hz

			// Параметры по умолчанию для Spasms
			SpasmsBase:      45.0, // Базовое значение схваток
			SpasmsAmplitude: 35.0, // Амплитуда схваток
			SpasmsFrequency: 0.4,  // Hz

			NoiseLevel: 0.05, // 5% шума
		},
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateNext генерирует следующую точку данных
func (g *sinusoidalGenerator) GenerateNext(timestamp float64) websocket.SensorData {
	// Генерируем BPM с использованием синусоиды + шум
	bpmSin := math.Sin(2 * math.Pi * g.params.BPMFrequency * timestamp)
	bpmNoise := g.generateNoise() * g.params.NoiseLevel * g.params.BPMAmplitude
	bpm := g.params.BPMBase + g.params.BPMAmplitude*bpmSin + bpmNoise

	// Генерируем Uterus с использованием косинуса + шум (смещение фазы для разнообразия)
	uterusCos := math.Cos(2 * math.Pi * g.params.UterusFrequency * timestamp)
	uterusNoise := g.generateNoise() * g.params.NoiseLevel * g.params.UterusAmplitude
	uterus := g.params.UterusBase + g.params.UterusAmplitude*uterusCos + uterusNoise

	// Добавляем небольшие скачки как в реальных данных
	if g.rng.Float64() < 0.02 { // 2% вероятность скачка
		bpm += g.generateNoise() * g.params.BPMAmplitude * 0.3
	}

	if g.rng.Float64() < 0.01 { // 1% вероятность скачка для uterus
		uterus += g.generateNoise() * g.params.UterusAmplitude * 0.5
	}

	spasmsSin := math.Sin(2 * math.Pi * g.params.SpasmsFrequency * timestamp)
	spasmsNoise := g.generateNoise() * g.params.NoiseLevel * g.params.SpasmsAmplitude
	spasms := g.params.SpasmsBase + g.params.SpasmsAmplitude*spasmsSin + spasmsNoise

	if g.rng.Float64() < 0.015 { // 1.5% вероятность скачка для spasms
		spasms += g.generateNoise() * g.params.SpasmsAmplitude * 0.4
	}

	return websocket.SensorData{
		BPM:    math.Max(0, bpm),    // Обеспечиваем неотрицательность
		Uterus: math.Max(0, uterus), // Обеспечиваем неотрицательность
		Spasms: math.Max(0, spasms), // Обеспечиваем неотрицательность
	}
}

// Reset сбрасывает генератор
func (g *sinusoidalGenerator) Reset() {
	g.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// SetParameters устанавливает новые параметры генерации
func (g *sinusoidalGenerator) SetParameters(params generator.GenerationParameters) {
	g.params = params
}

// generateNoise генерирует случайный шум в диапазоне [-1, 1]
func (g *sinusoidalGenerator) generateNoise() float64 {
	return (g.rng.Float64() - 0.5) * 2.0
}
