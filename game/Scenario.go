package game

import (
	"PixelSpace/models"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
)

const (
	distanciaMinima = 20  // Distancia mínima para recolectar un punto
	velocidadMax    = 3.0 // Velocidad máxima del personaje
	numeroPuntos    = 10  // Número de puntos fijos
	puntoSize       = 10  // Tamaño de los puntos
	canvasWidth     = 800
	canvasHeight    = 650
	margin          = 10
)

func load(filePath string) image.Image {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot read file:", err)
	}

	imgData, err := png.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
	}
	return imgData.(image.Image)
}

func distanciaEntrePersonajes(p1 models.Player, p2 Point) float64 {
	deltaX := float64(p1.X() - p2.X)
	deltaY := float64(p1.Y() - p2.Y)
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

func generarPuntos(width, height int, wg *sync.WaitGroup, mu *sync.Mutex) []Point {
	defer wg.Done()
	points := make([]Point, numeroPuntos)
	for i := 0; i < numeroPuntos; i++ {
		points[i] = Point{
			X:         rand.Intn(width),
			Y:         rand.Intn(height),
			Collected: false,
		}
	}
	return points // Devuelve la lista de puntos generados
}

func dibujarPuntos(img draw.Image, points []Point, obstacles []Obstacle) {
	for _, obstacle := range obstacles {
		for x := obstacle.X; x < obstacle.X+obstacle.Width; x++ {
			for y := obstacle.Y; y < obstacle.Y+obstacle.Height; y++ {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, color.RGBA{0, 0, 255, 255}) // Color azul para los obstáculos
				}
			}
		}
	}

	for _, p := range points {
		if !p.Collected {
			// Verificar si el punto no ha sido recolectado antes de dibujarlo en rojo
			for x := p.X - puntoSize/2; x <= p.X+puntoSize/2; x++ {
				for y := p.Y - puntoSize/2; y <= p.Y+puntoSize/2; y++ {
					if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
						img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Color rojo para los puntos no recogidos
					}
				}
			}
		}
	}
}
