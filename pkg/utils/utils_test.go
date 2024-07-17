package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	testee "github.com/robinmin/askllm/pkg/utils"
)

func TestNewInstance(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Field1 string `default:"value1"`
		Field2 int    `default:"42"`
	}

	t.Run("HappyPath", func(t *testing.T) {
		obj := testee.NewInstance[TestStruct]()
		assert.NotNil(t, obj)
		assert.Equal(t, "value1", obj.Field1)
		assert.Equal(t, 42, obj.Field2)
	})

	t.Run("ErrorOnSetDefaults", func(t *testing.T) {
		obj := testee.NewInstance[int]() // int does not have default values
		assert.Nil(t, obj)
	})
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	type TestConfig struct {
		Key string
	}

	t.Run("HappyPath", func(t *testing.T) {
		// prepare data file
		data := []byte("key: value\n")

		// Create a temporary file
		yamlFile, err := testee.WriteTempFile("HappyPath", data)
		assert.NoError(t, err)

		// Defer cleanup to ensure it happens even if the function returns early
		defer func() {
			err := testee.CleanupTempFile(yamlFile)
			assert.NoError(t, err)
		}()

		expectedConfig := &TestConfig{Key: "value"}

		// Load the config
		config, err := testee.LoadConfig[TestConfig](yamlFile)
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, config)

		err = os.Remove(yamlFile)
		assert.NoError(t, err)
	})

	t.Run("ErrorOnReadFile", func(t *testing.T) {
		config, err := testee.LoadConfig[TestConfig]("nonexistent.yaml")
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("ErrorOnUnmarshal", func(t *testing.T) {
		// prepare data file
		data := []byte("invalid yaml")

		// Create a temporary file
		yamlFile, err := testee.WriteTempFile("ErrorOnUnmarshal", data)
		assert.NoError(t, err)

		// Defer cleanup to ensure it happens even if the function returns early
		defer func() {
			err := testee.CleanupTempFile(yamlFile)
			assert.NoError(t, err)
		}()

		config, err := testee.LoadConfig[TestConfig](yamlFile)
		assert.Error(t, err)
		assert.Nil(t, config)

		err = os.Remove(yamlFile)
		assert.NoError(t, err)
	})
}

func TestSaveConfig(t *testing.T) {
	t.Parallel()

	type TestConfig struct {
		Key string
	}

	t.Run("HappyPath", func(t *testing.T) {
		yamlFile := "/tmp/config_cf689862_25b5_420f_bd39_745e855c00be.yaml"
		config := &TestConfig{Key: "value"}

		err := testee.SaveConfig(config, yamlFile)
		assert.NoError(t, err)

		configLoaded, err := testee.LoadConfig[TestConfig](yamlFile)
		assert.NoError(t, err)
		assert.Equal(t, configLoaded, config)

		err = os.Remove(yamlFile)
		assert.NoError(t, err)
	})

	t.Run("ErrorOnWriteFile", func(t *testing.T) {
		config := &TestConfig{Key: "value"}
		err := testee.SaveConfig(config, "/root/config.yaml") // permission denied
		assert.Error(t, err)
	})
}
