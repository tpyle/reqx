package context_test

import (
	"math"
	"testing"
	"text/template"

	"github.com/google/uuid"
	"github.com/tpyle/reqx/lib/requests/context"
)

func TestVariableSet_RenderTemplate(t *testing.T) {
	// Create a new VariableSet instance
	v := context.NewVariableSet()

	// Add some variables to the VariableSet
	v.SetVariable("name", "John")
	v.SetVariable("age", 30)

	// Create a new template
	tmp := template.New("test")
	tmp, err := tmp.Parse("My name is {{.name}} and I am {{.age}} years old.")
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Call the RenderTemplate method
	result, err := v.RenderTemplate(*tmp)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Check the result
	expected := "My name is John and I am 30 years old."
	if result != expected {
		t.Errorf("Unexpected result. Expected: %s, Got: %s", expected, result)
	}
}

func TestVariableSet_GenerateUUID(t *testing.T) {
	v := context.NewVariableSet()

	// Call the GenerateUUID method
	generatedUUID, err := v.GenerateUUID()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}

	// Check if the generated UUID is valid
	if generatedUUID == uuid.Nil {
		t.Error("Generated UUID is nil")
	}
}

func TestVariableSet_GenerateInt(t *testing.T) {
	v := context.NewVariableSet()

	// Call the GenerateInt method
	generatedInt, err := v.GenerateInt()
	if err != nil {
		t.Fatalf("Failed to generate int: %v", err)
	}

	if generatedInt == 0 {
		t.Error("Generated int is 0 (technically valid, but unlikely)")
	}

	// Check if the generated int is within the valid range
	if generatedInt < math.MinInt32 || generatedInt > math.MaxInt32 {
		t.Errorf("Generated int is out of range. Expected range: [%d, %d], Got: %d", math.MinInt32, math.MaxInt32, generatedInt)
	}
}
