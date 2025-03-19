package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"superpupergame/states" // Импортируем пакет states для машины состояний
)

// Game представляет главную структуру игры, реализующую интерфейс ebiten.Game
type Game struct {
	// stateMachine - машина состояний для управления состояниями игры
	stateMachine *states.StateMachine
}

// NewGame создает новый экземпляр игры
func NewGame() *Game {
	// Создаем новую игру
	game := &Game{
		// Инициализируем машину состояний
		stateMachine: states.NewStateMachine(),
	}

	// Создаем и добавляем состояние меню
	menuState := states.NewMenuState(game.stateMachine)
	game.stateMachine.Add("menu", menuState)

	// Создаем и добавляем игровое состояние
	playState := states.NewPlayState(game.stateMachine)
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
	// Делегируем обновление логики текущему состоянию
	return g.stateMachine.Update()
}

// Draw отрисовывает игру (реализация интерфейса ebiten.Game)
func (g *Game) Draw(screen *ebiten.Image) {
	// Делегируем отрисовку текущему состоянию
	g.stateMachine.Draw(screen)
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