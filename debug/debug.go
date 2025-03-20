package debug

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image/color"
    "fmt"
)

// Debug содержит глобальные настройки отладки
type Debug struct {
    Enabled        bool            // Флаг включения/выключения режима отладки
    ShowFPS        bool            // Показывать FPS
    ShowHitboxes   bool            // Показывать хитбоксы
    ShowPositions  bool            // Показывать координаты объектов
    DebugMessages  []string        // Список отладочных сообщений
}

// NewDebug создает новый экземпляр Debug
func NewDebug() *Debug {
    return &Debug{
        Enabled:       false,
        ShowFPS:       true,
        ShowHitboxes:  true,
        ShowPositions: true,
        DebugMessages: make([]string, 0),
    }
}

// Toggle переключает режим отладки
func (d *Debug) Toggle() {
    d.Enabled = !d.Enabled
}

// IsEnabled возвращает текущее состояние режима отладки
func (d *Debug) IsEnabled() bool {
    return d.Enabled
}

// DrawDebugInfo отрисовывает отладочную информацию
func (d *Debug) DrawDebugInfo(screen *ebiten.Image, tps float64) {
    if !d.Enabled {
        return
    }
    
    // Показываем FPS если включено
    if d.ShowFPS {
        ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 10, 10)
        ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f", tps), 10, 30)
    }
    
    // Показываем отладочные сообщения
    for i, msg := range d.DebugMessages {
        ebitenutil.DebugPrintAt(screen, msg, 10, 50+i*20)
    }
}

// DrawHitbox отрисовывает хитбокс объекта
func (d *Debug) DrawHitbox(screen *ebiten.Image, x, y, width, height float64) {
    if !d.ShowHitboxes {
        return
    }
    // Рисуем только контур
    ebitenutil.DrawLine(screen, x, y, x+width, y, color.RGBA{255, 0, 0, 255})
    ebitenutil.DrawLine(screen, x+width, y, x+width, y+height, color.RGBA{255, 0, 0, 255})
    ebitenutil.DrawLine(screen, x+width, y+height, x, y+height, color.RGBA{255, 0, 0, 255})
    ebitenutil.DrawLine(screen, x, y+height, x, y, color.RGBA{255, 0, 0, 255})
}

// AddMessage добавляет отладочное сообщение
func (d *Debug) AddMessage(msg string) {
    if !d.Enabled {
        return
    }
    
    // Добавляем сообщение в список
    d.DebugMessages = append(d.DebugMessages, msg)
    
    // Ограничиваем количество сообщений
    if len(d.DebugMessages) > 10 {
        d.DebugMessages = d.DebugMessages[len(d.DebugMessages)-10:]
    }
}

// ClearMessages очищает список отладочных сообщений
func (d *Debug) ClearMessages() {
    d.DebugMessages = make([]string, 0)
}