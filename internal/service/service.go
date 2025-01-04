package service

import (
	"math/rand"
	"time"

	"github.com/YattaDeSune/Game-of-life/pkg/life"
)

type LifeService struct {
	currentWorld *life.World
	nextWorld    *life.World
}

func New(height, width int) (*LifeService, error) {
	rand.NewSource(time.Now().UTC().UnixNano())

	currentWorld := life.NewWorld(height, width)

	currentWorld.Seed()

	newWorld := life.NewWorld(height, width)

	ls := LifeService{
		currentWorld: currentWorld,
		nextWorld:    newWorld,
	}

	return &ls, nil
}

// Новое состояние игры
func (ls *LifeService) NewState() *life.World {
	life.NextState(ls.currentWorld, ls.nextWorld)

	ls.currentWorld = ls.nextWorld

	return ls.currentWorld
}
