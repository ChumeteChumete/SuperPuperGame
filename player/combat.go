package player

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// UpdateCombat обрабатывает атаку игрока
func (p *Player) UpdateCombat() {
	// Определяем направление атаки по позиции курсора
	cursorX, cursorY := ebiten.CursorPosition()
	
	// Преобразование координат курсора из оконных в игровые
	gameWidth, gameHeight := 1280.0, 960.0
	w, h := ebiten.WindowSize()
	windowWidth, windowHeight := float64(w), float64(h)
	
	// Корректируем координаты курсора в соответствии с масштабом окна
	cursorX = int(float64(cursorX) * gameWidth / windowWidth)
	cursorY = int(float64(cursorY) * gameHeight / windowHeight)
	
	// Вычисляем вектор от игрока до курсора
	dx := float64(cursorX) - (p.X + 10)
	dy := float64(cursorY) - (p.Y + 10)
	
	// Нормализуем вектор направления
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		dx /= length
		dy /= length
	}
	
	// Запоминаем угол атаки
	p.AttackAngle = math.Atan2(dy, dx)
	
	// Обработка нажатия левой кнопки мыши для атаки
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !p.Attacking {
		// Проверяем время с последней атаки (кулдаун)
		if time.Since(p.LastAttackTime) >= 500*time.Millisecond {
			// Активируем атаку
			p.Attacking = true
			p.AttackTimer = 0
			p.LastAttackTime = time.Now()
			p.FrameX = 0  // Сбрасываем кадр анимации
			
			// Устанавливаем длительность атаки
			go func() {
				time.Sleep(300 * time.Millisecond)
				p.Attacking = false
				p.FrameX = 0
			}()
		}
	}
	
	// Обновляем таймер атаки
	if p.Attacking {
		p.AttackTimer += 1.0 / 60.0  // Увеличиваем таймер (60 FPS)
		
		// Обновляем анимацию атаки
		if p.FrameCount%10 == 0 {
			p.FrameX = (p.FrameX + 1) % FramesPerState
		}
	}
}

// DrawSword отрисовывает меч при атаке
func (p *Player) DrawSword(screen *ebiten.Image) {
	// Добавляем небольшое колебание для эффекта взмаха
	oscillation := math.Sin(p.AttackTimer * 10) * 0.5
	angle := p.AttackAngle + oscillation
	
	// Настраиваем параметры отрисовки
	opSword := &ebiten.DrawImageOptions{}
	
	// Центрируем меч относительно его середины
	swordWidth := float64(p.SwordImage.Bounds().Dx())
	opSword.GeoM.Translate(-swordWidth/2, 0)
	
	// Поворот меча на нужный угол (+90 градусов для правильной ориентации)
	opSword.GeoM.Rotate(angle + math.Pi/2)
	
	// Получаем координаты хитбокса игрока
	hitboxX, hitboxY, hitboxWidth, hitboxHeight := p.GetHitbox()
	
	// Смещение меча от игрока
	offsetDistance := 35 * ScaleFactor
	
	// Вычисляем позицию меча относительно центра хитбокса
	swordX := hitboxX + hitboxWidth/2 + math.Cos(angle)*offsetDistance
	swordY := hitboxY + hitboxHeight/2 + math.Sin(angle)*offsetDistance
	
	// Применяем смещение и отрисовываем меч
	opSword.GeoM.Translate(swordX, swordY)
	screen.DrawImage(p.SwordImage, opSword)
}

// AttackArea возвращает область атаки меча как прямоугольник
func (p *Player) AttackArea() (x, y, width, height float64) {
	if p.Attacking && !p.Dying {
		// Добавляем небольшое колебание для эффекта взмаха
		oscillation := math.Sin(p.AttackTimer * 10) * 0.5
		angle := p.AttackAngle + oscillation
		
		// Рассчитываем смещение меча от игрока
		offsetDistance := 15.0 * ScaleFactor
		offsetX := math.Cos(p.AttackAngle) * offsetDistance
		offsetY := math.Sin(p.AttackAngle) * offsetDistance
		
		// Получаем размеры меча
		swordLength := float64(p.SwordImage.Bounds().Dx())
		swordWidth := float64(p.SwordImage.Bounds().Dy())
		
		// Начальная позиция меча (рукоять)
		swordX := p.X + 10 + offsetX
		swordY := p.Y + 10 + offsetY
		
		// Конечная позиция меча (кончик)
		endX := swordX + math.Cos(angle)*(swordLength-10)
		endY := swordY + math.Sin(angle)*(swordLength-10)
		
		// Находим границы прямоугольника, описывающего область атаки
		minX := math.Min(swordX, endX)
		maxX := math.Max(swordX, endX)
		minY := math.Min(swordY, endY)
		maxY := math.Max(swordY, endY)
		
		// Добавляем толщину меча к размерам прямоугольника
		width := maxX - minX + swordWidth
		height := maxY - minY + swordWidth
		
		return minX - swordWidth/2, minY - swordWidth/2, width, height
	}
	
	// Если игрок не атакует или умирает, возвращаем нулевую область
	return 0, 0, 0, 0
}