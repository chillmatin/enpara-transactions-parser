package swaggerui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(openAPISpecJSON))
	})

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/openapi.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
}

const openAPISpecJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Enpara Transactions Parser API",
    "version": "1.0.0",
    "description": "Upload an Enpara PDF statement and convert it to json, csv, xlsx, or ofx"
  },
  "servers": [
    {
      "url": "http://localhost:8080"
    }
  ],
  "paths": {
    "/api/v1/convert": {
      "post": {
        "summary": "Convert uploaded statement PDF",
        "requestBody": {
          "required": true,
          "content": {
            "multipart/form-data": {
              "schema": {
                "type": "object",
                "required": ["file"],
                "properties": {
                  "file": {
                    "type": "string",
                    "format": "binary"
                  },
                  "format": {
                    "type": "string",
                    "enum": ["json", "csv", "xlsx", "ofx"],
                    "default": "json"
                  },
                  "type": {
                    "type": "string",
                    "enum": ["auto", "type1", "type2"],
                    "default": "auto"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Converted file"
          },
          "400": {
            "description": "Invalid request"
          },
          "422": {
            "description": "Conversion failed"
          }
        }
      }
    },
    "/api/v1/health": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {
                      "type": "string",
                      "example": "ok"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/formats": {
      "get": {
        "summary": "List supported formats",
        "responses": {
          "200": {
            "description": "Supported formats",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "formats": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      },
                      "example": ["json", "csv", "xlsx", "ofx"]
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`
