package response

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initialize()
		assert.Equal(t, true, initialized)
		initialized = false
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		value := Get("auth", 401, "unauthorized")
		expectedValue := "Password/Email kamu salah. Periksa dan coba lagi ya"
		assert.Equal(t, expectedValue, value)
		initialized = false
	})
	t.Run("key-not-found", func(t *testing.T) {
		value := Get("unknown", 401, "")
		expectedValue := "Unauthorized"
		assert.Equal(t, expectedValue, value)
		initialized = false
	})
}

func TestSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initialize()
		Set("test.123", "test")
		value := Get("test", 123, "")
		expectedValue := "test"
		assert.Equal(t, expectedValue, value)
		delete(defCfg, "test.123")
		initialized = false
	})
}

func TestSetConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initialize()
		SetConfig("test.123", "test-test")
		value := viper.GetString("test.123")
		expectedValue := "test-test"
		assert.Equal(t, expectedValue, value)
		delete(defCfg, "test.123")
		initialized = false
	})
}
