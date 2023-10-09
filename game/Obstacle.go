package game

// Obstacle representa un obstáculo en el juego
type Obstacle struct {
	X      int // Coordenada X del obstáculo
	Y      int // Coordenada Y del obstáculo
	Width  int // Ancho del obstáculo
	Height int // Altura del obstáculo
}

// NewObstacle crea un nuevo obstáculo con valores iniciales
func NewObstacle(x, y, width, height int) *Obstacle {
	return &Obstacle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Puedes definir otros métodos relacionados con los obstáculos aquí si es necesario
