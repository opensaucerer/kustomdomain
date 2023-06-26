package main

import (
	"log"
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
)

func main() {

	type Env struct {
		// Port for the server to listen on
		Port string `barfenv:"key=PORT;required=true"`
	}

	env := new(Env)

	barf.Env(env)

	allow := true
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  &allow, // enable request logging
		Recovery: &allow, // enable panic recovery so barf returns a 500 error instead of crashing
		Host:     "0.0.0.0",
	}); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	domains := map[string]string{
		"cendit.io":        "cenditio",
		"opensaucerer.com": "opensaucerer",
		"cendit.pro":       "cenditpro",
		"localhost":        "localhost",
		"164.92.156.239":   "digitaloceanipv4",
		"164.92.144.1":     "digitaloceangateway",
		"255.255.240.0":    "digitaloceanmask",
	}

	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {

		barf.Logger().Debug(r.Host)

		barf.Logger().Debug("referring host: " + r.Header.Get("ReferringHost"))

		if _, ok := domains[r.Header.Get("ReferringHost")]; !ok {

			barf.Response(w).Status(http.StatusNotFound).JSON(barf.Res{
				Status:  false,
				Data:    nil,
				Message: "Get thee behind me, oh Satan!",
			})
			return
		}

		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Welcome " + domains[r.Header.Get("ReferringHost")] + "!",
		})

	})

	barf.Get("/tls/ask", func(w http.ResponseWriter, r *http.Request) {

		q, _ := barf.Request(r).Query().JSON()

		log.Println(q)

		if _, ok := domains[q["domain"]]; ok {

			barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
				Status:  true,
				Data:    nil,
				Message: "Welcome son! Here's your TLS access",
			})
			return
		}

		barf.Response(w).Status(http.StatusNotFound).JSON(barf.Res{
			Status:  false,
			Data:    nil,
			Message: "Get thee behind me, oh TLS Satan!",
		})
	})

	// create & start server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
