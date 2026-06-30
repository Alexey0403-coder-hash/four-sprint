package spentcalories

import (
	"errors"
	"fmt"
	"log"
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

// parseTraining принимает строку и возвращает количество шагов, вид активности, продолжительность активности и ошибку.
func parseTraining(data string) (int, string, time.Duration, error) {
	part := strings.Split(data, ",") //Разбиваем строку по запятой на три части

	//Проверяем, что получилось три части.
	if len(part) != 3 {
		return 0, "", 0, errors.New("неверный формат: ожидается строка вида 'число строка время'")
	}

	//Количество шагов типа int.
	step, err := strconv.Atoi(part[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга числа: %w", err)
	}

	//Длительность ходьбы типа time.Duration.
	duration, err := time.ParseDuration(part[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга времени: %w", err)
	}

	return step, part[1], duration, nil
}

// distance принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {

	//Длина шага.
	stepLen := height * stepLengthCoefficient

	//Пройденная дистанция в метрах.
	distM := float64(steps) * stepLen

	//Пройденная дистанция в километрах.
	distKm := distM / mInKm

	return distKm
}

// meanSpeed функция принимает количество шагов, рост пользователя и продолжительность активности и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	//Продолжительность активности в часах.
	hours := duration.Hours()

	//Средняя скорость.
	midSpeed := distance(steps, height) / hours

	return midSpeed
}

// TrainingInfo принимает параметры тренировки и возвращает их в нужном формате или возвращает ошибку.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	hours := duration.Hours()

	var dist, speed, calories float64

	//Проверяем вид тренировки и рассчитываем ее параметры.
	switch trainType {
	case "Ходьба":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Бег":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainType, hours, dist, speed, calories), nil
}

// RunningSpentCalories принимает количество шагов, вес и рост и продолжительность бега, а возвращает количество потраченных каллорий и ошибку.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неверные данные: параметры должный быть больше нуля")
	}

	//Средняя скорость при беге.
	speed := meanSpeed(steps, height, duration)

	//Продолжительность бега в минутах.
	minutes := duration.Minutes()

	//Количество потраченных калорий при беге.
	trainCalories := (weight * speed * minutes) / minInH

	return trainCalories, nil
}

// WalkingSpentCalories принимает количество шагов, вес и рост и продолжительность ходьбы, а возвращает количество потраченных каллорий и ошибку.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неверные данные: параметры должный быть больше нуля")
	}

	//Средняя скорость при ходьбе.
	speed := meanSpeed(steps, height, duration)

	//Продолжительность ходьбы в минутах.
	minutes := duration.Minutes()

	//Количество потраченных калорий при ходьбе.
	walkCalories := ((weight * speed * minutes) / minInH) * walkingCaloriesCoefficient

	return walkCalories, nil
}
