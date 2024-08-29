package reflection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/neutrinocorp/geck/actuator"
	"github.com/neutrinocorp/geck/reflection"
)

func TestNewTypeFullNameAny(t *testing.T) {
	out := reflection.NewTypeFullNameAny(actuator.DiskActuator{})
	assert.Equal(t, "actuator.DiskActuator", out)
	out = reflection.NewTypeFullName[actuator.DiskActuator]()
	assert.Equal(t, "actuator.DiskActuator", out)
}
