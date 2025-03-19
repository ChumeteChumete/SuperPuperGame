// Пакет utils содержит вспомогательные функции
package utils

import (
	"math"
	"math/rand"
)

// Clamp ограничивает значение в диапазоне [min, max]
func Clamp(value, min, max float64) float64 {
	// Возвращаем значение, не меньшее min и не большее max
	return math.Max(min, math.Min(max, value))
}

// Distance вычисляет расстояние между двумя точками
func Distance(x1, y1, x2, y2 float64) float64 {
	// Вычисляем разницу координат
	dx := x2 - x1
	dy := y2 - y1
	
	// Возвращаем расстояние по теореме Пифагора
	return math.Sqrt(dx*dx + dy*dy)
}

// RandomRange возвращает случайное число в диапазоне [min, max)
func RandomRange(min, max float64) float64 {
	// Возвращаем случайное число в указанном диапазоне
	return min + rand.Float64()*(max-min)
}

// RandomInt возвращает случайное целое число в диапазоне [min, max)
func RandomInt(min, max int) int {
	// Возвращаем случайное целое число в указанном диапазоне
	return min + rand.Intn(max-min)
}

// Lerp выполняет линейную интерполяцию между a и b по параметру t
func Lerp(a, b, t float64) float64 {
	// Ограничиваем t в диапазоне [0, 1]
	t = Clamp(t, 0, 1)
	
	// Выполняем линейную интерполяцию
	return a + (b-a)*t
}

// Angle вычисляет угол между двумя точками
func Angle(x1, y1, x2, y2 float64) float64 {
	// Вычисляем разницу координат
	dx := x2 - x1
	dy := y2 - y1
	
	// Возвращаем угол в радианах
	return math.Atan2(dy, dx)
}