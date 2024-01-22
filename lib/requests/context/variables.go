package context

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"text/template"

	"github.com/google/uuid"
)

type VariableSet map[string]interface{}

func NewVariableSet() *VariableSet {
	return &VariableSet{}
}

// Generates a random UUID
func (v *VariableSet) GenerateUUID() (uuid.UUID, error) {
	return uuid.NewRandom()
}

// Generates a random number between -2^31 and 2^31-1
func (v *VariableSet) GenerateInt() (int32, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<32))
	if err != nil {
		return 0, err
	}
	return int32(n.Int64() - (1 << 31)), nil
}

func (v *VariableSet) GetVariable(name string) interface{} {
	return (*v)[name]
}

func (v *VariableSet) SetVariable(name string, value interface{}) {
	(*v)[name] = value
}

func (v *VariableSet) DeleteVariable(name string) {
	delete(*v, name)
}

func (v *VariableSet) RenderTemplate(tmp template.Template) (string, error) {
	buf := &bytes.Buffer{}
	err := tmp.Execute(buf, v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
