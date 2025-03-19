// Пакет ui содержит элементы пользовательского интерфейса
package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// Menu представляет меню игры
type Menu struct {
	// Title - заголовок меню
	Title string
	
	// Items - пункты меню
	Items []*Button
	
	// BackgroundColor - цвет фона меню
	BackgroundColor color.RGBA
}

// NewMenu создает новое меню
func NewMenu(title string, backgroundColor color.RGBA) *Menu {
	// Создаем и возвращаем новое меню
	return &Menu{
		Title:           title,
		Items:           make([]*Button, 0),
		BackgroundColor: backgroundColor,
	}
}

// AddItem добавляет новый пункт в меню
func (m *Menu) AddItem(text string, color color.RGBA, onClick func()) *Button {
	// Вычисляем позицию для нового пункта меню
	y := 300 + float64(len(m.Items)*100)
	
	// Создаем новую кнопку
	button := NewButton(540, y, 200, 50, text, color, onClick)
	
	// Добавляем кнопку в список пунктов меню
	m.Items = append(m.Items, button)
	
	return button
}

// Update обновляет логику меню
func (m *Menu) Update() {
	// Проверяем нажатие кнопки мыши
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Получаем позицию курсора
		x, y := ebiten.CursorPosition()
		
		// Проверяем, была ли нажата какая-либо кнопка
		for _, button := range m.Items {
			if button.Contains(x, y) {
				// Вызываем функцию обработки нажатия кнопки
				button.OnClick()
			}
		}
	}
}

// Draw отрисовывает меню
func (m *Menu) Draw(screen *ebiten.Image) {
	// Заполняем фон
	ebitenutil.DrawRect(screen, 0, 0, 1280, 960, m.BackgroundColor)
	
	// Отрисовываем заголовок меню
	ebitenutil.DebugPrintAt(screen, m.Title, 580, 200)
	
	// Отрисовываем пункты меню
	for _, button := range m.Items {
		button.Draw(screen)
	}
}