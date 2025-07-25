package cli

import (
	"sync"

	spfviper "github.com/spf13/viper"
)

var (
	v  *spfviper.Viper
	mu = sync.Mutex{}
)

// Instance provides viper instance, or lazy-loads a new one if one has not been
// defined.
func viper() *spfviper.Viper {
	if v != nil {
		return v
	}

	mu.Lock()
	defer mu.Unlock()
	if v == nil {
		v = spfviper.New()
	}
	return v
}

// Reset creates a new Viper v. This should really only be used
// for testing purposes.
func reset() {
	mu.Lock()
	defer mu.Unlock()
	v = spfviper.New()
}
