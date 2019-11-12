package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabasesArrayAsMap(t *testing.T) {
	dbs := []Database{
		Database{
			Name:     "superdb",
			Platform: "HPP 1.2.3",
		},
		Database{
			Name:     "superdb2",
			Platform: "PPH 1.2.3",
		},
	}

	dbsMap := DatabasesArrayAsMap(dbs)
	assert.Len(t, dbsMap, 2)
	assert.Equal(t, Database{
		Name:     "superdb",
		Platform: "HPP 1.2.3",
	}, dbsMap["superdb"])
	assert.Equal(t, Database{
		Name:     "superdb2",
		Platform: "PPH 1.2.3",
	}, dbsMap["superdb2"])
}

func TestHasEnterpriseLicense_NoEnterpriseLicense(t *testing.T) {
	assert.False(t, HasEnterpriseLicense(Database{
		Name: "superdb",
		Licenses: []License{
			License{Name: "Driving", Count: 10},
			License{Name: "Illegal query engine", Count: 9},
		},
	}))
}

func TestHasEnterpriseLicense_WithEnterpriseLicense(t *testing.T) {
	assert.True(t, HasEnterpriseLicense(Database{
		Name: "superdb",
		Licenses: []License{
			License{Name: "Oracle ENT", Count: 10},
			License{Name: "Illegal query engine", Count: 9},
		},
	}))
}

func TestHasEnterpriseLicense_NilLicense(t *testing.T) {
	assert.False(t, HasEnterpriseLicense(Database{
		Name:     "superdb",
		Licenses: nil,
	}))
}
