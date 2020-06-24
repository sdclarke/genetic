package genetic

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	batch *pixel.Batch
)

func MakeBatch() {
	batch = pixel.NewBatch(&pixel.TrianglesData{}, blackPic)
}

type Population struct {
	Brains       []*Brain
	minSteps     int
	obstacles    []pixel.Rect
	window       *pixelgl.Window
	staleness    int
	drawBest     bool
	generation   int
	goal         pixel.Vec
	mutationRate float64
}

func NewPopulation(size int, position pixel.Vec, moves int, bounds pixel.Rect, goal pixel.Vec, mutationRate float64, obstacles []pixel.Rect, window *pixelgl.Window) *Population {
	pop := &Population{
		minSteps:     moves,
		obstacles:    obstacles,
		window:       window,
		goal:         goal,
		mutationRate: mutationRate,
	}
	pop.Brains = make([]*Brain, size)
	for i := 0; i < size; i++ {
		pop.Brains[i] = NewBrain(position, moves, bounds)
	}
	return pop
}

func (p *Population) mutate() {
	for n, brain := range p.Brains {
		if n == 0 {
			continue
		}
		brain.Mutate(p.mutationRate)
	}
}

func (p *Population) NewGeneration() *Population {
	log.Printf("Staleness: %#v", p.staleness)
	p.calculateFitnesses()
	bestDotIndex := p.getBestDotIndex()
	newPopulation := &Population{
		minSteps:     p.minSteps,
		obstacles:    p.obstacles,
		window:       p.window,
		generation:   p.generation + 1,
		drawBest:     p.drawBest,
		goal:         p.goal,
		mutationRate: p.mutationRate,
	}
	if bestDotIndex == 0 {
		newPopulation.staleness = p.staleness + 1
	}
	newPopulation.Brains = make([]*Brain, len(p.Brains))
	fitnessSum := p.calculateFitnessSum()
	newPopulation.Brains[0] = p.Brains[bestDotIndex].Clone()
	newPopulation.Brains[0].SetBest(true)
	for i := 1; i < len(newPopulation.Brains); i++ {
		newPopulation.Brains[i] = p.selectParent(fitnessSum).Clone()
	}

	newPopulation.mutate()
	return newPopulation
}

func (p *Population) calculateFitnesses() {
	for _, brain := range p.Brains {
		brain.CalculateFitness(p.goal)
	}
}

func (p *Population) getBestDotIndex() int {
	maxFitness := 0.0
	maxIndex := 0
	for n, brain := range p.Brains {
		if brain.Fitness > maxFitness {
			maxFitness = brain.Fitness
			maxIndex = n
		}
	}
	if p.Brains[maxIndex].HasReachedGoal() {
		p.minSteps = p.Brains[maxIndex].NextMove
	}
	return maxIndex
}

func (p *Population) calculateFitnessSum() float64 {
	total := 0.0
	for _, brain := range p.Brains {
		total += brain.Fitness
	}
	return total
}

func (p *Population) selectParent(fitnessSum float64) *Brain {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Float64() * fitnessSum

	runningSum := 0.0

	for _, brain := range p.Brains {
		runningSum += brain.Fitness
		if runningSum > randomNum {
			return brain
		}
	}

	return nil
}

func (p *Population) Update() {
	batch.Clear()
	for _, brain := range p.Brains {
		if brain.NextMove > p.minSteps {
			brain.Kill()
		}
		matrix, err := brain.GetNextMove()
		if err != nil {
			var noMovesErr *NoMovesError
			if errors.As(err, &noMovesErr) {
				brain.Kill()
			} else {
				log.Fatalf("Unexpected error: %#v", err)
			}
		}
		brainPos := brain.GetPosition()
		x, y := brainPos.XY()
		if x < 0 || y < 0 || x > p.window.Bounds().W() || y > p.window.Bounds().H() {
			brain.SetPosition(pixel.V(pixel.Clamp(x, 0, p.window.Bounds().W()), pixel.Clamp(y, 0, p.window.Bounds().H())))
			matrix = brain.Matrix()
			brain.Kill()
		}
		if dist(brainPos, p.goal) < 10 {
			brain.SetReachedGoal(true)
		}
		for _, obstacle := range p.obstacles {
			if obstacle.Contains(brainPos) {
				brain.Kill()
			}
		}
		if best := brain.IsBest(); !p.drawBest || best {
			if best {
				brain.GetSprite().Draw(p.window, matrix)
			} else {
				brain.GetSprite().Draw(batch, matrix)
			}
		}
	}
	batch.Draw(p.window)
}

func (p *Population) AllDead() bool {
	for _, brain := range p.Brains {
		if !brain.IsDead() && !brain.HasReachedGoal() {
			return false
		}
	}
	return true
}

func (p *Population) SetDrawBest(b bool) {
	p.drawBest = b
}
