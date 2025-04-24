package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage разбирает строку с количеством шагов и длительностью,
// возвращая шаги и продолжительность прогулки. В случае ошибки — возвращает её.
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid data format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid step count: %v", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("step count must be positive")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration: %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("duration must be positive")
	}
	return steps, duration, nil
}

// DayActionInfo анализирует строку с активностью за день, вычисляет дистанцию
// и количество калорий на основе веса и роста. Возвращает отформатированную строку.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("parsing error: %v", err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("calorie calculation error: %v", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
}
