package bootstrap

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Route struct {
	TargetHost   string
	MatchPath    string
	ServiceState string
}

func (r *Route) GetRouteMatchPath() string {
	if len(r.MatchPath) == 0 {
		return ""
	}
	return r.MatchPath
}

func (r *Route) GetRouteTargetHost() string {
	if len(r.TargetHost) == 0 {
		return ""
	}
	return r.TargetHost
}

func BuildRoute(routes []Route, r gin.Engine) {
	routeMap := make(map[string]string)
	if routes != nil {
		for _, route := range routes {
			routePath := route.GetRouteMatchPath()
			routeTarget := route.GetRouteTargetHost()

			if len(routePath) == 0 || len(routeTarget) == 0 {
				continue
			}

			relativePath := routePath + "/*name"

			target := routeMap[routePath]
			if len(target) > 0 {
				continue
			}

			routeMap[routePath] = routeTarget
			r.Any(relativePath, func(c *gin.Context) {
				proxyUrl, _ := url.Parse(routeTarget)

				proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
				proxy.Transport = &http.Transport{
					// 响应超时时间
					ResponseHeaderTimeout: time.Duration(120) * time.Second,
				}

				proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
					// 正常响应502
					writer.WriteHeader(http.StatusBadGateway)
				}

				proxy.ServeHTTP(c.Writer, c.Request)
			})
		}
	}
}
