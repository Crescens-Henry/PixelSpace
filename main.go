package main

import (
	game2 "PixelSpace/game"
	"PixelSpace/models"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Game")

	background := load("./assets/map/Space Background.png")
	playerSprites := load("./assets/sprites/spriterickmorty.png")

	now := time.Now().UnixMilli()
	game := game2.NewGame(564,
		314,
		60,
		now,
		10)

	fpsInterval := int64(1000 / game.Fps())

	player := models.NewCharacter(100, 200, 40, 72, 0, 0, 4, 3, 0, 1, 2, 14, 0, 0)

	img := canvas.NewImageFromImage(background)
	img.FillMode = canvas.ImageFillOriginal

	sprite := image.NewRGBA(background.Bounds())

	playerImg := canvas.NewRasterFromImage(sprite)
	spriteSize := image.Pt(player.Width(), player.Height())

	c := container.New(layout.NewMaxLayout(), img, playerImg)
	w.SetContent(c)
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyDown:
			if player.Y() < int(game.CanvasHeight())-player.Height()-game.Margin() {
				player.SetYMov(player.Speed())
			}
			player.SetFrameY(player.DownY())
		case fyne.KeyUp:
			if player.Y() > 100 {
				player.SetYMov(-player.Speed())
			}
			player.SetFrameY(player.UpY())
		case fyne.KeyLeft:
			if player.X() > game.Margin() {
				player.SetXMov(-player.Speed())
			}
			player.SetFrameY(player.LeftY())
		case fyne.KeyRight:
			if player.X() < int(game.CanvasWidth())-player.Width()-game.Margin() {
				player.SetXMov(player.Speed())
			}
			player.SetFrameY(player.RightY())
		}
	})

	go func() {

		for {
			time.Sleep(time.Millisecond)

			now := time.Now().UnixMilli()
			elapsed := now - game.Then()

			if elapsed > fpsInterval {
				game.SetThen(now)

				spriteDP := image.Pt(player.Width()*player.FrameX(), player.Height()*player.FrameY())
				sr := image.Rectangle{spriteDP, spriteDP.Add(spriteSize)}

				dp := image.Pt(player.X(), player.Y())
				r := image.Rectangle{dp, dp.Add(spriteSize)}

				draw.Draw(sprite, sprite.Bounds(), image.Transparent, image.ZP, draw.Src)
				draw.Draw(sprite, r, playerSprites, sr.Min, draw.Src)
				playerImg = canvas.NewRasterFromImage(sprite)

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

	w.CenterOnScreen()
	w.ShowAndRun()
}
