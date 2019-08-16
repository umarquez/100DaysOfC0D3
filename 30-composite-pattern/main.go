package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
)

type Robot struct {
	StructureComponent
	MotionComponent
	VisionComponent
	WeaponComponent
	DefenceComponent
	SensorComponent
	Name    string
	LifeBar float64
}

func BuildRandomRobot() *Robot {
	r := new(Robot)
	r.LifeBar = float64(rand.Intn(150) + 50)

	switch rand.Intn(2) {
	case 0:
		r.StructureComponent = NewWoodStruct()
	case 1:
		r.StructureComponent = NewIronStructure()
	case 2:
		r.StructureComponent = NewCarbonFiberStructure()
	}

	switch rand.Intn(2) {
	case 0:
		r.MotionComponent = NewWheels()
	case 1:
		r.MotionComponent = NewTurbine()
	case 2:
		r.MotionComponent = NewTrackChain()
	}

	switch rand.Intn(1) {
	case 0:
		r.VisionComponent = NewDaylightCam()
	case 1:
		r.VisionComponent = NewNigthVisionCam()
	}

	switch rand.Intn(2) {
	case 0:
		r.WeaponComponent = NewLaser()
	case 1:
		r.WeaponComponent = NewGuns()
	case 2:
		r.WeaponComponent = NewBombs()
	}

	switch rand.Intn(2) {
	case 0:
		r.DefenceComponent = NewSteelShield()
	case 1:
		r.DefenceComponent = NewDiamondShield()
	case 2:
		r.DefenceComponent = NewCarbonFiberShield()
	}

	switch rand.Intn(1) {
	case 0:
		r.SensorComponent = NewMotionSensor()
	case 1:
		r.SensorComponent = NewProximitySensor()
	}
	return r
}

func PrintRobotCard(r *Robot) {
	cardData := make(map[string]interface{})
	cardData["Name"] = r.Name
	cardData["Weight"] = r.GetWeight()
	cardData["Strength"] = r.GetStrength()
	cardData["Speed"] = r.GetSpeed()

	var terrains []string
	if r.CanItGoOver(DirtType) {
		terrains = append(terrains, string(DirtType))
	}

	if r.CanItGoOver(WaterType) {
		terrains = append(terrains, string(WaterType))
	}

	if r.CanItGoOver(MudType) {
		terrains = append(terrains, string(MudType))
	}

	if r.CanItGoOver(SandType) {
		terrains = append(terrains, string(SandType))
	}

	if r.CanItGoOver(IceType) {
		terrains = append(terrains, string(IceType))
	}

	cardData["Terrains"] = terrains
	cardData["Vision"] = r.GetVisionType()
	cardData["Weapon"] = r.GetWeaponType()
	cardData["Shield"] = r.GetMaterial()
	cardData["Sensor"] = r.GetType()

	tplCard := `
=============[ {{.Name}} ]=============
- Peso: {{.Weight}}
- Resistencia: {{.Strength}}
- Velocidad: {{.Speed}}
- puede andar por:{{ range .Terrains }}
	- {{.}}{{ end }}
- Sistema de visión: {{.Vision}}
- Arma: {{.Weapon}}
- Escudo: {{.Shield}}
- Sensor: {{.Sensor}}
=======================================
`

	tpl := template.New("card")
	tpl, err := tpl.Parse(tplCard)
	if err != nil {
		log.Printf("error loading template, %v\n", err)
		return
	}

	err = tpl.Execute(os.Stdout, cardData)
	if err != nil {
		log.Printf("error executing template, %v\n", err)
		// return // descomentar si es neceario agregar cosas fuera del if
	}
}

func main() {
	for i := 5; i > 0; i-- {
		robot := BuildRandomRobot()
		robot.Name = fmt.Sprintf("HAL-%v", i)
		PrintRobotCard(robot)
	}
}
