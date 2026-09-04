//go:build !windows

package servicemon

import "errors"

var errUnsupportedPlatform = errors.New("controle de serviços do Windows só é suportado em Windows")

// Start não é suportado fora do Windows. Ver servicemon_other.go para
// o motivo deste stub existir.
func Start(windowsServiceName string) error {
	return errUnsupportedPlatform
}

// Stop não é suportado fora do Windows.
func Stop(windowsServiceName string) error {
	return errUnsupportedPlatform
}

// Restart não é suportado fora do Windows.
func Restart(windowsServiceName string) error {
	return errUnsupportedPlatform
}
