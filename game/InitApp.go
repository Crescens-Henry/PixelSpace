package game

import (
	"PixelSpace/models"
	"fmt"
	"image"
	"image/draw"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func RunGame() {
	myApp := app.New()
	w := myApp.NewWindow("Game")

	// Variable para rastrear si se ha ganado el juego
	gameWon := false
	var win fyne.Window

	background := load("./assets/map/earth.png")
	playerSprites := load("./assets/sprites/spriterickmorty.png")

	fps := 40

	now := time.Now().UnixMilli()
	then := now

	player := models.NewPlayer(100,
		200,
		40,
		72,
		0,
		0,
		4,
		20,
		0,
		0)

	var mu sync.Mutex // Mutex para proteger la manipulación de los puntos
	points := []Point{}
	obstacles := []Obstacle{}
	bestScore := 0
	score := 0

	img := canvas.NewImageFromImage(background)
	img.FillMode = canvas.ImageFillOriginal

	sprite := image.NewRGBA(background.Bounds())

	playerImg := canvas.NewRasterFromImage(sprite)
	spriteSize := image.Pt(player.Width(), player.Height())

	c := container.New(layout.NewMaxLayout(), img, playerImg)
	w.SetContent(c)

	scoreLabel := widget.NewLabel("Puntaje: 0")
	content := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, scoreLabel),
		c, scoreLabel)
	w.SetContent(content)

	// Función para reiniciar el juego
	restartGame := func() {
		player.SetX(100)
		player.SetY(200)
		score = 0
		scoreLabel.SetText(fmt.Sprintf("Puntaje: %d", score))
		points = []Point{}
		obstacles = []Obstacle{}
		gameWon = false // Reiniciar la variable de juego ganado
	}

	// Función para mostrar la ventana de alerta personalizada
	showCustomAlert := func() {
		content := container.NewVBox(
			widget.NewLabel("¡Felicidades! Has ganado el juego."),
			widget.NewButton("Reiniciar", func() {
				restartGame()
				gameWon = false // Restablecer el juego ganado
				win.Hide()      // Oculta la ventana de alerta
			}),
		)

		win = myApp.NewWindow("¡Has Ganado!") // Asigna la ventana de alerta a la variable win
		win.SetContent(content)
		win.Resize(fyne.NewSize(400, 200))
		win.CenterOnScreen()
		win.Show()
		win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
			if k.Name == fyne.KeyEscape {
				win.Hide() // Oculta la ventana de alerta al presionar Esc
			}
		})
	}

	// Goroutine para la generación y manejo de obstáculos
	go func() {
		for {
			if !gameWon {
				// Generar 3 obstáculos en diferentes lugares del juego
				obstacles = []Obstacle{
					{100, 100, 40, 40},
					{300, 150, 40, 40},
					{500, 100, 40, 40},
					{600, 400, 40, 40},
					{200, 500, 40, 40},
					{400, 350, 40, 40},
					{700, 200, 40, 40},
				}
			}

			// Esperar un tiempo antes de generar obstáculos nuevamente
			time.Sleep(10 * time.Second)

			// Restablecer el juego solo si se ha ganado
			if gameWon {
				// Mostrar ventana de alerta personalizada
				win.Canvas().SetOnTypedKey(nil) // Desactivar controles
				showCustomAlert()
			}
		}
	}()

	// Función para verificar si se ha ganado el juego
	checkGameWon := func() {
		if score >= 20 && !gameWon {
			gameWon = true
			showCustomAlert() // Muestra la ventana de alerta cuando se gana
		}
	}

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
			time.Sleep(5 * time.Second)
		}
	}()

	// Función para actualizar y mostrar el mejor puntaje
	//updateBestScoreLabel := func(bestScoreLabel *widget.Label) {
	//	bestScoreLabel.SetText(fmt.Sprintf("Mejor Puntaje: %d", bestScore))
	//}

	go func() {
		for {
			time.Sleep(time.Millisecond)

			now := time.Now().UnixMilli()
			elapsed := now - then

			if elapsed > int64(1000/fps) {
				then = now

				spriteDP := image.Pt(player.Width()*player.FrameX(), player.Height()*player.FrameY())
				sr := image.Rectangle{spriteDP, spriteDP.Add(spriteSize)}

				dp := image.Pt(player.X(), player.Y())
				r := image.Rectangle{dp, dp.Add(spriteSize)}

				draw.Draw(sprite, sprite.Bounds(), image.Transparent, image.ZP, draw.Src)
				draw.Draw(sprite, r, playerSprites, sr.Min, draw.Src)
				playerImg = canvas.NewRasterFromImage(sprite)

				mu.Lock()
				dibujarPuntos(sprite, points, obstacles)

				// Verificar colisiones con los obstáculos
				for _, obstacle := range obstacles {
					if player.X() < obstacle.X+obstacle.Width &&
						player.X()+player.Width() > obstacle.X &&
						player.Y() < obstacle.Y+obstacle.Height &&
						player.Y()+player.Height() > obstacle.Y {
						// Colisión con un obstáculo, reiniciar el juego
						if score > bestScore {
							bestScore = score
							//updateBestScoreLabel(scoreLabel) // Actualizar el mejor puntaje en la interfaz
						}
						restartGame()
					}
				}

				// Verificar colisiones con los puntos
				// Dentro de la goroutine principal
				for i := 0; i < len(points); i++ {
					if !points[i].Collected {
						distancia := distanciaEntrePersonajes(*player, points[i]) // Utiliza la instancia correcta del jugador
						if distancia < distanciaMinima {
							points[i].Collected = true
							score++
							scoreLabel.SetText(fmt.Sprintf("Puntaje: %d", score))
							if score > bestScore {
								bestScore = score
								// Actualiza el mejor puntaje en la interfaz aquí si es necesario
							}
						}
					}
				}

				mu.Unlock()

				if player.XMov() != 0 || player.YMov() != 0 {
					player.SetX(player.X() + player.XMov())
					player.SetY(player.Y() + player.YMov())
					player.SetFrameX((player.FrameX() + 1) % player.CyclesX())
					player.SetXMov(0)
					player.SetYMov(0)
				} else {
					player.SetFrameX(0)
				}

				c.Refresh()
			}
		}
	}()

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyDown:
			if player.Y() < canvasHeight-player.Height()-margin {
				player.SetYMov(player.Speed())
				player.SetFrameX(0)
			}
		case fyne.KeyUp:
			if player.Y() > margin {
				player.SetYMov(-player.Speed())
				player.SetFrameY(3)
			}
		case fyne.KeyLeft:
			if player.X() > margin {
				player.SetXMov(-player.Speed())
				player.SetFrameY(1)
			}
		case fyne.KeyRight:
			if player.X() < canvasWidth-player.Width()-margin {
				player.SetXMov(player.Speed())
				player.SetFrameY(2)
			}
		}
		checkGameWon() // Verificar si se ha ganado después de cada movimiento
	})

	bestScoreLabel := widget.NewLabel(fmt.Sprintf("Sigue asi!!"))
	bestScoreContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, bestScoreLabel),
		content, bestScoreLabel)

	w.CenterOnScreen()
	w.SetContent(bestScoreContainer)
	w.ShowAndRun()
	wg.Wait()
}
