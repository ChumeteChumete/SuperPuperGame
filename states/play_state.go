// Пакет states содержит реализацию состояний игры
package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math/rand"
	"superpupergame/enemy"
	"superpupergame/game" // Импортируем пакет с монеткой
	"superpupergame/player"
	"superpupergame/ui"
	"time"
	"math"
	"fmt"
)

// PlayState реализует игровое состояние
type PlayState struct {
	// stateMachine - ссылка на машину состояний для переключения состояний
	stateMachine *StateMachine
	
	// player - игрок
	player *player.Player
	
	// enemies - список врагов
	enemies []*enemy.Enemy
	
	// coins - список монеток
	coins []*game.Coin
	
	// coinCount - количество активных монеток
	coinCount int
	
	// maxCoins - максимальное количество монеток на экране
	maxCoins int
	
	// enemyCount - количество врагов в текущей волне
	enemyCount int
	
	// score - текущий счет
	score int
	
	// hud - элементы интерфейса
	hud *ui.HUD
}

// NewPlayState создает новое игровое состояние
func NewPlayState(stateMachine *StateMachine, player *player.Player) *PlayState {
	// Создаем игровое состояние
	return &PlayState{
		stateMachine: stateMachine,
		player:       player,
		enemyCount:   1,
		score:        0,
		hud:          ui.NewHUD(),
		coins:        make([]*game.Coin, 0),
		coinCount:    0,
		maxCoins:     5, // Максимальное количество монеток на экране
	}
}

// SpawnCoin создает новую монетку
func (p *PlayState) SpawnCoin() {
	// Создаем новую монетку если не превышен лимит
	if p.coinCount < p.maxCoins {
		p.coins = append(p.coins, game.NewCoin(1280, 960))
		p.coinCount++
	}
}

// Enter вызывается при входе в игровое состояние
func (p *PlayState) Enter() {
	
	// Сбрасываем параметры существующего игрока
    p.player.X = 640
    p.player.Y = 480
    p.player.Health = 100
    p.player.Dying = false
    p.player.Attacking = false
    p.player.Dashing = false
    p.player.DashCharges = p.player.MaxDashes
    p.player.DeathTimer = 0
	
	// Создаем первого врага
	p.enemies = []*enemy.Enemy{enemy.NewRandomEdgeEnemy()}
	
	// Очищаем список монеток
	p.coins = make([]*game.Coin, 0)
	p.coinCount = 0
	
	// Создаем начальные монетки
	for i := 0; i < 3; i++ {
		p.SpawnCoin()
	}
	
	// Сбрасываем счет
	p.score = 0
	
	// Сбрасываем количество врагов
	p.enemyCount = 1
}

// Update обновляет игровую логику
func (p *PlayState) Update() error {
	// Обновляем игрока
	p.player.Update()
	
	// Обновляем все монетки (анимация)
	for _, coin := range p.coins {
        coin.Update()
    }

	// Добавляем отладочную информацию о количестве объектов
	if p.player.DebugSystem != nil && p.player.DebugSystem.IsEnabled() {
        p.player.DebugSystem.ClearMessages()
        p.player.DebugSystem.AddMessage(fmt.Sprintf("Враги: %d", len(p.enemies)))
        p.player.DebugSystem.AddMessage(fmt.Sprintf("Монеты: %d/%d", p.coinCount, p.maxCoins))
        p.player.DebugSystem.AddMessage(fmt.Sprintf("Счёт: %d", p.score))
        p.player.DebugSystem.AddMessage(fmt.Sprintf("Волна: %d", p.enemyCount))
    }

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

	// Проверяем сбор монеток
	for i := len(p.coins) - 1; i >= 0; i-- {
		coin := p.coins[i]
		// Проверяем коллизию с игроком
		if coin.Collides(p.player.X, p.player.Y, 20, 20) { // Предполагаемый размер игрока
			// Увеличиваем счет
			p.score += 50
			
			// Удаляем монетку
			p.coins = append(p.coins[:i], p.coins[i+1:]...)
			p.coinCount--
			
			// Создаем новую монетку с небольшой задержкой
			go func() {
				time.Sleep(2 * time.Second)
				p.SpawnCoin()
			}()
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
					
					// С небольшим шансом создаем дополнительную монетку
					if p.coinCount < p.maxCoins && rand.Float64() < 0.3 {
						p.SpawnCoin()
					}
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
		
		// Создаем бонусную монетку после каждой волны
		p.SpawnCoin()
		
		// Небольшая пауза между волнами
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (p *PlayState) Draw(screen *ebiten.Image) {
    // Заполняем фон
    ebitenutil.DrawRect(screen, 0, 0, 1280, 960, color.RGBA{50, 50, 50, 255})

    // Отрисовываем игрока
    p.player.Draw(screen)

    // Отрисовываем монетки
    for _, coin := range p.coins {
        coin.Draw(screen)
    }

    // Отрисовываем врагов
    for _, e := range p.enemies {
        e.Draw(screen)
    }

    // Отрисовываем HUD (здоровье и счет)
    p.hud.Draw(screen, p.player.Health, p.score)

    // Отрисовываем хитбоксы в режиме отладки
    if p.player.DebugSystem != nil && p.player.DebugSystem.IsEnabled() && p.player.DebugSystem.ShowHitboxes {

        // Хитбокс игрока
        px, py, pw, ph := p.player.GetHitbox()
        p.player.DebugSystem.DrawHitbox(screen, px, py, pw, ph)

        // Хитбоксы монеток
        for _, coin := range p.coins {
            cx, cy, cw, ch := coin.GetHitbox()
            p.player.DebugSystem.DrawHitbox(screen, cx, cy, cw, ch)
        }

        // Хитбоксы врагов
        for _, e := range p.enemies {
            if e.Alive {
                ex, ey, ew, eh := e.GetHitbox()
                p.player.DebugSystem.DrawHitbox(screen, ex, ey, ew, eh)
            }
        }
    }
}

// Exit вызывается при выходе из игрового состояния
func (p *PlayState) Exit() {
	// Очистка ресурсов при выходе из состояния
}