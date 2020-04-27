package main

import (
	"github.com/faiface/pixel"
	"github.com/google/uuid"
)

func newUUID() string {
	uid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return uid.String()
}

func calcDirectionAngle(goPositionVector, mousePositionVector pixel.Vec) float64 {
	// вычитаем из радиус-вектора положения мыши радиус-вектор положения
	// игрового объекта. Затем берем угол.
	return mousePositionVector.Sub(goPositionVector).Angle()
}
