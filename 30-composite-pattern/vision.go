package main

type genericVision struct {
	radius float64
	angle  float64
	vType  VisionType
}

func (vision genericVision) PointTo(x, y int) {}

func (vision genericVision) GetVisionRadius() float64 {
	return vision.radius
}

func (vision genericVision) GetVisionAngle() float64 {
	return vision.angle
}

func (vision genericVision) GetVisionType() VisionType {
	return vision.vType
}

type DaylightCam struct {
	genericVision
}

func NewDaylightCam() DaylightCam {
	return DaylightCam{genericVision{
		radius: 20,
		angle:  120,
		vType:  DaylightVision,
	}}
}

type NightVisionCam struct {
	genericVision
}

func NewNigthVisionCam() NightVisionCam {
	return NightVisionCam{genericVision{
		radius: 10,
		angle:  80,
		vType:  NightVision,
	}}
}
