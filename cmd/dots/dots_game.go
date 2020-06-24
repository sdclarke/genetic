package main

import (
	_ "image/png"
	"log"
	"os"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sdclarke/genetic"

	"golang.org/x/image/colornames"
)

func usage() {
	log.Fatalf("Usage: dots <population size> <mutation rate> <draw only best (true/false)>")
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 760),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	popSize := 500
	mutationRate := 0.01
	obstacleCount := 2
	drawOnlyBest := false
	for i, arg := range os.Args {
		if i == 1 {
			popSizeInt64, err := strconv.ParseInt(arg, 10, 0)
			if err != nil {
				usage()
			}
			popSize = int(popSizeInt64)
		} else if i == 2 {
			mutationRate, err = strconv.ParseFloat(arg, 64)
			if err != nil {
				usage()
			}
		} else if i == 3 {
			drawOnlyBest, err = strconv.ParseBool(arg)
			if err != nil {
				usage()
			}
		}
	}

	goal := pixel.V(win.Bounds().W()/2, win.Bounds().H()-10)
	obstacleLocations := make([]pixel.Vec, obstacleCount)
	obstacles := make([]pixel.Rect, obstacleCount)
	obstacleSprites := make([]*pixel.Sprite, obstacleCount)

	obstacleLocations[0] = pixel.V(win.Bounds().W()/4, win.Bounds().H()/1.5)
	obstacles[0] = pixel.R(0, 0, win.Bounds().W()/1.4, 26).Moved(obstacleLocations[0].Sub(pixel.V(win.Bounds().W()/2.8, 13)))
	obstaclePic := pixel.MakePictureData(obstacles[0])
	for n := range obstaclePic.Pix {
		obstaclePic.Pix[n] = colornames.Red
	}
	obstacleSprites[0] = pixel.NewSprite(obstaclePic, obstaclePic.Bounds())

	obstacleLocations[1] = pixel.V(3*win.Bounds().W()/4, win.Bounds().H()/3.3)
	obstacles[1] = pixel.R(0, 0, win.Bounds().W()/1.4, 26).Moved(obstacleLocations[1].Sub(pixel.V(win.Bounds().W()/2.8, 13)))
	obstaclePic = pixel.MakePictureData(obstacles[1])
	for n := range obstaclePic.Pix {
		obstaclePic.Pix[n] = colornames.Red
	}
	obstacleSprites[1] = pixel.NewSprite(obstaclePic, obstaclePic.Bounds())

	win.SetSmooth(true)
	goalPic := pixel.MakePictureData(pixel.R(0, 0, 20, 20))
	for n := range goalPic.Pix {
		goalPic.Pix[n] = colornames.Green
	}
	goalSprite := pixel.NewSprite(goalPic, goalPic.Bounds())

	population := genetic.NewPopulation(popSize, pixel.V(win.Bounds().W()/2, 10), 200, win.Bounds(), goal, mutationRate, obstacles, win)
	population.SetDrawBest(drawOnlyBest)

	for !win.Closed() {
		win.Clear(colornames.White)
		goalMat := pixel.IM
		goalMat = goalMat.Moved(goal)
		goalSprite.Draw(win, goalMat)
		for i := 0; i < obstacleCount; i++ {
			obstacleMat := pixel.IM
			obstacleMat = obstacleMat.Moved(obstacleLocations[i])
			obstacleSprites[i].Draw(win, obstacleMat)
		}
		population.Update()
		if population.AllDead() {
			population = population.NewGeneration()
		}

		win.Update()
	}
}

func main() {
	genetic.MakeSprites()
	genetic.MakeBatch()
	pixelgl.Run(run)
}
