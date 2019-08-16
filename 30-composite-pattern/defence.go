package main

type genericDefence struct {
	material  Material
	integrity float64
}

func (defence genericDefence) GetMaterial() Material {
	return defence.material
}

func (defence genericDefence) GetIntegrity() float64 {
	return defence.integrity
}

type SteelShield struct {
	genericDefence
}

func NewSteelShield() SteelShield {
	return SteelShield{genericDefence{
		material:  SteelMaterial,
		integrity: 100,
	}}
}

type CarbonFiberShield struct {
	genericDefence
}

func NewCarbonFiberShield() CarbonFiberShield {
	return CarbonFiberShield{genericDefence{
		material:  CarbonFiberMaterial,
		integrity: 500,
	}}
}

type DiamondShield struct {
	genericDefence
}

func NewDiamondShield() DiamondShield {
	return DiamondShield{genericDefence{
		material:  DiamondMaterial,
		integrity: 1000,
	}}
}
