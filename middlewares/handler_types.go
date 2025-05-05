package middlewares

import (
	"net/http"

	"github.com/harshvardha/artOfSoftwareEngineering/controllers"
)

// authenticated endpoints handler
type authenticatedRequestHandler func(http.ResponseWriter, *http.Request, *controllers.IDAndRole, *string)
