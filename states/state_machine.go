// Пакет states содержит реализацию машины состояний для игры
package states

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// State - интерфейс для всех состояний игры
type State interface {
	// Enter вызывается при входе в состояние
	Enter()
	
	// Update обновляет логику состояния
	Update() error
	
	// Draw отрисовывает состояние на экране
	Draw(screen *ebiten.Image)
	
	// Exit вызывается при выходе из состояния
	Exit()
}

// StateMachine управляет переключением между состояниями
type StateMachine struct {
	// currentState - текущее активное состояние
	currentState State
	
	// states - карта всех доступных состояний
	states map[string]State
}

// NewStateMachine создает новую машину состояний
func NewStateMachine() *StateMachine {
	// Инициализируем машину состояний
	return &StateMachine{
		states: make(map[string]State),
	}
}

// Add добавляет новое состояние в машину состояний
func (sm *StateMachine) Add(name string, state State) {
	// Добавляем состояние в карту состояний
	sm.states[name] = state
}

// ChangeState меняет текущее состояние
func (sm *StateMachine) ChangeState(name string) {
	// Выходим из текущего состояния, если оно есть
	if sm.currentState != nil {
		sm.currentState.Exit()
	}
	
	// Получаем новое состояние из карты
	sm.currentState = sm.states[name]
	
	// Входим в новое состояние
	if sm.currentState != nil {
		sm.currentState.Enter()
	}
}

// Update обновляет текущее состояние
func (sm *StateMachine) Update() error {
	// Проверяем, что текущее состояние существует
	if sm.currentState != nil {
		// Обновляем текущее состояние
		return sm.currentState.Update()
	}
	return nil
}

// Draw отрисовывает текущее состояние
func (sm *StateMachine) Draw(screen *ebiten.Image) {
	// Проверяем, что текущее состояние существует
	if sm.currentState != nil {
		// Отрисовываем текущее состояние
		sm.currentState.Draw(screen)
	}
}