package models

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

func NewPlayer(x, y int,
	width int,
	height int,
	frameX int,
	frameY int,
	cyclesX int,
	speed int,
	xMov int,
	yMov int) *Player {
	return &Player{x: x, y: y, width: width, height: height, frameX: frameX, frameY: frameY, cyclesX: cyclesX, speed: speed, xMov: xMov, yMov: yMov}

}

func (p *Player) X() int {
	return p.x
}

func (p *Player) SetX(x int) {
	p.x = x
}

func (p *Player) Y() int {
	return p.y
}

func (p *Player) SetY(y int) {
	p.y = y
}

func (p *Player) Width() int {
	return p.width
}

func (p *Player) SetWidth(width int) {
	p.width = width
}

func (p *Player) Height() int {
	return p.height
}

func (p *Player) SetHeight(height int) {
	p.height = height
}

func (p *Player) FrameX() int {
	return p.frameX
}

func (p *Player) SetFrameX(frameX int) {
	p.frameX = frameX
}

func (p *Player) FrameY() int {
	return p.frameY
}

func (p *Player) SetFrameY(frameY int) {
	p.frameY = frameY
}

func (p *Player) CyclesX() int {
	return p.cyclesX
}

func (p *Player) SetCyclesX(cyclesX int) {
	p.cyclesX = cyclesX
}

func (p *Player) Speed() int {
	return p.speed
}

func (p *Player) SetSpeed(speed int) {
	p.speed = speed
}

func (p *Player) XMov() int {
	return p.xMov
}

func (p *Player) SetXMov(xMov int) {
	p.xMov = xMov
}

func (p *Player) YMov() int {
	return p.yMov
}

func (p *Player) SetYMov(yMov int) {
	p.yMov = yMov
}
