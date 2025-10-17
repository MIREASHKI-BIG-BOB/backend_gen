package generator

import (
	"backend_gen/internal/ports/generator"
	"backend_gen/internal/ports/websocket"
	"log"
	"math"
	"math/rand"
	"time"
)

// ctgGenerator реализует генерацию CTG данных
// Может быть здоровый плод (60%) или с гипоксией (40%)
type ctgGenerator struct {
	rng *rand.Rand

	// Параметры состояния плода
	hasHypoxia bool // true = гипоксия, false = здоровый

	// Текущие значения
	currentBPM    float64
	currentUterus float64
}

// NewCTGGenerator создает новый CTG генератор с указанным режимом
// hypoxiaMode: 0 = здоровый плод, 1 = гипоксия
func NewCTGGenerator(hypoxiaMode int) generator.DataGenerator {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Определяем состояние плода на основе переменной окружения
	hasHypoxia := hypoxiaMode == 1

	var initialBPM, initialUterus float64
	if hasHypoxia {
		// Параметры гипоксии (из анализа CSV: BPM ~148, Uterus ~17)
		initialBPM = 148.0
		initialUterus = 17.0
		log.Println("🔴 CTG Generator: HYPOXIA mode (BPM: ~148, Uterus: ~17) [HYPOXIA_MODE=1]")
	} else {
		// Параметры здорового плода (из анализа CSV: BPM ~140, Uterus ~14.8)
		initialBPM = 140.0
		initialUterus = 14.5
		log.Println("🟢 CTG Generator: HEALTHY mode (BPM: ~140, Uterus: ~14.5) [HYPOXIA_MODE=0]")
	}

	return &ctgGenerator{
		rng:           rng,
		hasHypoxia:    hasHypoxia,
		currentBPM:    initialBPM,
		currentUterus: initialUterus,
	}
}

// GenerateNext генерирует следующую точку данных
func (g *ctgGenerator) GenerateNext(timestamp float64) websocket.SensorData {
	g.update()

	// Генерируем Spasms для здорового плода (без схваток)
	// Небольшие естественные колебания в покое
	spasms := 20.0 + g.rng.Float64()*2 - 1.0 // 19-21, небольшие колебания

	return websocket.SensorData{
		BPMChild: g.currentBPM,
		Uterus:   g.currentUterus,
		Spasms:   math.Max(0, spasms),
	}
}

// update обновляет состояние генератора
func (g *ctgGenerator) update() {
	if g.hasHypoxia {
		// === ГИПОКСИЯ: Повышенный тонус матки и нестабильный пульс ===
		g.currentUterus += g.rng.Float64()*0.8 - 0.4 // Больше колебаний

		// Диапазон для гипоксии (из анализа: среднее 17.09)
		if g.currentUterus < 16.0 {
			g.currentUterus = 16.0
		}
		if g.currentUterus > 18.5 {
			g.currentUterus = 18.5
		}

		// Пульс при гипоксии: выше и с большей вариабельностью
		baseBPM := 148.0                   // Из анализа: среднее 148.25
		varDelta := g.rng.Float64()*12 - 6 // Диапазон ±6 (больше вариабельность)

		g.currentBPM = baseBPM + varDelta

		// Пределы для гипоксии
		if g.currentBPM < 140 {
			g.currentBPM = 140
		}
		if g.currentBPM > 156 {
			g.currentBPM = 156
		}

	} else {
		// === ЗДОРОВЫЙ ПЛОД: Низкий тонус и стабильный пульс ===
		g.currentUterus += g.rng.Float64()*0.6 - 0.3 // Небольшой дрейф

		// Диапазон для здорового плода (из анализа: среднее 14.80)
		if g.currentUterus < 13.5 {
			g.currentUterus = 13.5
		}
		if g.currentUterus > 15.5 {
			g.currentUterus = 15.5
		}

		// Пульс здорового плода
		baseBPM := 140.0                    // Из анализа: среднее 139.85
		varDelta := g.rng.Float64()*9 - 4.5 // Диапазон ±4.5 (норма)

		g.currentBPM = baseBPM + varDelta

		// Пределы для здорового
		if g.currentBPM < 135 {
			g.currentBPM = 135
		}
		if g.currentBPM > 145 {
			g.currentBPM = 145
		}
	}
}

// Reset сбрасывает генератор в начальное состояние (сохраняет текущий режим)
func (g *ctgGenerator) Reset() {
	g.rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// НЕ меняем режим! Сохраняем g.hasHypoxia как есть
	if g.hasHypoxia {
		g.currentBPM = 148.0
		g.currentUterus = 17.0
		log.Println("🔄 CTG Generator RESET: HYPOXIA mode")
	} else {
		g.currentBPM = 140.0
		g.currentUterus = 14.5
		log.Println("🔄 CTG Generator RESET: HEALTHY mode")
	}
}

// SetParameters устанавливает параметры генерации (для совместимости с интерфейсом)
func (g *ctgGenerator) SetParameters(params generator.GenerationParameters) {
	// Генератор настроен на здоровый плод с фиксированными параметрами
	// Этот метод оставлен для совместимости с интерфейсом
}
