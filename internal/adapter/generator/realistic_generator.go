package generator

import (
	"backend_gen/internal/ports/generator"
	"backend_gen/internal/ports/websocket"
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// RealDataPoint представляет точку реальных медицинских данных
type RealDataPoint struct {
	Time  float64
	Value float64
}

// DataPattern содержит загруженные реальные данные для использования в генерации
type DataPattern struct {
	BPMData    []RealDataPoint
	UterusData []RealDataPoint
	// Статистика для генерации разнообразия
	BPMStats    DataStats
	UterusStats DataStats
}

// DataStats содержит статистические характеристики данных
type DataStats struct {
	Min    float64
	Max    float64
	Mean   float64
	StdDev float64
	Count  int
}

// realisticGenerator генератор на основе реальных медицинских данных
type realisticGenerator struct {
	regularPatterns []DataPattern
	hypoxiaPatterns []DataPattern
	currentPattern  *DataPattern
	isHypoxia       bool
	startTime       time.Time
	rng             *rand.Rand
	params          generator.GenerationParameters
}

// NewRealisticGenerator создает новый генератор на основе реальных данных
func NewRealisticGenerator() generator.DataGenerator {
	g := &realisticGenerator{
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		startTime: time.Now(),
	}
	
	// Загружаем реальные данные при инициализации
	if err := g.loadRealData(); err != nil {
		// Если не удалось загрузить данные, используем параметры по умолчанию
		g.params = generator.GenerationParameters{
			BPMBase:      140.0,
			BPMAmplitude: 20.0,
			BPMFrequency: 0.02, // Более низкая частота для реалистичности
			
			UterusBase:      15.0,
			UterusAmplitude: 10.0,
			UterusFrequency: 0.01,
			
			SpasmsBase:      20.0,
			SpasmsAmplitude: 15.0,
			SpasmsFrequency: 0.015,
			
			NoiseLevel: 0.05,
		}
	}
	
	return g
}

// loadRealData загружает реальные медицинские данные из CSV файлов
func (g *realisticGenerator) loadRealData() error {
	// Загружаем образцы данных из regular и hypoxia папок
	regularDir := "regular"
	hypoxiaDir := "hypoxia"
	
	// Загружаем несколько паттернов из каждой категории
	regularPatterns, err := g.loadPatternsFromDir(regularDir, 5) // загружаем 5 паттернов
	if err != nil {
		return fmt.Errorf("failed to load regular patterns: %w", err)
	}
	g.regularPatterns = regularPatterns
	
	hypoxiaPatterns, err := g.loadPatternsFromDir(hypoxiaDir, 5)
	if err != nil {
		return fmt.Errorf("failed to load hypoxia patterns: %w", err)
	}
	g.hypoxiaPatterns = hypoxiaPatterns
	
	// Выбираем случайный паттерн для начала (70% шанс на regular, 30% на hypoxia)
	g.switchPattern()
	
	return nil
}

// loadPatternsFromDir загружает паттерны данных из указанной директории
func (g *realisticGenerator) loadPatternsFromDir(dir string, maxPatterns int) ([]DataPattern, error) {
	var patterns []DataPattern
	
	// Получаем список подпапок (номера случаев)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() || count >= maxPatterns {
			continue
		}
		
		caseDir := filepath.Join(dir, entry.Name())
		pattern, err := g.loadPattern(caseDir)
		if err != nil {
			continue // пропускаем поврежденные данные
		}
		
		patterns = append(patterns, pattern)
		count++
	}
	
	return patterns, nil
}

// loadPattern загружает один паттерн данных (BPM и Uterus) из папки случая
func (g *realisticGenerator) loadPattern(caseDir string) (DataPattern, error) {
	var pattern DataPattern
	
	// Загружаем BPM данные
	bpmDir := filepath.Join(caseDir, "bpm")
	bpmFiles, err := filepath.Glob(filepath.Join(bpmDir, "*.csv"))
	if err != nil || len(bpmFiles) == 0 {
		return pattern, fmt.Errorf("no BPM files found")
	}
	
	// Берем первый файл BPM
	bpmData, stats, err := g.loadCSVData(bpmFiles[0])
	if err != nil {
		return pattern, err
	}
	pattern.BPMData = bpmData
	pattern.BPMStats = stats
	
	// Загружаем Uterus данные
	uterusDir := filepath.Join(caseDir, "uterus")
	uterusFiles, err := filepath.Glob(filepath.Join(uterusDir, "*.csv"))
	if err != nil || len(uterusFiles) == 0 {
		return pattern, fmt.Errorf("no Uterus files found")
	}
	
	// Берем первый файл Uterus
	uterusData, stats, err := g.loadCSVData(uterusFiles[0])
	if err != nil {
		return pattern, err
	}
	pattern.UterusData = uterusData
	pattern.UterusStats = stats
	
	return pattern, nil
}

