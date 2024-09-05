package persistence

import (
	"context"

	"github.com/neutrinocorp/geck/actuator"
	"github.com/neutrinocorp/geck/internal/reflection"
)

type ActuatorSQL struct {
	ClientSQL ClientSQL
}

var _ actuator.Actuator = (*ActuatorSQL)(nil)

func NewActuatorSQL(clientSQL ClientSQL) ActuatorSQL {
	return ActuatorSQL{
		ClientSQL: clientSQL,
	}
}

func (a ActuatorSQL) State(ctx context.Context) (actuator.State, error) {
	row := a.ClientSQL.QueryRowContext(ctx, "SELECT version()")
	if err := row.Err(); err != nil {
		return actuator.State{
			Status:      actuator.StatusDown,
			Description: err.Error(),
		}, nil
	}
	var version string
	if err := row.Scan(&version); err != nil {
		return actuator.State{
			Status:      actuator.StatusDown,
			Description: err.Error(),
		}, nil
	}

	return actuator.State{
		Status: actuator.StatusUp,
		Details: map[string]any{
			"driver_name": reflection.NewTypeFullNameAny(a.ClientSQL.Driver()),
			"version":     version,
		},
	}, nil
}
