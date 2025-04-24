package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining разбирает строку с тренировкой: шаги, вид активности и длительность.
// Возвращает шаги, тип активности, длительность и ошибку при необходимости.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid step count: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("step count must be positive")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration: %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be positive")
	}

	return steps, activity, duration, nil
}

// distance рассчитывает дистанцию в километрах по числу шагов и росту пользователя.
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

// meanSpeed рассчитывает среднюю скорость движения на основе шагов, роста и времени.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// RunningSpentCalories рассчитывает количество калорий, сожжённых при беге.
// Возвращает ошибку, если входные параметры некорректны.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input: all values must be positive")
	}
	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories рассчитывает количество калорий, сожжённых при ходьбе,
// с учётом корректирующего коэффициента. Возвращает ошибку при неверных данных.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input: all values must be positive")
	}
	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * speed * minutes) / minInH * walkingCaloriesCoefficient
	return calories, nil
}

// TrainingInfo возвращает текстовое описание тренировки на основе входной строки.
// Учитывает вид активности, длительность, дистанцию, скорость и сожжённые калории.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		err = fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), dist, speed, calories,
	), nil
}
