package service

import (
	"net/http"
	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
)

type (

	Config struct {
		Server    *server.Config
		MetadataItems map[string]interface{}
	}
	MetadataService struct {
		config *Config
	}
)

func NewMetadataService(cfg *Config) *MetadataService {
        return &MetadataService{cfg}
}

func (s *MetadataService) Middleware(h http.Handler) http.Handler {
	return h
}

func (s *MetadataService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {

		status, res, err := j(r)
		if err != nil {
			server.LogWithFields(r).WithFields(logrus.Fields{
				"error": err,
			}).Error("problems with serving request")
			return http.StatusServiceUnavailable, nil, &jsonErr{"sorry, this service is unavailable"}
		}

		server.LogWithFields(r).Info("success!")
		return status, res, nil
	}
}

func (s *MetadataService) GetMetadataItem(r *http.Request) (int, interface{}, error) {
	res := s.config.MetadataItems[r.URL.Path]
	return http.StatusOK, res, nil
}

func (s *MetadataService) GetIndex(r *http.Request) (int, interface{}, error) {
	return http.StatusOK, "Mock EC2 Metadata Service", nil
}

// JSONEndpoints is a listing of all endpoints available in the MetadataService.
func (s *MetadataService) Endpoints() map[string]map[string]http.HandlerFunc {

	handlers := make(map[string]map[string]http.HandlerFunc)
	for url, value := range s.config.MetadataItems {
        server.Log.Info("adding route for url", url, " value ", value)

        handlers[url] = map[string]http.HandlerFunc {
			"GET": server.JSONToHTTP(s.GetMetadataItem).ServeHTTP,
		}
	}
	handlers["/"] = map[string]http.HandlerFunc {
			"GET": server.JSONToHTTP(s.GetIndex).ServeHTTP,
	}
	return handlers
}

func (s *MetadataService) Prefix() string {
	return "/"
}

type jsonErr struct {
	Err string `json:"error"`
}

func (e *jsonErr) Error() string {
	return e.Err
}