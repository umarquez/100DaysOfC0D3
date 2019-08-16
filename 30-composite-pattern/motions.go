package main

type genericMotion struct {
	speed           float64
	allowedTerrains []TerrainType
}

func (motion genericMotion) GetSpeed() float64 {
	return motion.speed
}

func (motion genericMotion) CanItGoOver(terrain TerrainType) bool {
	for _, t := range motion.allowedTerrains {
		if t == terrain {
			return true
		}
	}
	return false
}

type Wheels struct {
	genericMotion
}

func NewWheels() Wheels {
	return Wheels{genericMotion{
		speed: 10,
		allowedTerrains: []TerrainType{
			SandType,
			DirtType,
			MudType,
		},
	}}
}

type Turbine struct {
	genericMotion
}

func NewTurbine() Turbine {
	return Turbine{genericMotion{
		speed: 50,
		allowedTerrains: []TerrainType{
			WaterType,
		},
	}}
}

type TrackChain struct {
	genericMotion
}

func NewTrackChain() TrackChain {
	return TrackChain{genericMotion{
		speed: 8,
		allowedTerrains: []TerrainType{
			SandType,
			IceType,
			DirtType,
			MudType,
		},
	}}
}
