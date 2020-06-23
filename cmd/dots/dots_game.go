package main

import (
	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sdclarke/genetic"

	"golang.org/x/image/colornames"
)

var (
	popSize = 500
	goal    pixel.Vec
)

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

	goal = pixel.V(win.Bounds().W()/2, win.Bounds().H()-10)

	win.SetSmooth(true)
	goalPic := pixel.MakePictureData(pixel.R(0, 0, 20, 20))
	for n, _ := range goalPic.Pix {
		goalPic.Pix[n] = colornames.Green
	}
	goalSprite := pixel.NewSprite(goalPic, goalPic.Bounds())

	pic := make([]*pixel.PictureData, popSize)
	sprite := make([]*pixel.Sprite, popSize)

	for i := 0; i < popSize; i++ {
		pic[i] = pixel.MakePictureData(pixel.R(0, 0, 5, 5))
		//rand.Seed(time.Now().UnixNano())
		//var col color.RGBA
		//if i != 0 {
		//col = color.RGBA{R: uint8(rand.Intn(255/i) * i), G: uint8(rand.Intn(255/i) * i), B: uint8(rand.Intn(255/i) * i), A: 255}
		//} else {
		col := color.RGBA{A: 255}
		//}
		for n, _ := range pic[i].Pix {
			pic[i].Pix[n] = col
		}

		sprite[i] = pixel.NewSprite(pic[i], pic[i].Bounds())
	}

	//population := genetic.NewPopulation(popSize, win.Bounds().Center(), 500, win.Bounds(), goal)
	population := genetic.NewPopulation(popSize, pixel.V(win.Bounds().W()/2, 10), 200, win.Bounds(), goal)
	moves := make([]pixel.Vec, popSize)

	//last := time.Now()
	for !win.Closed() {
		//dt := time.Since(last).Seconds()
		//last = time.Now()

		moves = population.Update()
		win.Clear(colornames.White)
		goalMat := pixel.IM
		goalMat = goalMat.Moved(goal)
		goalSprite.Draw(win, goalMat)
		for n, move := range moves {
			mat := make([]pixel.Matrix, len(population.Brains))
			mat[n] = pixel.IM
			mat[n] = mat[n].Moved(move)
			sprite[n].Draw(win, mat[n])
		}
		if population.AllDead() {
			population = population.NewGeneration()
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
