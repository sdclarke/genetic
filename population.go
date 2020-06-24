package genetic

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Population struct {
	Brains     []*Brain
	minSteps   int
	obstacles  []pixel.Rect
	window     *pixelgl.Window
	staleness  int
	drawBest   bool
	generation int
}

func NewPopulation(size int, position pixel.Vec, moves int, bounds pixel.Rect, goal pixel.Vec, mutationRate float64, obstacles []pixel.Rect, window *pixelgl.Window) *Population {
	pop := &Population{
		minSteps:  moves,
		obstacles: obstacles,
		window:    window,
	}
	pop.Brains = make([]*Brain, size)
	for i := 0; i < size; i++ {
		pop.Brains[i] = NewBrain(position, moves, bounds, goal, mutationRate, false)
	}
	return pop
}

func (p *Population) mutate() {
	for n, brain := range p.Brains {
		if n == 0 {
			continue
		}
		brain.Mutate()
	}
}

func (p *Population) NewGeneration() *Population {
	log.Printf("Staleness: %#v", p.staleness)
	p.calculateFitnesses()
	bestDotIndex := p.getBestDotIndex()
	newPopulation := &Population{
		minSteps:   p.minSteps,
		obstacles:  p.obstacles,
		window:     p.window,
		generation: p.generation + 1,
		drawBest:   p.drawBest,
	}
	if bestDotIndex == 0 {
		newPopulation.staleness = p.staleness + 1
	}
	newPopulation.Brains = make([]*Brain, len(p.Brains))
	fitnessSum := p.calculateFitnessSum()
	newPopulation.Brains[0] = p.Brains[bestDotIndex].Clone(true)
	if p.staleness > 5 && bestDotIndex == 0 {
		for i := 1; i < len(newPopulation.Brains); i++ {
			newPopulation.Brains[i] = p.selectParent(fitnessSum, true).Clone(false)
		}
	} else {
		for i := 1; i < len(newPopulation.Brains); i++ {
			newPopulation.Brains[i] = p.selectParent(fitnessSum, false).Clone(false)
		}
	}

	newPopulation.mutate()
	return newPopulation
}

func (p *Population) calculateFitnesses() {
	for _, brain := range p.Brains {
		brain.CalculateFitness()
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

func (p *Population) selectParent(fitnessSum float64, stale bool) *Brain {
	if stale {
		return p.Brains[0]
	}
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
	for _, brain := range p.Brains {
		if brain.NextMove > p.minSteps {
			brain.Kill()
		}
		matrix, err := brain.GetNextMove()
		if err != nil {
			var noMovesErr *NoMovesError
			var hitWallErr *HitWallError
			if errors.As(err, &noMovesErr) {
				brain.Kill()
			} else if errors.As(err, &hitWallErr) {
				brain.Kill()
			} else {
				log.Fatalf("Unexpected error: %#v", err)
			}
		}
		for _, obstacle := range p.obstacles {
			if obstacle.Contains(brain.GetPosition()) {
				brain.Kill()
			}
		}
		if !p.drawBest || brain.IsBest() || p.generation == 0 {
			brain.GetSprite().Draw(p.window, matrix)
		}
	}
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
