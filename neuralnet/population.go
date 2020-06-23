package neuralnet

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
)

type Population struct {
	Brains   []*Brain
	minSteps int
}

func NewPopulation(size int, position pixel.Vec, moves int, bounds pixel.Rect, goal pixel.Vec) *Population {
	pop := &Population{
		minSteps: moves,
	}
	pop.Brains = make([]*Brain, size)
	for i := 0; i < size; i++ {
		pop.Brains[i] = NewBrain(position, moves, bounds, goal)
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
	p.calculateFitnesses()
	bestDotIndex := p.getBestDotIndex()
	newPopulation := &Population{
		minSteps: p.minSteps,
	}
	newPopulation.Brains = make([]*Brain, len(p.Brains))
	fitnessSum := p.calculateFitnessSum()
	newPopulation.Brains[0] = p.Brains[bestDotIndex].Clone()
	for i := 1; i < len(newPopulation.Brains); i++ {
		newPopulation.Brains[i] = p.selectParent(fitnessSum).Clone()
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

func (p *Population) Update() []pixel.Vec {
	moves := make([]pixel.Vec, len(p.Brains))
	var err error
	for n, brain := range p.Brains {
		if brain.NextMove > p.minSteps {
			brain.Kill()
		}
		moves[n], err = brain.GetNextMove()
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
	}
	return moves
}

func (p *Population) AllDead() bool {
	for _, brain := range p.Brains {
		if !brain.IsDead() && !brain.HasReachedGoal() {
			return false
		}
	}
	return true
}
