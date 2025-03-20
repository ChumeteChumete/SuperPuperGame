package enemy

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"superpupergame/utils"
)

type Enemy struct {
	X, Y   float64
	Speed  float64
	Alive  bool
}

func NewEnemy(x, y float64) *Enemy {
	return &Enemy{
		X:     x,
		Y:     y,
		Speed: 3.0,
		Alive: true,
	}
}

// NewRandomEdgeEnemy создаёт врага на случайном краю экрана
func NewRandomEdgeEnemy() *Enemy {
	rand.Seed(time.Now().UnixNano()) // Инициализируем генератор случайных чисел
	edge := rand.Intn(4)             // 0: верх, 1: право, 2: низ, 3: лево
	switch edge {
	case 0: // Верх
		return NewEnemy(float64(rand.Intn(1280)), 0)
	case 1: // Право
		return NewEnemy(1260, float64(rand.Intn(960)))
	case 2: // Низ
		return NewEnemy(float64(rand.Intn(1280)), 940)
	case 3: // Лево
		return NewEnemy(0, float64(rand.Intn(960)))
	default:
		return NewEnemy(0, 0) // На всякий случай
	}
}

func (e *Enemy) Update(targetX, targetY float64) {
	if !e.Alive {
		return
	}
	dx := targetX - e.X
	dy := targetY - e.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance > 0 {
		e.X += (dx / distance) * e.Speed
		e.Y += (dy / distance) * e.Speed
	}

	e.X = utils.Clamp(e.X, 0, 1260)
	e.Y = utils.Clamp(e.Y, 0, 940)
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	if e.Alive {
		ebitenutil.DrawRect(screen, e.X, e.Y, 20, 20, color.RGBA{255, 0, 0, 255})
	}
}

func (e *Enemy) GetHitbox() (x, y, width, height float64) {
    return e.X, e.Y, 20.0, 20.0 // Размеры врага
}