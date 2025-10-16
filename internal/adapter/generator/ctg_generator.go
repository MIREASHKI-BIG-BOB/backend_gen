package generator

import (
	"backend_gen/internal/ports/generator"
	"backend_gen/internal/ports/websocket"
	"math"
	"math/rand"
	"time"
)

// ctgGenerator реализует реалистичную генерацию CTG данных
// с эмуляцией схваток, гипоксии, вариабельности и децелераций
type ctgGenerator struct {
	timeMS int64
	rng    *rand.Rand

	// Состояние схватки
	inContraction bool
	phase         float64
	duration      float64
	intensity     float64

	// Состояние гипоксии
	hypoxia        bool
	hypoxiaStarted bool
	hypoxiaStep    int

	// Текущие значения
	currentBPM    float64
	currentUterus float64
}

// NewCTGGenerator создает новый CTG генератор
func NewCTGGenerator() generator.DataGenerator {
	return &ctgGenerator{
		timeMS:        0,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		currentBPM:    135,
		currentUterus: 17,
	}
}

// GenerateNext генерирует следующую точку данных
func (g *ctgGenerator) GenerateNext(timestamp float64) websocket.SensorData {
	g.update()

	// Генерируем Spasms на основе схваток
	// Когда идет схватка (Uterus > 28), spasms увеличивается
	spasms := 20.0
	if g.currentUterus > 28 {
		// Во время схватки spasms растет пропорционально интенсивности
		spasms = 20.0 + (g.currentUterus-28)*1.5
	} else {
		// В покое - небольшие колебания
		spasms = 20.0 + g.rng.Float64()*3 - 1.5
	}

	return websocket.SensorData{
		BPMChild: g.currentBPM,
		Uterus:   g.currentUterus,
		Spasms:   math.Max(0, spasms),
	}
}

// update обновляет состояние генератора (эквивалент Update() из realgen.go)
func (g *ctgGenerator) update() {
	g.timeMS += 120

	// Инициализация гипоксии в 50% случаев (один раз)
	if !g.hypoxiaStarted {
		g.hypoxia = g.rng.Float64() < 0.5
		g.hypoxiaStarted = true
	}

	// === Эмуляция схваток ===
	if !g.inContraction && g.rng.Float64() < 0.0008 {
		g.inContraction = true
		g.phase = 0
		g.duration = 150 + g.rng.Float64()*100 // 150–250 сек
		g.intensity = 30 + g.rng.Float64()*6
		if g.hypoxia {
			// более частые и затяжные схватки при гипоксии
			g.duration += 30
			g.intensity += 2
		}
	}

	if g.inContraction {
		step := 0.12 / g.duration
		g.phase += step

		x := g.phase
		if x > 1 {
			g.inContraction = false
			g.currentUterus = 17 + g.rng.Float64()*1.5
		} else {
			shape := math.Sin(math.Pi * math.Pow(x, 0.6))
			g.currentUterus = 17 + shape*(g.intensity-17)
		}
	} else {
		g.currentUterus += g.rng.Float64()*0.4 - 0.2
		if g.currentUterus < 15 {
			g.currentUterus = 15
		}
		if g.currentUterus > 19 {
			g.currentUterus = 19
		}
	}

	// === Пульс (BPM) ===
	baseBPM := 130.0

	// гипоксия прогрессирует
	if g.hypoxia {
		g.hypoxiaStep++
		switch {
		case g.hypoxiaStep < 2500:
			baseBPM += 15 // начальная тахикардия
		case g.hypoxiaStep < 5000:
			baseBPM -= 10 // спад
		default:
			baseBPM -= 25 // устойчивая брадикардия
			if baseBPM < 95 {
				baseBPM = 95
			}
		}
	}

	// вариабельность
	varDelta := 5.0
	if g.hypoxia && g.hypoxiaStep > 3000 {
		varDelta = g.rng.Float64()*2 + 0.5 // падение вариабельности
	} else {
		varDelta = g.rng.Float64()*10 - 5
	}

	g.currentBPM = baseBPM + varDelta

	// === Децелерации при схватке и гипоксии ===
	if g.currentUterus > 28 {
		if g.hypoxia && g.hypoxiaStep > 2000 {
			g.currentBPM -= float64(g.rng.Intn(30) + 10) // поздние глубокие децелерации
		} else {
			g.currentBPM -= float64(g.rng.Intn(10) + 5)
		}
	}

	if g.currentBPM < 85 {
		g.currentBPM = 85
	}
}

// Reset сбрасывает генератор в начальное состояние
func (g *ctgGenerator) Reset() {
	g.timeMS = 0
	g.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	g.inContraction = false
	g.phase = 0
	g.duration = 0
	g.intensity = 0
	g.hypoxia = false
	g.hypoxiaStarted = false
	g.hypoxiaStep = 0
	g.currentBPM = 135
	g.currentUterus = 17
}

// SetParameters устанавливает параметры генерации (для совместимости с интерфейсом)
func (g *ctgGenerator) SetParameters(params generator.GenerationParameters) {
	// CTG генератор использует фиксированные параметры из медицинской практики
	// Этот метод оставлен для совместимости с интерфейсом, но не используется
}

