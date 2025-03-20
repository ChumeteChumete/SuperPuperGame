package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"superpupergame/debug" // Новый импорт для пакета отладки
	"superpupergame/player" // Импортируем пакет player
	"superpupergame/states"
)

// Game представляет главную структуру игры, реализующую интерфейс ebiten.Game
type Game struct {
	stateMachine *states.StateMachine
	player       *player.Player // Добавляем поле для игрока
	debugSystem  *debug.Debug   // Добавляем поле для системы отладки
}

// NewGame создает новый экземпляр игры
func NewGame() *Game {

	// Создаем систему отладки
	debugSystem := debug.NewDebug()

	gamePlayer := player.NewPlayer(640, 480, debugSystem)
	
	// Создаем новую игру
	game := &Game{
		stateMachine: states.NewStateMachine(),
		player:       gamePlayer,
		debugSystem:  debugSystem,
	}

	// Создаем и добавляем состояние меню
	menuState := states.NewMenuState(game.stateMachine)
	game.stateMachine.Add("menu", menuState)

	// Создаем и добавляем игровое состояние
	playState := states.NewPlayState(game.stateMachine, game.player) // Передаем игрока в игровое состояние
	game.stateMachine.Add("playing", playState)

	// Создаем и добавляем состояние смерти
	deathState := states.NewDeathState(game.stateMachine)
	game.stateMachine.Add("death", deathState)

	// Устанавливаем начальное состояние (меню)
	game.stateMachine.ChangeState("menu")

	return game
}

// Update обновляет игровую логику (реализация интерфейса ebiten.Game)
func (g *Game) Update() error {
	// Переключение режима отладки по F1
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.debugSystem.Toggle()
	}

	// Переключение показа FPS по F2
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) && g.debugSystem.IsEnabled() {
		g.debugSystem.ShowFPS = !g.debugSystem.ShowFPS
	}
		
	// Переключение показа хитбоксов по F3
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) && g.debugSystem.IsEnabled() {
		g.debugSystem.ShowHitboxes = !g.debugSystem.ShowHitboxes
	}
	
	// Делегируем обновление логики текущему состоянию
	return g.stateMachine.Update()
}

// Draw отрисовывает игру (реализация интерфейса ebiten.Game)
func (g *Game) Draw(screen *ebiten.Image) {
	// Делегируем отрисовку текущему состоянию
	g.stateMachine.Draw(screen)

	// Если включен режим отладки, выполняем отладочные действия
	if g.debugSystem.IsEnabled() {
		// Отображаем отладочную информацию
        g.debugSystem.DrawDebugInfo(screen, ebiten.ActualTPS())
        
        // Добавляем информацию о текущем состоянии
		currentState := g.stateMachine.GetCurrentStateName()
		g.debugSystem.AddMessage("Текущее состояние: " + currentState)
	}

	// Отрисовка игрока - это должно происходить в соответствующем состоянии,
	// но для простоты можно оставить здесь
	// g.player.Draw(screen)
}

// Layout определяет логический размер игры (реализация интерфейса ebiten.Game)
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Возвращаем фиксированный размер игрового мира
	return 1280, 960
}

func main() {
	// Создаем новую игру
	game := NewGame()
	
	// Настраиваем окно игры
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("SuperPuperGame")
	
	// Запускаем игровой цикл
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}