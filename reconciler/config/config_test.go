// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"
	"time"

	"github.com/agntcy/dir/reconciler/tasks/policy"
	dbconfig "github.com/agntcy/dir/server/database/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_NoFile_ReturnsDefaults(t *testing.T) {
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Database defaults
	assert.Equal(t, dbconfig.DefaultType, cfg.Database.Type)
	assert.Equal(t, dbconfig.DefaultPostgresHost, cfg.Database.Postgres.Host)
	assert.Equal(t, dbconfig.DefaultPostgresPort, cfg.Database.Postgres.Port)
	assert.Equal(t, dbconfig.DefaultPostgresDatabase, cfg.Database.Postgres.Database)
	assert.Equal(t, dbconfig.DefaultPostgresSSLMode, cfg.Database.Postgres.SSLMode)

	// Task defaults
	assert.True(t, cfg.Regsync.Enabled)
	assert.True(t, cfg.Indexer.Enabled)
	assert.False(t, cfg.Name.Enabled)

	// Policy parameters come from the task's own defaults, not a second copy.
	assert.False(t, cfg.Policy.Enabled)
	assert.Equal(t, policy.DefaultInterval, cfg.Policy.Interval)
	assert.Equal(t, policy.DefaultRecordTimeout, cfg.Policy.RecordTimeout)
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	t.Setenv("RECONCILER_REGSYNC_ENABLED", "false")
	t.Setenv("RECONCILER_INDEXER_ENABLED", "false")
	t.Setenv("RECONCILER_NAME_ENABLED", "true")
	t.Setenv("RECONCILER_INDEXER_INTERVAL", "2h")
	t.Setenv("RECONCILER_NAME_INTERVAL", "30m")
	t.Setenv("RECONCILER_POLICY_ENABLED", "true")
	t.Setenv("RECONCILER_POLICY_INTERVAL", "5m")
	t.Setenv("RECONCILER_POLICY_RECORD_TIMEOUT", "45s")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.False(t, cfg.Regsync.Enabled)
	assert.False(t, cfg.Indexer.Enabled)
	assert.True(t, cfg.Name.Enabled)
	assert.Equal(t, 2*time.Hour, cfg.Indexer.Interval)
	assert.Equal(t, 30*time.Minute, cfg.Name.Interval)
	assert.True(t, cfg.Policy.Enabled)
	assert.Equal(t, 5*time.Minute, cfg.Policy.Interval)
	assert.Equal(t, 45*time.Second, cfg.Policy.RecordTimeout)
}
