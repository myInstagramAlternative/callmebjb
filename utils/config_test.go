package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	//	"gopkg.in/yaml.v3"

	// Import the package containing the function to be tested
	"callmebjb/utils"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary YAML configuration file for testing
	testConfig := `
modem:
  port: "/dev/tty.usbserial-0001"
  baudrate: 115200
server:
  listen: "localhost"
  port: "8080"
`
	tmpfile, err := os.CreateTemp("", "config_test_*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(testConfig)
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Call the function to be tested
	var cfg utils.Config
	utils.ReadConfig(&cfg, tmpfile.Name())
	// Define the expected configuration struct
	expectedCfg := utils.Config{
		Modem: struct {
			Port     string `yaml:"port"`
			BaudRate int    `yaml:"baudrate"`
		}{
			Port:     "/dev/tty.usbserial-0001",
			BaudRate: 115200,
		},
		Server: struct {
			Listen string `yaml:"listen"`
			Port   string `yaml:"port"`
		}{
			Listen: "localhost",
			Port:   "8080",
		},
	}

	// Assert that the actual configuration matches the expected configuration
	assert.Equal(t, expectedCfg, cfg)
}
