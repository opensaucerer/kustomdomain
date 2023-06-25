package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/opensaucerer/barf"
)

func main() {

	type Env struct {
		// Port for the server to listen on
		Port string `barfenv:"key=PORT;required=true"`
	}

	env := new(Env)

	allow := true
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  &allow, // enable request logging
		Recovery: &allow, // enable panic recovery so barf returns a 500 error instead of crashing
	}); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// barf tries to be as unobtrusive as possible, so your route handlers still
	// inherit the standard http.ResponseWriter and *http.Request parameters
	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {

		barf.Logger().Debug(r.Host)

		if strings.Split(r.Host, ".")[0] == "opensaucerer" {

			barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
				Status:  true,
				Data:    nil,
				Message: "Welcome " + strings.Split(r.Host, ".")[0],
			})
			return
		}

		barf.Response(w).Status(http.StatusNotFound).JSON(barf.Res{
			Status:  false,
			Data:    nil,
			Message: "Get thee behind me, oh Satan!",
		})
	})

	// create & start server
	if err := barf.Beck(); err != nil {
		// barf exposes a logger instance
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
