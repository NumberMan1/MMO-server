package core

import (
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"math"
)

// LookRotation 方向向量转欧拉角
func LookRotation(fromDir *vector3.Vector3) *vector3.Vector3 {
	rad2Deg := 57.29578
	eulerAngles := vector3.NewVector3(0, 0, 0)
	eulerAngles.X = math.Acos(math.Sqrt((fromDir.X*fromDir.X+fromDir.Z*fromDir.Z)/(fromDir.X*fromDir.X+fromDir.Y*fromDir.Y+fromDir.Z*fromDir.Z))) * rad2Deg
	if fromDir.Y > 0 {
		eulerAngles.X = 360 - eulerAngles.X
	}
	//AngleY = arc tan(x/z)
	eulerAngles.Y = math.Atan2(fromDir.X, fromDir.Z) * rad2Deg
	if eulerAngles.Y < 0 {
		eulerAngles.Y += 180
	}
	if fromDir.X < 0 {
		eulerAngles.Y += 180
	}
	//AngleZ = 0
	eulerAngles.Z = 0
	return eulerAngles
}
