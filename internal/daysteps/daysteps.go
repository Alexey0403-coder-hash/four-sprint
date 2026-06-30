package daysteps

import (
	"errors"
	"fmt"
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

// parsePackage принимает строку и возвращает количество шагов, продолжительность активности и ошибку.
func parsePackage(data string) (int, time.Duration, error) {
	part := strings.Split(data, ",") //Разбиваем строку по запятой на две части

	//Проверяем, что получилось две части.
	if len(part) != 2 {
		return 0, 0, errors.New("неверный формат: ожидается строка вида 'число время'")
	}

	//Количество шагов типа int.
	step, err := strconv.Atoi(part[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга числа: %w", err)
	}

	//Длительность ходьбы типа time.Duration.
	duration, err := time.ParseDuration(part[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга времени: %w", err)
	}

	return step, duration, nil
}

// DayActionInfo принимает количество шагов, продолжительность ходьбы, вес и рост и возвращает количество шагов, дистанцию и количество калорий в нужном формате.
func DayActionInfo(data string, weight, height float64) string {

	//Получаем данные о количестве шагов, продолжительности прогулки и наличии ошибки
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	//проверяем, что количество шагов больше 0
	if steps <= 0 {
		return ""
	}

	//Пройденная дистанция в метрах.
	distM := float64(steps) * stepLength

	//Пройденная дистанция в километрах.
	distKm := distM / float64(mInKm)

	//Количество калорий, потраченных при ходьбе.
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distKm, calories)
}
