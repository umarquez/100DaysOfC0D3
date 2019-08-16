package main

type genericStructure struct {
	weight   float64
	strength float64
}

func (gs genericStructure) GetWeight() float64 {
	return gs.weight
}

func (gs genericStructure) GetStrength() float64 {
	return gs.strength
}

type WoodStructure struct {
	genericStructure
}

func NewWoodStruct() WoodStructure {
	return WoodStructure{
		genericStructure{
			weight:   50,
			strength: 100,
		},
	}
}

type IronStructure struct {
	genericStructure
}

func NewIronStructure() IronStructure {
	return IronStructure{
		genericStructure{
			weight:   100,
			strength: 500,
		},
	}
}

type CarbonFiberStructure struct {
	genericStructure
}

func NewCarbonFiberStructure() CarbonFiberStructure {
	return CarbonFiberStructure{
		genericStructure{
			weight:   100,
			strength: 1000,
		},
	}
}
