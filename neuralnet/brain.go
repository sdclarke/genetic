package neuralnet

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
)

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
}

func NewBrain(position pixel.Vec, moves int, windowBounds pixel.Rect, goal pixel.Vec) *Brain {
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
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < moves; i++ {
		brain.moves[i] = pixel.Unit(rand.Float64() * 2 * math.Pi)
	}
	return brain
}

func (b *Brain) GetNextMove() (pixel.Vec, error) {
	if b.dead || b.reachedGoal {
		return b.position, nil
	}
	if b.NextMove >= len(b.moves) {
		return b.GetPosition(), &NoMovesError{}
	}
	if b.firstMove {
		b.firstMove = false
		return b.GetPosition(), nil
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
		return pixel.V(pixel.Clamp(x, 0, b.windowBounds.W()), pixel.Clamp(y, 0, b.windowBounds.H())), &HitWallError{}
	}
	b.position = newPosition
	if dist(b.position, b.goal) < 10 {
		b.reachedGoal = true
	}
	return b.GetPosition(), nil
}

func (b *Brain) Clone() *Brain {
	newBrain := NewBrain(b.startPosition, len(b.moves), b.windowBounds, b.goal)
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
		if rand.Float64() < 0.001 {
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
