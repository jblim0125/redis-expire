package injectors

import (
	"github.com/jblim0125/redredis-expire/controllers"
)

// Version version injector
type Version struct{}

// Init version controller create
func (Version) Init(in *Injector) *controllers.Version {
	return controllers.Version{}.New()
}
