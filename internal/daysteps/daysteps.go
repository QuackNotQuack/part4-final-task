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
		return 0, 0, fmt.Errorf("ошибка: неверный формат данных")
	}

	// Проверка: шаги и длительность должны быть без пробелов
	if strings.Contains(parts[0], " ") {
		return 0, 0, fmt.Errorf("ошибка: количество шагов содержит пробелы")
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("ошибка преобразования количества шагов")
	}

	if strings.Contains(parts[1], " ") {
		return 0, 0, fmt.Errorf("ошибка: продолжительность содержит пробелы")
	}
	duration, err := time.ParseDuration(parts[1])
	if err != nil || duration <= 0 {
		return 0, 0, fmt.Errorf("ошибка преобразования времени")
	}

	return steps, duration, nil
}

// DayActionInfo анализирует строку с активностью за день, вычисляет дистанцию
// и количество калорий на основе веса и роста. Возвращает отформатированную строку.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("ошибка при парсинге: %v", err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("ошибка при расчёте калорий: %v", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
}
