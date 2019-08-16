package main

type genericSensor struct {
	cType   SensorType
	cRadius float64
}

func (sensor genericSensor) GetType() SensorType {
	return sensor.cType
}

func (sensor genericSensor) GetRadius() float64 {
	return sensor.cRadius
}

type MotionSensor struct {
	genericSensor
}

func NewMotionSensor() MotionSensor {
	return MotionSensor{genericSensor{
		cType:   MotionType,
		cRadius: 5,
	}}
}

type ProximitySensor struct {
	genericSensor
}

func NewProximitySensor() ProximitySensor {
	return ProximitySensor{genericSensor{
		cType:   ProximityType,
		cRadius: 10,
	}}
}
