package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gizmo/web"
)

type (
	SecurityCredentials struct {
		User            string `json:"User"`
		AccessKeyId     string `json:"AccessKeyId"`
		SecretAccessKey string `json:"SecretAccessKey"`
		Token           string `json:"Token"`
		Expiration      string `json:"Expiration"`
	}

	MetadataValues struct {
		Hostname            string              `json:"hostname"`
		InstanceId          string              `json:"instance-id"`
		InstanceType        string              `json:"instance-type"`
		SecurityCredentials SecurityCredentials `json:"security-credentials"`
	}

	Config struct {
		Server           *server.Config
		MetadataValues   *MetadataValues
		MetadataPrefixes []string
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

func (s *MetadataService) GetHostName(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res := fmt.Sprint(s.config.MetadataValues.Hostname)
	fmt.Fprintf(w, res)
	server.Log.Info("GetHostName returning: ", res)
	return
}

func (s *MetadataService) GetInstanceId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res := fmt.Sprint(s.config.MetadataValues.InstanceId)
	fmt.Fprintf(w, res)
	server.Log.Info("GetInstanceId returning: ", res)
	return
}

func (s *MetadataService) GetInstanceType(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res := fmt.Sprint(s.config.MetadataValues.InstanceType)
	fmt.Fprintf(w, res)
	server.Log.Info("GetInstanceType returning: ", res)
	return
}

func (s *MetadataService) GetIAM(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res := fmt.Sprint("security-credentials/")
	fmt.Fprintf(w, res)
	server.Log.Info("GetIAM returning: ", res)
	return
}

func (s *MetadataService) GetSecurityCredentials(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res := fmt.Sprint(s.config.MetadataValues.SecurityCredentials.User)
	server.Log.Info("GetSecurityCredentials returning: ", res)
	fmt.Fprintf(w, res)
	return
}

func (s *MetadataService) GetSecurityCredentialDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	username := web.Vars(r)["username"]

	if username == s.config.MetadataValues.SecurityCredentials.User {
		details, err := json.MarshalIndent(s.config.MetadataValues.SecurityCredentials, "", "\t")
		if err != nil {
			server.Log.Error("error converting security credentails to json: ", err)
			http.Error(w, "", http.StatusNotFound)

			return
		} else {
			server.Log.Info("GetSecurityCredentialDetails returning: ", details)

			w.Write(details)
			return
		}
	} else {
		server.Log.Error("error, IAM user not found")
		http.Error(w, "", http.StatusNotFound)
	}

	return
}

func (s *MetadataService) GetMetadataIndex(w http.ResponseWriter, r *http.Request) {

	index := []string{"hostname",
		"instance-id",
		"instance-type",
		"iam/"}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res := fmt.Sprint(strings.Join(index, "\n"))
	server.Log.Info("GetMetadataIndex returning: ", res)
	fmt.Fprintf(w, res)
	return
}

func (s *MetadataService) GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mock EC2 Metadata Service")
	return
}

// Endpoints is a listing of all endpoints available in the MetadataService.
func (service *MetadataService) Endpoints() map[string]map[string]http.HandlerFunc {

	handlers := make(map[string]map[string]http.HandlerFunc)
	for index, value := range service.config.MetadataPrefixes {
		server.Log.Info("adding Metadata prefix (", index, ") ", value)

		handlers[value+"/"] = map[string]http.HandlerFunc{
			"GET": service.GetMetadataIndex,
		}
		handlers[value+"/hostname"] = map[string]http.HandlerFunc{
			"GET": service.GetHostName,
		}
		handlers[value+"/instance-id"] = map[string]http.HandlerFunc{
			"GET": service.GetInstanceId,
		}
		handlers[value+"/instance-type"] = map[string]http.HandlerFunc{
			"GET": service.GetInstanceType,
		}
		handlers[value+"/iam/"] = map[string]http.HandlerFunc{
			"GET": service.GetIAM,
		}
		handlers[value+"/iam/security-credentials/"] = map[string]http.HandlerFunc{
			"GET": service.GetSecurityCredentials,
		}
		handlers[value+"/iam/security-credentials/{username}"] = map[string]http.HandlerFunc{
			"GET": service.GetSecurityCredentialDetails,
		}
	}
	handlers["/"] = map[string]http.HandlerFunc{
		"GET": service.GetIndex,
	}
	return handlers
}

func (s *MetadataService) Prefix() string {
	return "/"
}

type error struct {
	Err string
}

func (e *error) Error() string {
	return e.Err
}
