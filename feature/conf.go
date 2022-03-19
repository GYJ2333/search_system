package feature

type Conf struct {
	Kind []Kind `yaml:"kind"`
}
type Kind struct {
	Name    string   `yaml:"name"`
	Feature []string `yaml:"feature"`
}