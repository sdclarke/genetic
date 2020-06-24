package genetic

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

var (
	blackSprite *pixel.Sprite
	greenSprite *pixel.Sprite
)

func MakeSprites() {
	pic := pixel.MakePictureData(pixel.R(0, 0, 5, 5))
	for n := range pic.Pix {
		pic.Pix[n] = colornames.Black
	}
	blackSprite = pixel.NewSprite(pic, pic.Bounds())

	pic = pixel.MakePictureData(pixel.R(0, 0, 20, 20))
	for n := range pic.Pix {
		pic.Pix[n] = colornames.Green
	}
	greenSprite = pixel.NewSprite(pic, pic.Bounds())
}

type Brain struct {
	position      pixel.Vec
	startPosition pixel.Vec
	velocity      pixel.Vec
	acceleration  pixel.Vec
	goal          pixel.Vec
	moves         []pixel.Vec
	NextMove      int
	dead          bool
	firstMove     bool
	windowBounds  pixel.Rect
	reachedGoal   bool
	Fitness       float64
	mutationRate  float64
	sprite        *pixel.Sprite
	best          bool
}

func NewBrain(position pixel.Vec, moves int, windowBounds pixel.Rect, goal pixel.Vec, mutationRate float64, best bool) *Brain {
	brain := &Brain{
		position:      position,
		startPosition: position,
		velocity:      pixel.V(0, 0),
		acceleration:  pixel.V(0, 0),
		goal:          goal,
		moves:         make([]pixel.Vec, moves),
		NextMove:      0,
		dead:          false,
		firstMove:     true,
		windowBounds:  windowBounds,
		mutationRate:  mutationRate,
		best:          best,
	}
	brain.sprite = blackSprite
	if best {
		brain.sprite = greenSprite
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < moves; i++ {
		brain.moves[i] = pixel.Unit(rand.Float64() * 2 * math.Pi)
	}
	return brain
}

func (b *Brain) GetNextMove() (pixel.Matrix, error) {
	if b.dead || b.reachedGoal {
		return b.matrix(), nil
	}
	if b.NextMove >= len(b.moves) {
		return b.matrix(), &NoMovesError{}
	}
	if b.firstMove {
		b.firstMove = false
		return b.matrix(), nil
	}
	b.acceleration = b.moves[b.NextMove]
	b.NextMove++
	b.velocity = b.velocity.Add(b.acceleration)

	if mag := b.velocity.X*b.velocity.X + b.velocity.Y*b.velocity.Y; mag > 25 {
		b.velocity.Scaled(5 / math.Sqrt(mag))
	}

	newPosition := b.position.Add(b.velocity)
	x, y := newPosition.XY()
	if x < 0 || y < 0 || x > b.windowBounds.W() || y > b.windowBounds.H() {
		b.position = pixel.V(pixel.Clamp(x, 0, b.windowBounds.W()), pixel.Clamp(y, 0, b.windowBounds.H()))
		return b.matrix(), &HitWallError{}
	}
	b.position = newPosition
	if dist(b.position, b.goal) < 10 {
		b.reachedGoal = true
	}
	return b.matrix(), nil
}

func (b *Brain) matrix() pixel.Matrix {
	mat := pixel.IM
	return mat.Moved(b.position)
}

func (b *Brain) Clone(best bool) *Brain {
	newBrain := NewBrain(b.startPosition, len(b.moves), b.windowBounds, b.goal, b.mutationRate, best)
	for n, move := range b.moves {
		newBrain.moves[n] = move
	}
	return newBrain
}

func (b *Brain) Kill() {
	b.dead = true
}

func (b *Brain) IsDead() bool {
	return b.dead
}

func (b *Brain) HasReachedGoal() bool {
	return b.reachedGoal
}

func (b *Brain) GetPosition() pixel.Vec {
	return b.position
}

func (b *Brain) Mutate() {
	rand.Seed(time.Now().UnixNano())
	for n, _ := range b.moves {
		if rand.Float64() < b.mutationRate {
			b.moves[n] = pixel.Unit(rand.Float64() * 2 * math.Pi)
		}
	}
}

func dist(a, b pixel.Vec) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (b *Brain) CalculateFitness() float64 {
	if b.reachedGoal {
		b.Fitness = float64(1.0)/16.0 + float64(10000.0)/float64(b.NextMove*b.NextMove)
	} else {
		distance := dist(b.position, b.goal)
		b.Fitness = 1.0 / (distance * distance)
	}
	return b.Fitness
}

func (b *Brain) GetSprite() *pixel.Sprite {
	return b.sprite
}

func (b *Brain) IsBest() bool {
	return b.best
}
