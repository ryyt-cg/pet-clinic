package ds

import (
	"crypto/tls"
	"gin-petclinic-service/config/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
)

type HttpServer struct {
	r *gin.Engine
}

func NewHttpServer(r *gin.Engine) *HttpServer {
	return &HttpServer{
		r: r,
	}
}

// HttpRouter is a method on the HttpServer struct that sets up and returns a new http.Server instance.
// The http.Server instance uses the gin.Engine instance from the HttpServer struct to handle HTTP requests.
func (httpSrv *HttpServer) HttpRouter() *http.Server {
	// Create a TLS configuration
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	srv := &http.Server{
		Addr:    app.Config.Server.HttpPort,
		Handler: httpSrv.r,
		// Good practice: enforce timeouts for servers you create!
		//ReadTimeout:  10 * 60, // 10 minutes
		//WriteTimeout: 10 * 60, // 10 minutes
		//IdleTimeout:  10 * 60, // 10 minutes

		TLSConfig: tlsConfig,
	}

	err := http2.ConfigureServer(srv, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure HTTP/2 server")
	}

	return srv
}
