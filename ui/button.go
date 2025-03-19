// Пакет ui содержит элементы пользовательского интерфейса
package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// Button представляет кнопку в пользовательском интерфейсе
type Button struct {
	// X - координата X верхнего левого угла кнопки
	X float64
	
	// Y - координата Y верхнего левого угла кнопки
	Y float64
	
	// Width - ширина кнопки
	Width float64
	
	// Height - высота кнопки
	Height float64
	
	// Text - текст на кнопке
	Text string
	
	// Color - цвет кнопки
	Color color.RGBA
	
	// OnClick - функция, вызываемая при нажатии на кнопку
	OnClick func()
}

// NewButton создает новую кнопку
func NewButton(x, y, width, height float64, text string, color color.RGBA, onClick func()) *Button {
	// Создаем и возвращаем новую кнопку
	return &Button{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Text:    text,
		Color:   color,
		OnClick: onClick,
	}
}

// Contains проверяет, содержит ли кнопка указанную точку
func (b *Button) Contains(x, y int) bool {
	// Проверяем, находится ли точка внутри кнопки
	return float64(x) >= b.X && float64(x) <= b.X+b.Width &&
		float64(y) >= b.Y && float64(y) <= b.Y+b.Height
}

// Draw отрисовывает кнопку
func (b *Button) Draw(screen *ebiten.Image) {
	// Рисуем прямоугольник кнопки
	ebitenutil.DrawRect(screen, b.X, b.Y, b.Width, b.Height, b.Color)
	
	// Рассчитываем положение текста для центрирования
	textX := b.X + (b.Width-float64(len(b.Text)*6))/2
	textY := b.Y + b.Height/2
	
	// Отрисовываем текст кнопки
	ebitenutil.DebugPrintAt(screen, b.Text, int(textX), int(textY))
}