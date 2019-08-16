package main

type genericWeapon struct {
	damage float64
	wType  WeaponType
}

func (weapon genericWeapon) PointTo(x, y int) {}

func (weapon genericWeapon) Fire() float64 {
	return weapon.damage
}

func (weapon genericWeapon) GetWeaponType() WeaponType {
	return weapon.wType
}

type Laser struct {
	genericWeapon
}

func NewLaser() Laser {
	return Laser{genericWeapon{
		damage: 15,
		wType:  LaserWeapon,
	}}
}

type Bombs struct {
	genericWeapon
}

func NewBombs() Bombs {
	return Bombs{genericWeapon{
		damage: 8,
		wType:  BombWeapon,
	}}
}

type Guns struct {
	genericWeapon
}

func NewGuns() Guns {
	return Guns{genericWeapon{
		damage: 2,
		wType:  GunWeapon,
	}}
}