// loadCSVData загружает данные из CSV файла
func (g *realisticGenerator) loadCSVData(filename string) ([]RealDataPoint, DataStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, DataStats{}, err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, DataStats{}, err
	}
	
	if len(records) < 2 { // header + at least one data row
		return nil, DataStats{}, fmt.Errorf("insufficient data")
	}
	
	var data []RealDataPoint
	var values []float64
	
	// Пропускаем заголовок
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 2 {
			continue
		}
		
		time, err := strconv.ParseFloat(records[i][0], 64)
		if err != nil {
			continue
		}
		
		value, err := strconv.ParseFloat(records[i][1], 64)
		if err != nil {
			continue
		}
		
		data = append(data, RealDataPoint{Time: time, Value: value})
		values = append(values, value)
	}
	
	// Вычисляем статистику
	stats := g.calculateStats(values)
	
	return data, stats, nil
}

// calculateStats вычисляет статистические характеристики данных
func (g *realisticGenerator) calculateStats(values []float64) DataStats {
	if len(values) == 0 {
		return DataStats{}
	}
	
	var stats DataStats
	stats.Count = len(values)
	
	// Min/Max/Sum
	min := values[0]
	max := values[0]
	sum := 0.0
	
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	
	stats.Min = min
	stats.Max = max
	stats.Mean = sum / float64(len(values))
	
	// Standard deviation
	variance := 0.0
	for _, v := range values {
		diff := v - stats.Mean
		variance += diff * diff
	}
	stats.StdDev = math.Sqrt(variance / float64(len(values)))
	
	return stats
}

// switchPattern переключается на другой паттерн данных
func (g *realisticGenerator) switchPattern() {
	// 70% шанс на regular, 30% на hypoxia
	if g.rng.Float64() < 0.7 && len(g.regularPatterns) > 0 {
		g.isHypoxia = false
		idx := g.rng.Intn(len(g.regularPatterns))
		g.currentPattern = &g.regularPatterns[idx]
	} else if len(g.hypoxiaPatterns) > 0 {
		g.isHypoxia = true
		idx := g.rng.Intn(len(g.hypoxiaPatterns))
		g.currentPattern = &g.hypoxiaPatterns[idx]
	}
}

// GenerateNext генерирует следующую точку данных на основе реальных паттернов
func (g *realisticGenerator) GenerateNext(timestamp float64) websocket.SensorData {
	// Иногда переключаемся на другой паттерн для разнообразия
	if g.rng.Float64() < 0.001 { // 0.1% шанс переключения на каждом такте
		g.switchPattern()
	}
	
	var bpm, uterus float64
	
	if g.currentPattern != nil {
		bpm = g.generateFromPattern(g.currentPattern.BPMData, g.currentPattern.BPMStats, timestamp)
		uterus = g.generateFromPattern(g.currentPattern.UterusData, g.currentPattern.UterusStats, timestamp)
	} else {
		// Fallback на синусоидальную генерацию если данные не загружены
		bpm = g.generateSinusoidal(timestamp, g.params.BPMBase, g.params.BPMAmplitude, g.params.BPMFrequency)
		uterus = g.generateSinusoidal(timestamp, g.params.UterusBase, g.params.UterusAmplitude, g.params.UterusFrequency)
	}
	
	// Генерация spasms (пока используем простую модель)
	spasms := g.generateSinusoidal(timestamp, g.params.SpasmsBase, g.params.SpasmsAmplitude, g.params.SpasmsFrequency)
	
	// Добавляем немного шума для реалистичности
	bpm += g.generateNoise() * 2.0
	uterus += g.generateNoise() * 1.0
	spasms += g.generateNoise() * 1.5
	
	return websocket.SensorData{
		BPMChild: math.Max(0, bpm),
		Uterus:   math.Max(0, uterus),
		Spasms:   math.Max(0, spasms),
	}
}

// generateFromPattern генерирует значение на основе реального паттерна данных
func (g *realisticGenerator) generateFromPattern(data []RealDataPoint, stats DataStats, timestamp float64) float64 {
	if len(data) == 0 {
		return stats.Mean
	}
	
	// Используем циклическое воспроизведение данных
	dataLen := float64(len(data))
	normalizedTime := math.Mod(timestamp*10, dataLen) // масштабируем время
	
	// Линейная интерполяция между точками
	idx := int(normalizedTime)
	if idx >= len(data)-1 {
		return data[len(data)-1].Value
	}
	
	// Интерполируем между текущей и следующей точкой
	t := normalizedTime - float64(idx)
	value := data[idx].Value*(1-t) + data[idx+1].Value*t
	
	// Добавляем небольшие вариации на основе статистики
	variation := g.generateNoise() * stats.StdDev * 0.1
	
	return value + variation
}

// generateSinusoidal генерирует синусоидальное значение (fallback)
func (g *realisticGenerator) generateSinusoidal(timestamp, base, amplitude, frequency float64) float64 {
	sin := math.Sin(2 * math.Pi * frequency * timestamp)
	noise := g.generateNoise() * g.params.NoiseLevel * amplitude
	return base + amplitude*sin + noise
}

// generateNoise генерирует случайный шум в диапазоне [-1, 1]
func (g *realisticGenerator) generateNoise() float64 {
	return (g.rng.Float64() - 0.5) * 2.0
}

// Reset сбрасывает генератор
func (g *realisticGenerator) Reset() {
	g.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	g.startTime = time.Now()
	g.switchPattern() // Выбираем новый паттерн
}

// SetParameters устанавливает параметры генерации
func (g *realisticGenerator) SetParameters(params generator.GenerationParameters) {
	g.params = params
}
