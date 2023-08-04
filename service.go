package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NYTimes/gizmo/server"
)

type (
	SecurityCredentials struct {
		AccessKeyId     string `json:"AccessKeyId"`
		SecretAccessKey string `json:"SecretAccessKey"`
		Token           string `json:"Token"`
		Expiration      string `json:"Expiration"`
		Code            string `json:"Code"`
	}

	Network struct {
		Interfaces map[string][]string `json:"interfaces"`
	}

	MetadataValues struct {
		AmiId               string              `json:"ami-id"`
		AmiLaunchIndex      string              `json:"ami-launch-index"`
		AmiManifestPath     string              `json:"ami-manifest-path"`
		AvailabilityZone    string              `json:"availability-zone"`
		Hostname            string              `json:"hostname"`
		InstanceAction      string              `json:"instance-action"`
		InstanceId          string              `json:"instance-id"`
		InstanceType        string              `json:"instance-type"`
		LocalHostName       string              `json:"local-hostname"`
		LocalIpv4           string              `json:"local-ipv4"`
		Mac                 string              `json:"mac"`
		Profile             string              `json:"profile"`
		ReservationId       string              `json:"reservation-id"`
		User                string              `json:"User"`
		SecurityGroups      []string            `json:"security-groups"`
		SecurityCredentials SecurityCredentials `json:"security-credentials"`
		Network             Network             `json:"network"`
	}

	Config struct {
		Server           *server.Config
		MetadataValues   *MetadataValues
		MetadataPrefixes []string
		UserdataValues   map[string]string
		UserdataPrefixes []string
		NetworkPrefixes  []string
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

func movedPermanently(redirectPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
	}
}

func (s *MetadataService) GetAmiId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiId)
}

func (s *MetadataService) GetAmiLaunchIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiLaunchIndex)
}

func (s *MetadataService) GetAmiManifestPath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiManifestPath)
}

func (s *MetadataService) GetAvailabilityZone(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AvailabilityZone)
}

func (s *MetadataService) GetHostName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Hostname)
}

func (s *MetadataService) GetInstanceAction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceAction)
}

func (s *MetadataService) GetInstanceId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceId)
}

func (s *MetadataService) GetInstanceType(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceType)
}

func (s *MetadataService) GetLocalHostName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.LocalHostName)
}

func (s *MetadataService) GetLocalIpv4(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.LocalIpv4)
}

func (s *MetadataService) GetIAM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "security-credentials/")
}

func (s *MetadataService) GetMac(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Mac)
}

func (s *MetadataService) GetProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Profile)
}

func (s *MetadataService) GetReservationId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.ReservationId)
}

func (s *MetadataService) GetSecurityCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.User)
}

func (s *MetadataService) GetSecurityGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, strings.Join(s.config.MetadataValues.SecurityGroups, "\n"))
}

func (s *MetadataService) GetSecurityGroupIds(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, strings.Join(s.config.MetadataValues.Network.Interfaces["00:00:00:00:00:00"], "\n"))
}

func (s *MetadataService) GetSecurityCredentialDetails(w http.ResponseWriter, r *http.Request) {
	username := server.Vars(r)["username"]

	if username != s.config.MetadataValues.User {
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
		handlers[value+"/ami-id"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiId),
		}
		handlers[value+"/ami-launch-index"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiLaunchIndex),
		}
		handlers[value+"/ami-manifest-path"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiManifestPath),
		}
		handlers[value+"/placement/availability-zone"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAvailabilityZone),
		}
		handlers[value+"/hostname"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetHostName),
		}
		handlers[value+"/instance-action"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceAction),
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
		handlers[value+"/iam/security-credentials"] = map[string]http.HandlerFunc{
			"GET": movedPermanently(value + "/iam/security-credentials/"),
		}
		handlers[value+"/iam/security-credentials/"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityCredentials),
		}
		handlers[value+"/iam/security-credentials/{username}"] = map[string]http.HandlerFunc{
			"GET": service.GetSecurityCredentialDetails,
		}
		handlers[value+"/local-hostname"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetLocalHostName),
		}
		handlers[value+"/local-ipv4"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetLocalIpv4),
		}
		handlers[value+"/mac"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetMac),
		}
		handlers[value+"/profile"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetProfile),
		}
		handlers[value+"/reservation-id"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetReservationId),
		}
		handlers[value+"/security-groups"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityGroups),
		}
		handlers[value+"/network/interfaces/macs/00:00:00:00:00:00/security-group-ids"] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityGroupIds),
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
