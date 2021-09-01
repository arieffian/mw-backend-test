package response

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	defCfg      map[string]string
	initialized = false
)

// initialize this response configuration
func initialize() {
	viper.SetEnvPrefix("MW_TEST_RESPONSE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	defCfg = make(map[string]string)

	// general
	defCfg["general.500"] = "Internal server error"

	for k := range defCfg {
		err := viper.BindEnv(k)
		if err != nil {
			log.Errorf("Failed to bind env \"%s\" into respose configuration. Got %s", k, err)
		}
	}

	initialized = true
}

// SetConfig put response configuration key value
func SetConfig(key, value string) {
	viper.Set(key, value)
}

// Get fetch response configuration as string value
func Get(key string, httpCode int, typeResponse string) string {
	if !initialized {
		initialize()
	}

	// newKey := key + "." + strconv.Itoa(httpCode)
	// more safe if we use fmt.Sprinf to avoid parse variable
	newKey := fmt.Sprintf("%s.%d", key, httpCode)
	if len(typeResponse) > 0 {
		newKey += "." + typeResponse
	}

	ret := viper.GetString(newKey)
	if len(ret) == 0 {
		if ret, ok := defCfg[newKey]; ok {
			return ret
		}
		log.Debugf("%s config key not found", newKey)
	}

	ret = http.StatusText(httpCode)

	return ret
}

// Set response configuration key value
func Set(key, value string) {
	defCfg[key] = value
}
