package tinysa

type Model string

const (
	ModelBasic Model = "tinySA"
	ModelUltra Model = "tinySA4"
)

type deviceModel struct {
	model  Model
	width  int
	height int
}

var deviceModels = map[string]deviceModel{
	"tinySA":  {ModelBasic, 320, 280},
	"tinySA4": {ModelUltra, 480, 320},
}
