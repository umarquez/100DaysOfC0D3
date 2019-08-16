package main

type StructureComponent interface {
	GetWeight() float64
	GetStrength() float64
}

type MotionComponent interface {
	GetSpeed() float64
	CanItGoOver(TerrainType) bool
}

type VisionComponent interface {
	PointTo(x, y int)
	GetVisionRadius() float64
	GetVisionAngle() float64
	GetVisionType() VisionType
}

type WeaponComponent interface {
	PointTo(x, y int)
	Fire() float64
	GetWeaponType() WeaponType
}

type DefenceComponent interface {
	GetMaterial() Material
	GetIntegrity() float64
}

type SensorComponent interface {
	GetType() SensorType
	GetRadius() float64
}
