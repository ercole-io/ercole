package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffFeature(t *testing.T) {
	differences := DiffFeature([]Feature{
		Feature{Name: "Spatial queryes", Status: true},
		Feature{Name: "High heavy database", Status: false},
		Feature{Name: "Power saving", Status: false},
		Feature{Name: "Crackked", Status: true},
		Feature{Name: "Star wars support", Status: true},
		Feature{Name: "Star wars support SP3", Status: false},
	}, []Feature{
		Feature{Name: "Spatial queryes", Status: true},
		Feature{Name: "High heavy database", Status: true},
		Feature{Name: "Power saving", Status: false},
		Feature{Name: "Crackked", Status: false},
		Feature{Name: "Dark Engine", Status: true},
	})

	assert.Equal(t, differences["Spatial queryes"], DiffFeatureActive)
	assert.Equal(t, differences["High heavy database"], DiffFeatureActivated)
	assert.Equal(t, differences["Power saving"], DiffFeatureInactive)
	assert.Equal(t, differences["Crackked"], DiffFeatureDeactivated)
	assert.Equal(t, differences["Dark Engine"], DiffFeatureActivated)
	assert.Equal(t, differences["Star wars support"], DiffFeatureDeactivated)
	assert.Equal(t, differences["Star wars support SP3"], DiffFeatureInactive)
}
