// Пакет states содержит реализацию состояний игры
package states

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"math"
	"superpupergame/ui"
)

// DeathState реализует состояние смерти игрока
type DeathState struct {
	// stateMachine - ссылка на машину состояний для переключения состояний
	stateMachine *StateMachine
	
	// player - игрок, который умер
	player interface {
		Draw(screen *ebiten.Image)
		UpdateDeathAnimation()
	}
	
	// score - итоговый счет
	score int
	
	// deathTimer - таймер с момента смерти
	deathTimer float64
	
	// buttons - кнопки на экране смерти
	buttons []*ui.Button
}

// NewDeathState создает новое состояние смерти
func NewDeathState(stateMachine *StateMachine) *DeathState {
	// Создаем состояние смерти
	deathState := &DeathState{
		stateMachine: stateMachine,
		buttons:      make([]*ui.Button, 0),
	}
	
	// Добавляем кнопку "Начать заново"
	deathState.buttons = append(deathState.buttons, ui.NewButton(
		540, 400, 200, 50,
		"Restart",
		color.RGBA{0, 200, 0, 255},
		func() {
			// Переключаемся на состояние игры
			stateMachine.ChangeState("playing")
		},
	))
	
	// Добавляем кнопку "В меню"
	deathState.buttons = append(deathState.buttons, ui.NewButton(
		540, 500, 200, 50,
		"Main Menu",
		color.RGBA{100, 100, 100, 255},
		func() {
			// Переключаемся на состояние меню
			stateMachine.ChangeState("menu")
		},
	))
	
	return deathState
}

// Enter вызывается при входе в состояние смерти
func (d *DeathState) Enter() {
	// Сбрасываем таймер смерти
	d.deathTimer = 0
	
	// Получаем текущие данные из игрового состояния
	if playState, ok := d.stateMachine.states["playing"].(*PlayState); ok {
		d.player = playState.player
		d.score = playState.score
	}
}

// Update обновляет логику состояния смерти
func (d *DeathState) Update() error {
	// Увеличиваем таймер смерти
	d.deathTimer += 1.0 / 60.0
	
	// Обновляем анимацию смерти
	d.player.UpdateDeathAnimation()
	
	// Обработка кнопок на экране смерти
	if d.deathTimer > 2 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Получаем позицию курсора
		x, y := ebiten.CursorPosition()
		
		// Проверяем, была ли нажата какая-либо кнопка
		for _, button := range d.buttons {
			if button.Contains(x, y) {
				// Вызываем функцию обработки нажатия кнопки
				button.OnClick()
			}
		}
	}
	
	return nil
}

// Draw отрисовывает состояние смерти
func (d *DeathState) Draw(screen *ebiten.Image) {
	// Заполняем фон
	ebitenutil.DrawRect(screen, 0, 0, 1280, 960, color.RGBA{50, 50, 50, 255})
	
	// Отрисовываем умирающего игрока
	d.player.Draw(screen)
	
	// Затемняем экран
	overlay := ebiten.NewImage(1280, 960)
	overlay.Fill(color.RGBA{0, 0, 0, uint8(math.Min(192, d.deathTimer*80))})
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(overlay, op)
	
	// Отображаем текст и кнопки после задержки
	if d.deathTimer > 1 {
		// Отрисовываем текст "Game Over"
		ebitenutil.DebugPrintAt(screen, "GAME OVER", 580, 300)
		
		// Показываем итоговый счёт
		scoreText := fmt.Sprintf("Final Score: %d", d.score)
		ebitenutil.DebugPrintAt(screen, scoreText, 580, 350)
		
		// Отрисовываем кнопки после короткой задержки
		if d.deathTimer > 2 {
			for _, button := range d.buttons {
				button.Draw(screen)
			}
		}
	}
}

// Exit вызывается при выходе из состояния смерти
func (d *DeathState) Exit() {
	// Очистка ресурсов при выходе из состояния
}