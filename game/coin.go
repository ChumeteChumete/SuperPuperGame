package game

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image"
    "log"
    "math/rand"
    "time"
)

// Coin представляет монетку в игре
type Coin struct {
    x, y         float64
    spriteSheet  *ebiten.Image
    frameWidth   int
    frameHeight  int
    frameCount   int
    currentFrame int
    animSpeed    float64
    lastUpdated  time.Time
}

// NewCoin создаёт новую монетку с случайной позицией
func NewCoin(screenWidth, screenHeight float64) *Coin {
    img, _, err := ebitenutil.NewImageFromFile("assets/coin.png")
    if err != nil {
        log.Fatal("Ошибка загрузки изображения монетки:", err)
    }
    
    // Получаем размер изображения
    width, height := img.Size()
    
    // Определяем количество кадров в спрайт-листе
    // Предполагаем, что спрайты расположены в одну линию горизонтально
    frameCount := 15 // Предполагаемое количество кадров, измените если необходимо
    frameWidth := width / frameCount
    
    return &Coin{
        x:           rand.Float64() * screenWidth,
        y:           rand.Float64() * screenHeight,
        spriteSheet: img,
        frameWidth:  frameWidth,
        frameHeight: height,
        frameCount:  frameCount,
        currentFrame: 0,
        animSpeed:   0.15, // Скорость анимации - количество секунд между кадрами
        lastUpdated: time.Now(),
    }
}

// Update обновляет состояние монетки, включая анимацию
func (c *Coin) Update() {
    now := time.Now()
    
    // Обновляем анимацию по таймеру
    if now.Sub(c.lastUpdated).Seconds() >= c.animSpeed {
        c.currentFrame = (c.currentFrame + 1) % c.frameCount
        c.lastUpdated = now
    }
}

// Draw отрисовывает текущий кадр монетки на экране
func (c *Coin) Draw(screen *ebiten.Image) {
    op := &ebiten.DrawImageOptions{}
    
    // Вычисляем позицию текущего кадра в спрайт-листе
    sx := c.frameWidth * c.currentFrame
    sy := 0
    
    // Вырезаем текущий кадр из спрайт-листа
    // Используем image.Rect вместо ebiten.Rect
    r := image.Rect(sx, sy, sx+c.frameWidth, sy+c.frameHeight)
    subImage := c.spriteSheet.SubImage(r).(*ebiten.Image)
    
    // Позиционируем кадр в игровом мире
    op.GeoM.Translate(c.x, c.y)
    
    // Отрисовываем текущий кадр
    screen.DrawImage(subImage, op)
}

// Collides проверяет, пересекается ли монетка с заданным прямоугольником
func (c *Coin) Collides(x, y float64, width, height int) bool {
    return x < c.x+float64(c.frameWidth) &&
           x+float64(width) > c.x &&
           y < c.y+float64(c.frameHeight) &&
           y+float64(height) > c.y
}