package particles

import (
	"math"
	"slices"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type Particle struct {
	Lifetime int64
	Speed    float64

	X float64
	Y float64
}

type NextPosition func(particle *Particle, deltaMS int64)
type ParticleRender func(row, col int, count [][]ParticleCellStat, p *ParticleParams) window.Cell
type Reset func(particle *Particle, params *ParticleParams)

type ParticleCellStat struct {
	count         int
	totalLifetime int64
}

type ParticleParams struct {
	MaxLife  int64
	MaxSpeed float64

	ParticleCount int

	XSTD float64
	X    int
	Y    int

	nextPosition NextPosition
	render       ParticleRender
	reset        Reset
}

type ParticleSystem struct {
	ParticleParams
	particles []*Particle

	lastTime int64
}

func NewParticleSystem(params ParticleParams) ParticleSystem {
	particles := make([]*Particle, 0)
	for i := 0; i < params.ParticleCount; i++ {
		particles = append(particles, &Particle{})
	}

	return ParticleSystem{
		ParticleParams: params,
		lastTime:       time.Now().UnixMilli(),
		particles:      particles,
	}
}

func (ps *ParticleSystem) Start() {
	for _, p := range ps.particles {
		ps.reset(p, &ps.ParticleParams)
	}
}

func (ps *ParticleSystem) Update() {
	now := time.Now().UnixMilli()
	delta := now - ps.lastTime
	ps.lastTime = now

	for _, p := range ps.particles {
		ps.nextPosition(p, delta)

		if p.Y >= float64(ps.Y) || p.X >= float64(ps.X) || p.Lifetime <= 0 {
			ps.reset(p, &ps.ParticleParams)
		}
	}
}

func (ps *ParticleSystem) Render() (window.Location, [][]window.Cell) {
	counts := make([][]ParticleCellStat, 0)

	for row := 0; row < ps.Y; row++ {
		count := make([]ParticleCellStat, 0)
		for col := 0; col < ps.X; col++ {
			count = append(count, ParticleCellStat{})
		}
		counts = append(counts, count)
	}

	for _, p := range ps.particles {
		row := int(math.Floor(p.Y))
		col := int(math.Round(p.X))

		counts[row][col].count++
		counts[row][col].totalLifetime += p.Lifetime
	}

	out := make([][]window.Cell, 0)
	for r, row := range counts {
		outRow := make([]window.Cell, 0)
		for c := range row {
			outRow = append(outRow, ps.render(r, c, counts, &ps.ParticleParams))
		}

		out = append(out, outRow)
	}

	slices.Reverse(out)
	return window.NewLocation(0, 0), out
}
