// Package vector3 三维向量：(x,y,z)
package vector3

import (
	"math"
)

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (m *Vector3) Equal(v Vector3) bool {
	return m.X == v.X && m.Y == v.Y && m.Z == v.Z
}

// Set 三维向量：设值
func (m *Vector3) Set(x, y, z float64) {
	m.X = x
	m.Y = y
	m.Z = z
}

// Clone 三维向量：拷贝
func (m *Vector3) Clone() Vector3 {
	return NewVector3(m.X, m.Y, m.Z)
}

// Add 三维向量：加上
// this = this + v
func (m *Vector3) Add(v Vector3) {
	m.X += v.X
	m.Y += v.Y
	m.Z += v.Z
}

// Sub 三维向量：减去
// this = this - v
func (m *Vector3) Sub(v Vector3) {
	m.X -= v.X
	m.Y -= v.Y
	m.Z -= v.Z
}

// Multiply 三维向量：数乘
func (m *Vector3) Multiply(scalar float64) {
	m.X *= scalar
	m.Y *= scalar
	m.Z *= scalar
}

// Divide 三维向量：除法
func (m *Vector3) Divide(scalar float64) {
	if scalar == 0 {
		panic("分母不能为零！")
	}
	m.Multiply(1 / scalar)
}

// Dot 三维向量：点积
func (m *Vector3) Dot(v Vector3) float64 {
	return m.X*v.X + m.Y*v.Y + m.Z*v.Z
}

// Cross 三维向量：叉积
func (m *Vector3) Cross(v Vector3) {
	x, y, z := m.X, m.Y, m.Z
	m.X = y*v.Z - z*v.Y
	m.Y = z*v.X - x*v.Z
	m.Z = x*v.Y - y*v.X
}

// Length 三维向量：长度
func (m *Vector3) Length() float64 {
	return math.Sqrt(m.X*m.X + m.Y*m.Y + m.Z*m.Z)
}

// LengthSq 三维向量：长度平方
func (m *Vector3) LengthSq() float64 {
	return m.X*m.X + m.Y*m.Y + m.Z*m.Z
}

// Normalize 三维向量：单位化
func (m *Vector3) Normalize() {
	m.Divide(m.Length())
}

// NewVector3 返回：新向量
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

// Zero3 返回：零向量(0,0,0)
func Zero3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 0}
}

// XAxis3 X 轴 单位向量
func XAxis3() Vector3 {
	return Vector3{X: 1, Y: 0, Z: 0}
}

// YAxis3 Y 轴 单位向量
func YAxis3() Vector3 {
	return Vector3{X: 0, Y: 1, Z: 0}
}

// ZAxis3 Z 轴 单位向量
func ZAxis3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 1}
}
func XYAxis3() Vector3 {
	return Vector3{X: 1, Y: 1, Z: 0}
}
func XZAxis3() Vector3 {
	return Vector3{X: 1, Y: 0, Z: 1}
}
func YZAxis3() Vector3 {
	return Vector3{X: 0, Y: 1, Z: 1}
}
func XYZAxis3() Vector3 {
	return Vector3{X: 1, Y: 1, Z: 1}
}

// Add3 返回：a + b 向量
func Add3(a, b Vector3) Vector3 {
	return Vector3{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

// Sub3 返回：a - b 向量
func Sub3(a, b Vector3) Vector3 {
	return Vector3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

// Dot 返回：a * b 向量
func Dot(a, b Vector3) Vector3 {
	return Vector3{X: a.X * b.X, Y: a.Y * b.Y, Z: a.Z * b.Z}
}

// Cross3 返回：a X b 向量 (X 叉乘)
func Cross3(a, b Vector3) Vector3 {
	return Vector3{X: a.Y*b.Z - a.Z*b.Y, Y: a.Z*b.X - a.X*b.Z, Z: a.X*b.Y - a.Y*b.X}
}

func AddArray3(vs []Vector3, dv Vector3) []Vector3 {
	for i := range vs {
		vs[i].Add(dv)
	}
	return vs
}

func Multiply3(v Vector3, scalars []float64) []Vector3 {
	vs := make([]Vector3, 0)
	for _, value := range scalars {
		vector := v.Clone()
		vector.Multiply(value)
		vs = append(vs, vector)
	}
	return vs
}

// Normalize3 返回：单位化向量
func Normalize3(a Vector3) Vector3 {
	b := a.Clone()
	b.Normalize()
	return b
}

// GetDistance 求两点间距离
func GetDistance(a Vector3, b Vector3) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2) + math.Pow(a.Z-b.Z, 2))
}
