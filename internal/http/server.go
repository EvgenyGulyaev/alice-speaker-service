package http

import (
	"aliceSpeakerService/internal/http/middleware"
	"aliceSpeakerService/internal/http/routes"
	"net/http"
	"strings"

	"github.com/go-www/silverlining"
)

type Server struct {
	port string
}

func GetServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) StartHandle() error {
	return silverlining.ListenAndServe(s.port, func(ctx *silverlining.Context) {
		HandleRequest(ctx)
	})
}

func HandleRequest(ctx *silverlining.Context) {
	updateHeaders(ctx)
	path := string(ctx.Path())
	switch ctx.Method() {
	case silverlining.MethodGET:
		handleGet(ctx, path)
	case silverlining.MethodPOST:
		handlePost(ctx, path)
	case silverlining.MethodPATCH:
		handlePatch(ctx, path)
	case silverlining.MethodOPTIONS:
		ctx.WriteHeader(http.StatusNoContent)
	default:
		routes.NotFound(ctx)
	}
}

func handleGet(ctx *silverlining.Context, path string) {
	switch path {
	case "/health":
		routes.GetHealth(ctx)
	case "/api/accounts":
		middleware.RequireServiceToken(func(c *silverlining.Context) {
			routes.GetAccounts(c)
		})(ctx)
	default:
		if parts := apiPathParts(path); len(parts) == 4 && parts[0] == "api" && parts[1] == "accounts" && parts[3] == "resources" {
			middleware.RequireServiceToken(func(c *silverlining.Context) {
				routes.GetAccountResources(c, parts[2])
			})(ctx)
			return
		}
		routes.NotFound(ctx)
	}
}

func handlePost(ctx *silverlining.Context, path string) {
	body, err := ctx.Body()
	if err != nil {
		routes.GetError(ctx, &routes.Error{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}

	switch path {
	case "/api/accounts":
		middleware.RequireServiceToken(func(c *silverlining.Context) {
			routes.PostAccount(c, body)
		})(ctx)
	case "/api/announce/scenario":
		middleware.RequireServiceToken(func(c *silverlining.Context) {
			routes.PostAnnounceScenario(c, body)
		})(ctx)
	default:
		if parts := apiPathParts(path); len(parts) == 4 && parts[0] == "api" && parts[1] == "accounts" && parts[3] == "refresh" {
			middleware.RequireServiceToken(func(c *silverlining.Context) {
				routes.PostAccountRefresh(c, parts[2])
			})(ctx)
			return
		}
		if parts := apiPathParts(path); len(parts) == 4 && parts[0] == "api" && parts[1] == "accounts" && parts[3] == "import-cookies" {
			middleware.RequireServiceToken(func(c *silverlining.Context) {
				routes.PostAccountImportCookies(c, parts[2], body)
			})(ctx)
			return
		}
		if parts := apiPathParts(path); len(parts) == 4 && parts[0] == "api" && parts[1] == "accounts" && parts[3] == "cleanup-scenarios" {
			middleware.RequireServiceToken(func(c *silverlining.Context) {
				routes.PostAccountCleanupScenarios(c, parts[2], body)
			})(ctx)
			return
		}
		routes.NotFound(ctx)
	}
}

func handlePatch(ctx *silverlining.Context, path string) {
	body, err := ctx.Body()
	if err != nil {
		routes.GetError(ctx, &routes.Error{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}
	if parts := apiPathParts(path); len(parts) == 3 && parts[0] == "api" && parts[1] == "accounts" {
		middleware.RequireServiceToken(func(c *silverlining.Context) {
			routes.PatchAccount(c, parts[2], body)
		})(ctx)
		return
	}
	routes.NotFound(ctx)
}

func apiPathParts(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return nil
	}
	return parts
}

func updateHeaders(ctx *silverlining.Context) {
	ctx.ResponseHeaders().Set("Access-Control-Allow-Origin", "*")
	ctx.ResponseHeaders().Set("Access-Control-Allow-Credentials", "true")
	ctx.ResponseHeaders().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.ResponseHeaders().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
}
