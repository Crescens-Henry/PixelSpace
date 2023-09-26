package game

type Game struct {
	canvasWidth  float32
	canvasHeight float32
	fps          int
	then         int64
	margin       int
}

func NewGame(canvasWidth float32, canvasHeight float32, fps int, then int64, margin int) *Game {
	return &Game{canvasWidth: canvasWidth, canvasHeight: canvasHeight, fps: fps, then: then, margin: margin}
}
func (g *Game) CanvasWidth() float32 {
	return g.canvasWidth
}

func (g *Game) SetCanvasWidth(canvasWidth float32) {
	g.canvasWidth = canvasWidth
}

func (g *Game) CanvasHeight() float32 {
	return g.canvasHeight
}

func (g *Game) SetCanvasHeight(canvasHeight float32) {
	g.canvasHeight = canvasHeight
}

func (g *Game) Fps() int {
	return g.fps
}

func (g *Game) SetFps(fps int) {
	g.fps = fps
}

func (g *Game) Then() int64 {
	return g.then
}

func (g *Game) SetThen(then int64) {
	g.then = then
}

func (g *Game) Margin() int {
	return g.margin
}

func (g *Game) SetMargin(margin int) {
	g.margin = margin
}
