package security

import (
	"fmt"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	fmt.Println(Init())
}

func TestGet(t *testing.T) {
	fmt.Println(Get())
}

func TestGetEnv(t *testing.T) {
	fmt.Println(os.Getenv(EnvMasterKey))
}

func TestSetEnv(t *testing.T) {
	fmt.Println(os.Setenv(EnvMasterKey, "OYITmq0b/MqsVnZ5epEOkQoxz0lTxt/h0rtpei+ad4g="))
}
