package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/chillmatin/enpara-transactions-parser/internal/handlers"
	"github.com/chillmatin/enpara-transactions-parser/internal/swaggerui"
	"github.com/gin-gonic/gin"
)

func main() {
	defaultHost := getEnvOrDefault("ENPARA_API_HOST", "localhost")
	defaultPort := getEnvIntOrDefault("ENPARA_API_PORT", 8080)
	defaultSwaggerEnabled := getEnvBoolOrDefault("ENPARA_API_SWAGGER", false)

	host := flag.String("host", defaultHost, "API bind host")
	port := flag.Int("port", defaultPort, "API bind port")
	enableSwagger := flag.Bool("swagger", defaultSwaggerEnabled, "Enable Swagger UI at /swagger")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: enpara-api [flags]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Start the HTTP API for converting Enpara PDF statements.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Endpoints: /api/v1/convert, /api/v1/formats, /api/v1/health\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nEnvironment Variables:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  ENPARA_API_HOST     Default host (current default: %s)\n", defaultHost)
		fmt.Fprintf(flag.CommandLine.Output(), "  ENPARA_API_PORT     Default port (current default: %d)\n", defaultPort)
		fmt.Fprintf(flag.CommandLine.Output(), "  ENPARA_API_SWAGGER  Enable Swagger by default (current default: %t)\n", defaultSwaggerEnabled)
		fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  enpara-api\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  enpara-api --host 0.0.0.0 --port 8080 --swagger\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  ENPARA_API_SWAGGER=true ENPARA_API_PORT=9090 enpara-api\n")
	}
	flag.Parse()

	if *port < 1 || *port > 65535 {
		log.Fatalf("invalid port %d", *port)
	}

	address := fmt.Sprintf(":%d", *port)
	bindHost := strings.TrimSpace(*host)
	if bindHost != "" {
		address = fmt.Sprintf("%s:%d", bindHost, *port)
	}

	displayHost := bindHost
	if displayHost == "" || displayHost == "0.0.0.0" || displayHost == "::" {
		displayHost = "localhost"
	}

	router := gin.Default()

	router.POST("/api/v1/convert", handlers.HandleConvert)
	router.GET("/api/v1/formats", handlers.HandleFormats)
	router.GET("/api/v1/health", handlers.HandleHealth)
	if *enableSwagger {
		swaggerui.RegisterRoutes(router)
		log.Printf("swagger ui enabled at /swagger")
	}

	log.Printf("starting api server on http://%s:%d", displayHost, *port)
	if err := router.Run(address); err != nil {
		log.Fatalf("start api server: %v", err)
	}
}

func getEnvOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvIntOrDefault(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("invalid %s=%q, using default %d", name, value, fallback)
		return fallback
	}

	return parsed
}

func getEnvBoolOrDefault(name string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("invalid %s=%q, using default %t", name, value, fallback)
		return fallback
	}

	return parsed
}
