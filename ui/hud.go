// Пакет ui содержит элементы пользовательского интерфейса
package ui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// HUD представляет элементы интерфейса во время игры
type HUD struct {
	// Здесь могут быть различные настройки HUD
}

// NewHUD создает новый экземпляр HUD
func NewHUD() *HUD {
	// Создаем и возвращаем новый HUD
	return &HUD{}
}

// Draw отрисовывает все элементы HUD
func (h *HUD) Draw(screen *ebiten.Image, health float64, score int) {
	// Отрисовываем полоску здоровья
	h.DrawHealthBar(screen, 20, 20, 200, 20, health)
	
	// Отрисовываем счет
	scoreText := fmt.Sprintf("Score: %d", score)
	ebitenutil.DebugPrintAt(screen, scoreText, 20, 50)
}

// DrawHealthBar отрисовывает полоску здоровья
func (h *HUD) DrawHealthBar(screen *ebiten.Image, x, y, width, height float64, health float64) {
	// Фон полоски здоровья (серый)
	ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{100, 100, 100, 255})
	
	// Заполнение полоски здоровья (зеленый -> желтый -> красный)
	healthWidth := (health / 100) * width
	var healthColor color.RGBA
	
	// Выбираем цвет в зависимости от уровня здоровья
	if health > 70 {
		healthColor = color.RGBA{0, 200, 0, 255} // Зеленый
	} else if health > 30 {
		healthColor = color.RGBA{200, 200, 0, 255} // Желтый
	} else {
		healthColor = color.RGBA{200, 0, 0, 255} // Красный
	}
	
	// Отрисовываем заполненную часть полоски здоровья
	ebitenutil.DrawRect(screen, x, y, healthWidth, height, healthColor)
	
	// Рисуем рамку полоски здоровья
	ebitenutil.DrawLine(screen, x, y, x+width, y, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawLine(screen, x, y+height, x+width, y+height, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawLine(screen, x, y, x, y+height, color.RGBA{200, 200, 200, 255})
	ebitenutil.DrawLine(screen, x+width, y, x+width, y+height, color.RGBA{200, 200, 200, 255})
}