// Пакет states содержит реализацию состояний игры
package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
	"superpupergame/enemy"
	"superpupergame/player"
	"superpupergame/ui"
	"time"
)

// PlayState реализует игровое состояние
type PlayState struct {
	// stateMachine - ссылка на машину состояний для переключения состояний
	stateMachine *StateMachine
	
	// player - игрок
	player *player.Player
	
	// enemies - список врагов
	enemies []*enemy.Enemy
	
	// enemyCount - количество врагов в текущей волне
	enemyCount int
	
	// score - текущий счет
	score int
	
	// hud - элементы интерфейса
	hud *ui.HUD
}

// NewPlayState создает новое игровое состояние
func NewPlayState(stateMachine *StateMachine) *PlayState {
	// Создаем игровое состояние
	return &PlayState{
		stateMachine: stateMachine,
		enemyCount:   1,
		score:        0,
		hud:          ui.NewHUD(),
	}
}

// Enter вызывается при входе в игровое состояние
func (p *PlayState) Enter() {
	// Создаем игрока в центре экрана
	p.player = player.NewPlayer(640, 480)
	
	// Создаем первого врага
	p.enemies = []*enemy.Enemy{enemy.NewRandomEdgeEnemy()}
	
	// Сбрасываем счет
	p.score = 0
	
	// Сбрасываем количество врагов
	p.enemyCount = 1
}

// Update обновляет игровую логику
func (p *PlayState) Update() error {
	// Обновляем игрока
	p.player.Update()

	// Обрабатываем взаимодействие с врагами
	for _, e := range p.enemies {
		// Обновляем врага, передавая позицию игрока как цель
		e.Update(p.player.X+10, p.player.Y+10)

		// Вычисляем расстояние между игроком и врагом
		dx := p.player.X + 10 - e.X
		dy := p.player.Y + 10 - e.Y
		distance := math.Sqrt(dx*dx + dy*dy)
		
		// Проверяем столкновение с врагом
		if distance < 20 && e.Alive {
			// Уменьшаем здоровье при контакте с врагом
			p.player.Health -= 25
			
			// Проверяем, умер ли игрок
			if p.player.Health <= 0 {
				// Запускаем анимацию смерти
				p.player.StartDeathAnimation()
				
				// Переходим в состояние смерти
				p.stateMachine.ChangeState("death")
				return nil
			}
			
			// Отталкиваем игрока от врага
			pushDirection := math.Atan2(dy, dx)
			pushDistance := 50.0
			p.player.X += math.Cos(pushDirection) * pushDistance
			p.player.Y += math.Sin(pushDirection) * pushDistance
		}
	}

	// Подсчитываем живых врагов и проверяем атаки
	liveEnemies := 0
	for _, e := range p.enemies {
		if e.Alive {
			liveEnemies++
			
			// Получаем область атаки игрока
			attackX, attackY, attackW, attackH := p.player.AttackArea()
			
			// Если игрок атакует
			if attackX != 0 {
				// Вычисляем границы области атаки
				attackLeft := attackX
				attackRight := attackX + attackW
				attackTop := attackY
				attackBottom := attackY + attackH

				// Вычисляем границы врага
				enemyLeft := e.X
				enemyRight := e.X + 20
				enemyTop := e.Y
				enemyBottom := e.Y + 20

				// Проверяем пересечение областей атаки и врага
				if attackLeft < enemyRight && attackRight > enemyLeft &&
					attackTop < enemyBottom && attackBottom > enemyTop {
					// Уничтожаем врага
					e.Alive = false
					
					// Увеличиваем счет
					p.score += 100
				}
			}
		}
	}

	// Если все враги уничтожены, создаем новую волну
	if liveEnemies == 0 {
		// Увеличиваем количество врагов
		p.enemyCount++
		
		// Очищаем список врагов
		p.enemies = nil
		
		// Создаем новых врагов
		for i := 0; i < p.enemyCount; i++ {
			p.enemies = append(p.enemies, enemy.NewRandomEdgeEnemy())
		}
		
		// Восстанавливаем немного здоровья при уничтожении всех врагов
		p.player.Health = math.Min(p.player.Health+10, 100)
		
		// Небольшая пауза между волнами
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// Draw отрисовывает игровое состояние
func (p *PlayState) Draw(screen *ebiten.Image) {
	// Заполняем фон
	ebitenutil.DrawRect(screen, 0, 0, 1280, 960, color.RGBA{50, 50, 50, 255})
	
	// Отрисовываем игрока
	p.player.Draw(screen)
	
	// Отрисовываем врагов
	for _, e := range p.enemies {
		e.Draw(screen)
	}
	
	// Отрисовываем HUD (здоровье и счет)
	p.hud.Draw(screen, p.player.Health, p.score)
}

// Exit вызывается при выходе из игрового состояния
func (p *PlayState) Exit() {
	// Очистка ресурсов при выходе из состояния
}