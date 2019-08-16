package main

type TerrainType string

const (
	DirtType  TerrainType = "dirt"
	SandType  TerrainType = "sand"
	WaterType TerrainType = "water"
	IceType   TerrainType = "ice"
	MudType   TerrainType = "mud"
)

type VisionType string

const (
	DaylightVision VisionType = "daylight"
	NightVision    VisionType = "night"
	XRayVision     VisionType = "x-ray"
	InfraredVision VisionType = "infrared"
	SonarVision    VisionType = "sonar"
)

type Material string

const (
	WoodMaterial        Material = "wood"
	SteelMaterial       Material = "steel"
	DiamondMaterial     Material = "diamond"
	CarbonFiberMaterial Material = "carbon-fiber"
	AdamantiumMaterial  Material = "adamantium"
)

type SensorType string

const (
	ProximityType SensorType = "proximity"
	MotionType    SensorType = "motion"
	HeatType      SensorType = "heat"
	SoundType     SensorType = "sound"
	RadioFrqType  SensorType = "rf"
)

type WeaponType string

const LaserWeapon WeaponType = "laser"
const GunWeapon WeaponType = "gun"
const BombWeapon WeaponType = "bomb"
const HammerWeapon WeaponType = "hammer"
const MissileWeapon WeaponType = "missile"
