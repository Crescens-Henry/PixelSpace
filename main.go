package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	distanciaMinima = 30  // Distancia mínima para recolectar un punto
	velocidadMax    = 3.0 // Velocidad máxima del personaje
	numeroPuntos    = 10  // Número de puntos fijos
	puntoSize       = 10  // Tamaño de los puntos
	canvasWidth     = 800
	canvasHeight    = 650
	margin          = 10
)

type Player struct {
	x, y    int
	width   int
	height  int
	frameX  int
	frameY  int
	cyclesX int
	speed   int
	xMov    int
	yMov    int
}

type Point struct {
	X, Y      int
	Collected bool
}

type Obstacle struct {
	X, Y   int
	Width  int
	Height int
}

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

func distanciaEntrePersonajes(p1 Player, p2 Point) float64 {
	deltaX := float64(p1.x - p2.X)
	deltaY := float64(p1.y - p2.Y)
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
	mu.Lock()
	defer mu.Unlock()
	return points
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

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Game")

	background := load("./assets/map/118804.png")
	playerSprites := load("./assets/sprites/spriterickmorty.png")

	fps := 30

	now := time.Now().UnixMilli()
	then := now

	player := Player{
		x:       100,
		y:       200,
		width:   40,
		height:  68,
		frameX:  0,
		frameY:  2, // Cambiar el frame inicial para que mire hacia abajo
		cyclesX: 4,
		speed:   20,
		xMov:    0,
		yMov:    0,
	}

	var mu sync.Mutex // Mutex para proteger la manipulación de los puntos
	points := []Point{}
	obstacles := []Obstacle{}
	bestScore := 0
	score := 0

	img := canvas.NewImageFromImage(background)
	img.FillMode = canvas.ImageFillOriginal

	sprite := image.NewRGBA(background.Bounds())

	playerImg := canvas.NewRasterFromImage(sprite)
	spriteSize := image.Pt(player.width, player.height)

	c := container.New(layout.NewMaxLayout(), img, playerImg)
	w.SetContent(c)

	scoreLabel := widget.NewLabel("Puntaje: 0")
	content := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, scoreLabel),
		c, scoreLabel)
	w.SetContent(content)

	// Función para reiniciar el juego
	restartGame := func() {
		player.x = 100
		player.y = 200
		score = 0
		scoreLabel.SetText(fmt.Sprintf("Puntaje: %d", score))
		points = []Point{}
		obstacles = []Obstacle{}
	}

	// Goroutine para la generación y manejo de obstáculos
	go func() {
		for {
			// Generar 3 obstáculos en la parte superior del juego
			obstacles = []Obstacle{
				{100, 100, 40, 40},
				{300, 150, 40, 40},
				{500, 100, 40, 40},
				{200, 300, 40, 40},
				{400, 350, 40, 40},
				{600, 400, 40, 40},
				{250, 500, 40, 40},
				{350, 200, 40, 40},
				{550, 250, 40, 40},
				{700, 500, 40, 40},
				{100, 500, 40, 40},
				{700, 100, 40, 40},
			}

			// Esperar un tiempo antes de generar obstáculos nuevamente
			time.Sleep(10 * time.Second)
		}
	}()

	// Goroutine para la generación de puntos
	var wg sync.WaitGroup
	go func() {
		for {
			wg.Add(1)
			newPoints := generarPuntos(canvasWidth, canvasHeight, &wg, &mu)
			mu.Lock()
			points = append(points, newPoints...)
			mu.Unlock()

			// Esperar un tiempo antes de generar más puntos
			time.Sleep(15 * time.Second)
		}
	}()

	// Función para actualizar y mostrar el mejor puntaje
	updateBestScoreLabel := func(bestScoreLabel *widget.Label) {
		bestScoreLabel.SetText(fmt.Sprintf("Mejor Puntaje: %d", bestScore))
	}

	go func() {
		for {
			time.Sleep(time.Millisecond)

			now := time.Now().UnixMilli()
			elapsed := now - then

			if elapsed > int64(1000/fps) {
				then = now

				spriteDP := image.Pt(player.width*player.frameX, player.height*player.frameY)
				sr := image.Rectangle{spriteDP, spriteDP.Add(spriteSize)}

				dp := image.Pt(player.x, player.y)
				r := image.Rectangle{dp, dp.Add(spriteSize)}

				draw.Draw(sprite, sprite.Bounds(), image.Transparent, image.ZP, draw.Src)
				draw.Draw(sprite, r, playerSprites, sr.Min, draw.Src)
				playerImg = canvas.NewRasterFromImage(sprite)

				mu.Lock()
				dibujarPuntos(sprite, points, obstacles)

				// Verificar colisiones con los obstáculos
				for _, obstacle := range obstacles {
					if player.x < obstacle.X+obstacle.Width &&
						player.x+player.width > obstacle.X &&
						player.y < obstacle.Y+obstacle.Height &&
						player.y+player.height > obstacle.Y {
						// Colisión con un obstáculo, reiniciar el juego
						if score > bestScore {
							bestScore = score
							updateBestScoreLabel(scoreLabel) // Actualizar el mejor puntaje en la interfaz
						}
						restartGame()
					}
				}

				// Verificar colisiones con los puntos
				for i := 0; i < len(points); i++ {
					if !points[i].Collected {
						distancia := distanciaEntrePersonajes(player, points[i])
						if distancia < distanciaMinima {
							points[i].Collected = true
							score++
							scoreLabel.SetText(fmt.Sprintf("Puntaje: %d", score))
							if score > bestScore {
								bestScore = score
								updateBestScoreLabel(scoreLabel) // Actualizar el mejor puntaje en la interfaz
							}
						}
					}
				}
				mu.Unlock()

				if player.xMov != 0 || player.yMov != 0 {
					player.x += player.xMov
					player.y += player.yMov
					player.frameX = (player.frameX + 1) % player.cyclesX
					player.xMov = 0
					player.yMov = 0
				} else {
					player.frameX = 0
				}

				c.Refresh()
			}
		}
	}()

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyDown:
			if player.y < canvasHeight-player.height-margin {
				player.yMov = player.speed
				player.frameY = 2
			}
		case fyne.KeyUp:
			if player.y > margin {
				player.yMov = -player.speed
				player.frameY = 0
			}
		case fyne.KeyLeft:
			if player.x > margin {
				player.xMov = -player.speed
				player.frameY = 1
			}
		case fyne.KeyRight:
			if player.x < canvasWidth-player.width-margin {
				player.xMov = player.speed
				player.frameY = 3
			}
		}
	})

	bestScoreLabel := widget.NewLabel(fmt.Sprintf("Mejor Puntaje: %d", bestScore))
	bestScoreContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, bestScoreLabel),
		content, bestScoreLabel)

	w.CenterOnScreen()
	w.SetContent(bestScoreContainer)
	w.ShowAndRun()
	wg.Wait()
}
