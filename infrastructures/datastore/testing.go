package datastore

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jblim0125/redredis-expire/common/appdata"
)

// MakeTestDataStore creates a DataStore for use with unit tests.
func (DataStore) MakeTestDataStore(tb testing.TB, log *logrus.Logger) *DataStore {
	ds := makeUnmigratedTestSQLStore(tb, log)
	err := ds.Migrate()
	require.NoError(tb, err)
	return ds
}

func makeUnmigratedTestSQLStore(tb testing.TB, log *logrus.Logger) *DataStore {
	ds, _ := DataStore{}.New(log)
	conf := &appdata.DatastoreConfiguration{
		Database: appdata.Sqlite,
		Endpoint: appdata.EndpointInfo{
			Path: "file:test.db?cache=shared&mode=memory",
		},
		Debug: appdata.DatastoreDebug{
			LogLevel:      "info",
			SlowThreshold: "30sec",
		},
	}
	ds.Connect(conf)
	return ds
}

// CloseConnection closes underlying database connection.
func CloseConnection(tb testing.TB, ds *DataStore) {
	db, err := ds.Orm.DB()
	assert.NoError(tb, err)
	err = db.Close()
	assert.NoError(tb, err)
}
