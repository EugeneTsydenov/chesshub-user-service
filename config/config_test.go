package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	databasePasswordEnvKey = "DATABASE_PASSWORD"
	configPathEnvKey       = "CONFIG_PATH"
)

func TestLoad(t *testing.T) {
	t.Run("should successfully loads config yaml and env in local, develop, prod environments", func(t *testing.T) {
		environments := []string{LocalEnv, DevelopEnv, ProdEnv}
		for _, env := range environments {
			env := env

			t.Run(fmt.Sprintf("should successfully loads config in %s environment", env), func(t *testing.T) {
				t.Cleanup(clean)
				setEnv(t)

				c, err := Load(env, os.Getenv(configPathEnvKey))

				assert.NoError(t, err)
				assert.NotNil(t, c)
				assert.Equal(t, env, c.App.Env)
				assert.NotEqual(t, "", c.Database.Password)
				assert.Equal(t, "postgres", c.Database.User)
			})
		}
	})

	t.Run("should load with error when env is invalid", func(t *testing.T) {
		t.Cleanup(clean)
		setEnv(t)

		c, err := Load("nonexistent", os.Getenv(configPathEnvKey))

		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Contains(t, err.Error(), "invalid environment")
	})

	t.Run("should return error when config path is empty", func(t *testing.T) {
		t.Cleanup(clean)

		c, err := Load("", "")

		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("should return error when config path is invalid", func(t *testing.T) {
		t.Cleanup(clean)

		invalidPath := "/non/existing/path"
		t.Setenv(configPathEnvKey, invalidPath)

		c, err := Load(LocalEnv, os.Getenv(configPathEnvKey))

		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Contains(t, err.Error(), "failed to read base config")
	})

	t.Run("should work with absolute config path", func(t *testing.T) {
		t.Cleanup(clean)

		wd, err := os.Getwd()
		require.NoError(t, err)

		absPath := filepath.Join(wd, "../config")
		t.Setenv(configPathEnvKey, absPath)
		t.Setenv(databasePasswordEnvKey, "postgres")

		c, err := Load(LocalEnv, absPath)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, LocalEnv, c.App.Env)
	})

	t.Run("should generate correct DSN string", func(t *testing.T) {
		d := DatabaseConfig{
			Host:     "testhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
			SSLMode:  "disable",
		}

		expectedDSN := "host=testhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
		assert.Equal(t, expectedDSN, d.DSN())
	})
}

func setEnv(t *testing.T) {
	t.Setenv(databasePasswordEnvKey, "postgres")
	t.Setenv(configPathEnvKey, "../config")
}

func clean() {
	_ = os.Unsetenv("DATABASE_PASSWORD")
	_ = os.Unsetenv("CONFIG_PATH")

	viper.Reset()
}

func TestLoadDefaultConfig(t *testing.T) {
	t.Run("should successfully loads default config", func(t *testing.T) {
		t.Cleanup(clean)

		err := loadDefaultCfg("../config")
		assert.NoError(t, err)
	})

	t.Run("should load with error when config path is invalid", func(t *testing.T) {
		t.Cleanup(clean)

		err := loadDefaultCfg("/non/existing/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read base config")
	})
}

func TestLoadConfigByEnv(t *testing.T) {
	t.Run("should successfully loads config by env", func(t *testing.T) {
		t.Cleanup(clean)

		environments := []string{LocalEnv, DevelopEnv, ProdEnv}

		for _, env := range environments {
			t.Run(fmt.Sprintf("load config for %s environment", env), func(t *testing.T) {
				err := loadEnvCfg(env, "../config")
				assert.NoError(t, err)
			})
		}
	})

	t.Run("should load with error when env is invalid", func(t *testing.T) {
		t.Cleanup(clean)

		err := loadEnvCfg("nonexistent", "../config")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid environment")
	})
}

func TestLoadEnvVariables(t *testing.T) {
	t.Cleanup(clean)

	err := loadEnvVariables()
	assert.NoError(t, err)
}
