package middleware

import "kkagitala/go-rest-api/service"

// Middleware describes a service middleware.
type Middleware func(service service.Service) service.Service
