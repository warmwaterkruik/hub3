package experimental

import (
	"net/http"

	"github.com/delving/rapid-saas/config"
	"github.com/go-chi/render"
)

// listNameSpaces list all currently defined NameSpace object
func listNameSpaces(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, config.Config.NameSpaceMap.ByPrefix())
	return
}
