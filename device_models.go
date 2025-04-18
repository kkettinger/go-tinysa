package tinysa

// Model represents the type of device model as a string.
type Model string

// Supported device models.
const (
	ModelBasic Model = "tinySA"
	ModelUltra Model = "tinySA4"
)

// deviceModel holds metadata for a specific device model.
type deviceModel struct {
	model  Model
	width  int
	height int
}

// deviceModels maps model names to their corresponding deviceModel configurations.
var deviceModels = map[string]deviceModel{
	"tinySA":  {ModelBasic, 320, 280},
	"tinySA4": {ModelUltra, 480, 320},
}
