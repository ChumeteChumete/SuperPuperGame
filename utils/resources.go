// Пакет utils содержит вспомогательные функции
package utils

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	_ "image/png" // Для поддержки PNG изображений
	"log"
)

// ImageCache - кэш загруженных изображений
var ImageCache = make(map[string]*ebiten.Image)

// LoadImage загружает изображение из файла и кэширует его
func LoadImage(path string) *ebiten.Image {
	// Проверяем, есть ли изображение в кэше
	if img, ok := ImageCache[path]; ok {
		return img
	}
	
	// Загружаем изображение из файла
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("Ошибка загрузки изображения %s: %v", path, err)
		// Возвращаем пустое изображение в случае ошибки
		emptyImg := ebiten.NewImage(10, 10)
		return emptyImg
	}
	
	// Кэшируем изображение
	ImageCache[path] = img
	
	return img
}

// LoadEmbeddedImage загружает встроенное изображение из байтового массива
func LoadEmbeddedImage(data []byte) *ebiten.Image {
	// Декодируем изображение из байтов
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("Ошибка декодирования встроенного изображения: %v", err)
		// Возвращаем пустое изображение в случае ошибки
		emptyImg := ebiten.NewImage(10, 10)
		return emptyImg
	}
	
	// Преобразуем в изображение Ebiten
	ebitenImg := ebiten.NewImageFromImage(img)
	
	return ebitenImg
}

// GetImageDimensions возвращает размеры изображения
func GetImageDimensions(img *ebiten.Image) (int, int) {
	return img.Bounds().Dx(), img.Bounds().Dy()
}