package player

import (
	"math"
	"time"
	
	"github.com/hajimehoshi/ebiten/v2"
)

// UpdateMovement обрабатывает движение игрока и рывки
func (p *Player) UpdateMovement() {
	// Сбрасываем направление движения
	p.DirX, p.DirY = 0, 0
	
	// Обработка клавиш направления
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.DirY = -1         // Движение вверх
		p.FrameY = 2        // Устанавливаем кадр анимации для направления вверх
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.DirY = 1          // Движение вниз
		p.FrameY = 3        // Устанавливаем кадр анимации для направления вниз
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.DirX = -1         // Движение влево
		p.FrameY = 4        // Устанавливаем кадр анимации для направления влево
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.DirX = 1          // Движение вправо
		p.FrameY = 4        // Устанавливаем кадр анимации для направления вправо
	}
	
	// Нормализация вектора движения при диагональном движении
	if p.DirX != 0 && p.DirY != 0 {
		// Вычисляем длину вектора направления
		magnitude := math.Sqrt(p.DirX*p.DirX + p.DirY*p.DirY)
		// Нормализуем вектор для равномерной скорости по диагонали
		p.DirX = p.DirX / magnitude
		p.DirY = p.DirY / magnitude
		
		// Выбираем приоритетное направление для анимации
		if math.Abs(p.DirX) > math.Abs(p.DirY) {
			// Движение больше по горизонтали
			if p.DirX < 0 {
				p.FrameY = 3  // Влево
			} else {
				p.FrameY = 0  // Вправо
			}
		} else {
			// Движение больше по вертикали
			if p.DirY < 0 {
				p.FrameY = 2  // Вверх
			} else {
				p.FrameY = 0  // Вниз
			}
		}
	}
	
	// Применяем текущую скорость движения
	currentSpeed := p.Speed
	if p.Dashing {
		currentSpeed = p.DashSpeed  // Используем скорость рывка
	}
	
	// Обновляем позицию игрока
	p.X += p.DirX * currentSpeed
	p.Y += p.DirY * currentSpeed
	
	// Обработка рывка (dash)
	p.handleDash()
}

// handleDash обрабатывает логику рывка
func (p *Player) handleDash() {
	// Проверяем возможность рывка:
	// 1. Нажата клавиша пробел
	// 2. Игрок не выполняет рывок в данный момент
	// 3. Игрок движется (есть направление)
	// 4. Есть заряды рывка
	if ebiten.IsKeyPressed(ebiten.KeySpace) && 
	   !p.Dashing && 
	   (p.DirX != 0 || p.DirY != 0) && 
	   p.DashCharges > 0 {
		// Активируем рывок
		p.Dashing = true
		// Уменьшаем количество зарядов
		p.DashCharges--
		// Запускаем восстановление заряда
		go p.rechargeDash()
		// Ограничиваем длительность рывка
		go func() {
			time.Sleep(200 * time.Millisecond)
			p.Dashing = false
		}()
	}
}

// rechargeDash восстанавливает заряд рывка через указанное время
func (p *Player) rechargeDash() {
	// Ждем 5 секунд перед восстановлением заряда
	time.Sleep(5 * time.Second)
	// Восстанавливаем заряд если не достигнут максимум
	if p.DashCharges < p.MaxDashes {
		p.DashCharges++
	}
}

// DrawDashCharges отрисовывает индикаторы зарядов рывка
func (p *Player) DrawDashCharges(screen *ebiten.Image) {
	// Импорт нужно добавить в файл: "image/color"
	// Отрисовка зарядов рывка
	// import "image/color"
	// for i := 0; i < p.DashCharges; i++ {
	//    img := ebiten.NewImage(10, 5)
	//    img.Fill(color.RGBA{255, 255, 0, 255})
	//    opCharge := &ebiten.DrawImageOptions{}
	//    opCharge.GeoM.Translate(p.X+float64(i*15), p.Y-15)
	//    screen.DrawImage(img, opCharge)
	// }
}