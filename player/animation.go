package player

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// UpdateAnimation обновляет анимацию игрока
func (p *Player) UpdateAnimation() {
	// Обновляем анимацию движения
	if p.DirX != 0 || p.DirY != 0 {
		// Увеличиваем счетчик кадров
		p.FrameCount++
		
		// Меняем кадр анимации каждые 10 обновлений
		if p.FrameCount % 10 == 0 {
			p.FrameX = (p.FrameX + 1) % FramesPerState
		}
	} else {
		// Если игрок стоит на месте, используем первый кадр
		p.FrameX = 0
	}
}

// UpdateDeathAnimation обновляет анимацию смерти
func (p *Player) UpdateDeathAnimation() {
	// Если игрок не умирает, не выполняем анимацию
	if !p.Dying {
		return
	}
	
	// Увеличиваем таймер смерти
	p.DeathTimer += 1.0 / 60.0
	
	// Замедленная анимация падения
	if p.DeathTimer > 0.5 && p.FrameCount%20 == 0 {
		p.FrameX = (p.FrameX + 1) % FramesPerState
	}
	
	// Обновляем счетчик кадров
	p.FrameCount++
}

// StartDeathAnimation запускает анимацию смерти
func (p *Player) StartDeathAnimation() {
	p.Dying = true      // Устанавливаем флаг смерти
	p.DeathTimer = 0    // Сбрасываем таймер
	p.FrameX = 0        // Начинаем с первого кадра
	p.FrameY = 0        // Используем кадры для основного состояния
}

// DrawSprite отрисовывает спрайт игрока с учетом текущего состояния
func (p *Player) DrawSprite(screen *ebiten.Image) {
	// Вырезаем текущий кадр из спрайт-листа
	r := image.Rect(
		p.FrameX*FrameWidth, p.FrameY*FrameHeight, 
		(p.FrameX+1)*FrameWidth, (p.FrameY+1)*FrameHeight,
	)
	subImage := p.SpriteSheet.SubImage(r).(*ebiten.Image)
	
	// Рассчитываем масштабированные размеры
	scaledWidth := float64(FrameWidth) * ScaleFactor
	scaledHeight := float64(FrameHeight) * ScaleFactor
	
	// Настраиваем параметры отрисовки
	op := &ebiten.DrawImageOptions{}
	
	// Смещение для центрирования спрайта
	offsetX := -scaledWidth/2 + 16.0  // Корректировка по X (влево)
	offsetY := -scaledHeight/2 + 15.0  // Корректировка по Y (вверх)
	
	// Центрируем спрайт
	op.GeoM.Translate(-float64(FrameWidth)/2, -float64(FrameHeight)/2)
	
	// Масштабируем спрайт
	op.GeoM.Scale(ScaleFactor, ScaleFactor)
	
	// Отражаем спрайт для направления влево
	if p.DirX < 0 && !p.Dying {
		op.GeoM.Scale(-1, 1) // Отражение по горизонтали
	}
	
	// Добавляем эффекты для анимации смерти
	if p.Dying {
		// Эффект падения (поворот)
		rotationAngle := p.DeathTimer * 0.3 // Медленное падение
		if rotationAngle > math.Pi/2 {
			rotationAngle = math.Pi/2 // Ограничиваем поворот до 90 градусов
		}
		op.GeoM.Rotate(rotationAngle)
		
		// Эффект затухания (прозрачность)
		opacity := 1.0 - math.Min(p.DeathTimer/3.0, 0.5)
		op.ColorM.Scale(1, 1, 1, opacity)
	}
	
	// Перемещаем спрайт в позицию игрока с корректировкой
	op.GeoM.Translate(p.X + scaledWidth/2 + offsetX, p.Y + scaledHeight/2 + offsetY)
	
	// Отрисовываем спрайт
	screen.DrawImage(subImage, op)
}