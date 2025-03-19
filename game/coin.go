package game

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "log"
    "math/rand"
)

// Coin представляет монетку в игре
type Coin struct {
    x, y  float64
    image *ebiten.Image
}

// NewCoin создаёт новую монетку с случайной позицией
func NewCoin(screenWidth, screenHeight float64) *Coin {
    img, _, err := ebitenutil.NewImageFromFile("assets/coin.png")
    if err != nil {
        log.Fatal("Ошибка загрузки изображения монетки:", err)
    }
    return &Coin{
        x:     rand.Float64() * screenWidth,
        y:     rand.Float64() * screenHeight,
        image: img,
    }
}

// Draw отрисовывает монетку на экране
func (c *Coin) Draw(screen *ebiten.Image) {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(c.x, c.y)
    screen.DrawImage(c.image, op)
}

// Collides проверяет, пересекается ли монетка с заданным прямоугольником
func (c *Coin) Collides(x, y float64, width, height int) bool {
    coinWidth, coinHeight := c.image.Size()
    return x < c.x+float64(coinWidth) &&
           x+float64(width) > c.x &&
           y < c.y+float64(coinHeight) &&
           y+float64(height) > c.y
}