package particles

import (
	"math"
	"math/rand"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type Coffee struct {
	ParticleSystem
	id int
	z  int
}

func (c *Coffee) Id() int {
	return c.id
}

func (c *Coffee) Z() int {
	return c.z
}

func reset(p *Particle, params *ParticleParams) {
	p.Lifetime = int64(math.Floor(float64(params.MaxLife) * rand.Float64()))
	p.Speed = params.MaxSpeed * rand.Float64()

	maxX := math.Floor(float64(params.X) / 2)
	x := math.Max(-maxX, math.Min(rand.NormFloat64()*params.XSTD, maxX))
	p.X = x + maxX
	p.Y = 0
}

func nextPos(particle *Particle, deltaMS int64) {
	particle.Lifetime -= deltaMS
	if particle.Lifetime <= 0 {
		return
	}

	percent := (float64(deltaMS) / 1000.0)
	particle.Y += particle.Speed * percent
}

var dirs = [][]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},

	{0, -1},
	{0, 1},

	{1, 0},
	{1, 1},
	{1, -1},
}

func countParticles(row, col int, counts [][]ParticleCellStat) int {
	count := 0
	for _, dir := range dirs {
		r := row + dir[0]
		c := col + dir[1]

		if r < 0 || r >= len(counts) || c < 0 || c >= len(counts[0]) {
			continue
		}
		count += counts[row+dir[0]][col+dir[1]].count
	}
	return count
}

func NewCoffee(width, height int, scale float64) Coffee {
	assert.Assert(width%2 == 1, "width of particle system MUST be odd")

	startTime := time.Now().UnixMilli()
	ascii := func(row, col int, counts [][]ParticleCellStat, params *ParticleParams) window.Cell {
		stats := counts[row][col]
		if stats.count == 0 {
			return window.DefaultCell(' ')
		}

		direction := row +
			int(((time.Now().UnixMilli()-startTime)/2000)%2)

		/**
		  white = FF FF FF
		  yellow = FF FF 00
		  orange = FF 88 00
		  red = FF 00 00
		*/

		normalLife := math.Min(0.9999, math.Max(0.0001, float64(stats.totalLifetime)/float64(params.MaxLife)))
		halfCol := float64(params.X / 2)

		colDist := math.Abs(halfCol - float64(col))
		//rowSq := float64(row * row)
		colNormal := math.Max(0, 1-colDist/halfCol)

		normal := normalLife * colNormal
		green := byte(255 * normal)
		color := window.NewColor(255, green, 0, true)

		if countParticles(row, col, counts) > 3 {
			if direction%2 == 0 {
				return window.ForegroundCell('{', color)
			}
			return window.ForegroundCell('}', color)
		}
		return window.ForegroundCell('.', color)
	}

	_ = ascii

	/*
	   asciiFire := func(row, col int, counts [][]int) window.Cell {
	       count := counts[row][col]
	       if count == 0 {
	           return window.DefaultCell(' ')
	       }
	       if count < 4 {
	           return window.DefaultCell('░')
	       }
	       if count  < 6 {
	           return "▒"
	       }
	       if count < 9 {
	           return "▓"
	       }
	       return "█"
	   }
	*/

	return Coffee{
		ParticleSystem: NewParticleSystem(
			ParticleParams{
				MaxLife:       2 * 6000,
				MaxSpeed:      1.75 * 3,
				ParticleCount: 1000,

				reset:        reset,
				render:       ascii,
				nextPosition: nextPos,

				XSTD: scale,
				X:    width,
				Y:    height,
			},
		),
		z:  1,
		id: window.GetNextId(),
	}
}
