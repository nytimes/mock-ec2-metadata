package service

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		UserdataValues   map[string]string
		UserdataPrefixes []string
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

// middleware for adding plaintext content type
func plainText(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		h(w, r)
	}
}

func (s *MetadataService) GetHostName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Hostname)
}

func (s *MetadataService) GetInstanceId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceId)
}

func (s *MetadataService) GetInstanceType(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceType)
}

func (s *MetadataService) GetIAM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "security-credentials/")
}

func (s *MetadataService) GetSecurityCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.SecurityCredentials.User)
}

func (s *MetadataService) GetSecurityCredentialDetails(w http.ResponseWriter, r *http.Request) {
	username := web.Vars(r)["username"]

	if username != s.config.MetadataValues.SecurityCredentials.User {
		server.Log.Error("error, IAM user not found")
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(s.config.MetadataValues.SecurityCredentials)
	if err != nil {
		server.Log.Error("error converting security credentails to json: ", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	server.LogWithFields(r).Info("GetSecurityCredentialDetails returning: %#v",
		s.config.MetadataValues.SecurityCredentials)
}

func (s *MetadataService) GetMetadataIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `hostname
instance-id
instance-type
iam`)
}

func (s *MetadataService) GetUserData(w http.ResponseWriter, r *http.Request) {

	for index, value := range s.config.UserdataValues {
		fmt.Fprintf(w, fmt.Sprint(index+"="+value+"\n"))
	}
}

func (s *MetadataService) GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mock EC2 Metadata Service")
}

// Endpoints is a listing of all endpoints available in the MetadataService.
func (service *MetadataService) Endpoints() map[string]map[string]http.HandlerFunc {
	handlers := map[string]map[string]http.HandlerFunc{}

	for index, value := range service.config.MetadataPrefixes {
		server.Log.Info("adding Metadata prefix (", index, ") ", value)
		handlers[value+"/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetMetadataIndex),
		}
		handlers[value+"/hostname"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetHostName),
		}
		handlers[value+"/instance-id"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceId),
		}
		handlers[value+"/instance-type"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceType),
		}
		handlers[value+"/iam/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetIAM),
		}
		handlers[value+"/iam/security-credentials/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityCredentials),
		}
		handlers[value+"/iam/security-credentials/{username}"] = map[string]http.HandlerFunc{
			"GET": service.GetSecurityCredentialDetails,
		}
	}

	for index, value := range service.config.UserdataPrefixes {
		server.Log.Info("adding Userdata prefix (", index, ") ", value)

		handlers[value+"/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetUserData),
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
