// Пакет states содержит реализацию состояний игры
package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"superpupergame/ui"
	"os"
)

// MenuState реализует состояние главного меню игры
type MenuState struct {
	// stateMachine - ссылка на машину состояний для переключения состояний
	stateMachine *StateMachine
	
	// menuItems - элементы меню
	menuItems []*ui.Button
}

// NewMenuState создает новое состояние меню
func NewMenuState(stateMachine *StateMachine) *MenuState {
	// Создаем состояние меню
	menuState := &MenuState{
		stateMachine: stateMachine,
		menuItems:    make([]*ui.Button, 0),
	}
	
	// Добавляем кнопку "Начать новую игру"
	menuState.menuItems = append(menuState.menuItems, ui.NewButton(
		540, 300, 200, 50,
		"Start New Game",
		color.RGBA{0, 200, 0, 255},
		func() {
			// Переключаемся на состояние игры
			stateMachine.ChangeState("playing")
		},
	))
	
	// Добавляем кнопку "Настройки"
	menuState.menuItems = append(menuState.menuItems, ui.NewButton(
		540, 400, 200, 50,
		"Settings",
		color.RGBA{100, 100, 100, 255},
		func() {
			// Заглушка для настроек
		},
	))
	
	// Добавляем кнопку "Выход"
	menuState.menuItems = append(menuState.menuItems, ui.NewButton(
		540, 500, 200, 50,
		"Exit",
		color.RGBA{200, 0, 0, 255},
		func() {
			os.Exit(0)
		},
	))
	
	return menuState
}

// Enter вызывается при входе в состояние меню
func (m *MenuState) Enter() {
	// Инициализация состояния меню
}

// Update обновляет логику меню
func (m *MenuState) Update() error {
	// Проверяем нажатие кнопки мыши
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Получаем позицию курсора
		x, y := ebiten.CursorPosition()
		
		// Проверяем, была ли нажата какая-либо кнопка
		for _, button := range m.menuItems {
			if button.Contains(x, y) {
				// Вызываем функцию обработки нажатия кнопки
				button.OnClick()
			}
		}
	}
	
	return nil
}

// Draw отрисовывает меню
func (m *MenuState) Draw(screen *ebiten.Image) {
	// Заполняем фон
	ebitenutil.DrawRect(screen, 0, 0, 1280, 960, color.RGBA{50, 50, 50, 255})
	
	// Отрисовываем заголовок игры
	ebitenutil.DebugPrintAt(screen, "SuperPuperGame", 580, 200)
	
	// Отрисовываем кнопки меню
	for _, button := range m.menuItems {
		button.Draw(screen)
	}
}

// Exit вызывается при выходе из состояния меню
func (m *MenuState) Exit() {
	// Очистка ресурсов при выходе из состояния
}