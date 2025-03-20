package player

import (
	"image"
	"fmt"
	"log"
	"os"
	"time"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"superpupergame/debug" // Импортируем пакет debug
	"superpupergame/utils"
)

// Константы для настройки спрайта и анимации
const (
	FrameWidth     = 32   // Ширина одного кадра в пикселях
	FrameHeight    = 32   // Высота одного кадра в пикселях
	FramesPerState = 4    // Количество кадров в одном состоянии анимации
	ScaleFactor    = 3.0  // Коэффициент масштабирования спрайта
)

// Направления движения игрока
const (
	DirNone = iota  // Без движения
	DirUp           // Вверх
	DirDown         // Вниз
	DirLeft         // Влево
	DirRight        // Вправо
)

// Player представляет основную структуру игрока
type Player struct {
	X, Y           float64      // Позиция игрока на экране
	Speed          float64      // Базовая скорость движения
	Health         float64      // Текущее здоровье
	MaxHealth      float64      // Максимальное здоровье
	
	// Состояния
	Dying          bool         // Флаг смерти для анимации
	Attacking      bool         // Флаг выполнения атаки
	Dashing        bool         // Флаг выполнения рывка
	
	// Атрибуты атаки
	AttackAngle    float64      // Угол атаки (в радианах)
	AttackTimer    float64      // Таймер атаки
	LastAttackTime time.Time    // Время последней атаки
	
	// Атрибуты рывка
	DashSpeed      float64      // Скорость при рывке
	DashCharges    int          // Текущее количество зарядов рывка
	MaxDashes      int          // Максимальное количество зарядов рывка
	DirX, DirY     float64      // Компоненты вектора направления движения
	
	// Анимация
	FrameX         int          // Текущий кадр по X (индекс в спрайт-листе)
	FrameY         int          // Текущий кадр по Y (направление/состояние)
	FrameCount     int          // Счётчик для анимации
	DeathTimer     float64      // Таймер для анимации смерти
	
	// Ресурсы
	SpriteSheet    *ebiten.Image // Спрайт-лист игрока
	SwordImage     *ebiten.Image // Изображение меча
	
	// Система отладки
	DebugSystem    *debug.Debug // Ссылка на систему отладки
}

// NewPlayer создаёт и инициализирует нового игрока с указанными координатами
func NewPlayer(x, y float64, debugSystem *debug.Debug) *Player {
	// Загружаем изображение меча
	sword := loadImage("assets/sword.png")
	
	// Загружаем спрайт-лист игрока
	spriteSheet := loadImage("assets/player_sprites.png")
	
	// Создаём и возвращаем новый экземпляр игрока
	return &Player{
		X:              x,
		Y:              y,
		Speed:          2.0,                // Базовая скорость движения
		DashSpeed:      5.0,                // Скорость при рывке
		DirX:           0,                  // Начальное направление по X
		DirY:           0,                  // Начальное направление по Y
		DashCharges:    2,                  // Начальное количество зарядов рывка
		MaxDashes:      2,                  // Максимальное количество зарядов
		LastAttackTime: time.Time{},        // Нулевое время последней атаки
		SwordImage:     sword,              // Изображение меча
		SpriteSheet:    spriteSheet,        // Спрайт-лист игрока
		FrameX:         0,                  // Начальный кадр по X
		FrameY:         0,                  // Начальное состояние: стояние вправо
		Health:         100,                // Начальное здоровье
		MaxHealth:      100,                // Максимальное здоровье
		Dying:          false,              // Флаг смерти
		DeathTimer:     0,                  // Таймер смерти
		DebugSystem:    debugSystem,        // Система отладки
	}
}

// Update обновляет состояние игрока на каждом кадре
func (p *Player) Update() {
	// Если игрок умирает, обновляем только анимацию смерти
	if p.Dying {
		p.UpdateDeathAnimation()
		return
	}
	
	// Обновляем движение игрока (перенесено в movement.go)
	p.UpdateMovement()
	
	// Обновляем атаку игрока (перенесено в combat.go)
	p.UpdateCombat()
	
	// Обновляем анимацию игрока (перенесено в animation.go)
	p.UpdateAnimation()
	
	// Ограничиваем позицию игрока границами экрана
	p.X = utils.Clamp(p.X, 0, 1260)
	p.Y = utils.Clamp(p.Y, 0, 940)
}

// Draw отрисовывает игрока на экране
func (p *Player) Draw(screen *ebiten.Image) {
    // Отрисовка спрайта игрока
    p.DrawSprite(screen)

    // Отрисовка полосок заряда рывка
    if !p.Dying {
        p.DrawDashCharges(screen)
        if p.Attacking {
            p.DrawSword(screen)
        }
    }

    // Отрисовка отладочной информации (без хитбокса)
    if p.DebugSystem.IsEnabled() && p.DebugSystem.ShowPositions {
        debugInfo := fmt.Sprintf("X: %.1f, Y: %.1f", p.X, p.Y)
        ebitenutil.DebugPrintAt(screen, debugInfo, int(p.X), int(p.Y)-15)
        healthInfo := fmt.Sprintf("HP: %.1f/%.1f", p.Health, p.MaxHealth)
        ebitenutil.DebugPrintAt(screen, healthInfo, int(p.X), int(p.Y)-30)
        dashInfo := fmt.Sprintf("Dash: %d/%d", p.DashCharges, p.MaxDashes)
        ebitenutil.DebugPrintAt(screen, dashInfo, int(p.X), int(p.Y)-45)
    }
}

// loadImage загружает изображение из файла и возвращает его как ebiten.Image
func loadImage(path string) *ebiten.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Не удалось открыть файл %s: %v", path, err)
	}
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Не удалось декодировать изображение %s: %v", path, err)
	}
	
	return ebiten.NewImageFromImage(img)
}

// GetHitbox возвращает координаты и размеры хитбокса игрока
func (p *Player) GetHitbox() (x, y, width, height float64) {
    hitboxWidth := 20.0 	// Ширина хитбокса
    hitboxHeight := 30.0 	// Высота хитбокса
    x = p.X + 6        		// Смещение хитбокса по X
    y = p.Y + 2        		// Смещение хитбокса по Y
    return x, y, hitboxWidth, hitboxHeight
}