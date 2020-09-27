package secret_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	rpccodes "google.golang.org/grpc/codes"
	"io/ioutil"
	"net/http"
	"path"
)

var (
	requestTypeLogin      = "Login"
	requestTypeEncrypt    = "Encrypt"
	requestTypeDisconnect = "Disconnect"
	requestTypeLogout     = "Logout"
	requestTypeStatus     = "Status"
	apiRoute              = "portalintf"
)

// Request - базовая структура запроса с API secret
type AuthRequest struct {
	Request
	Proxy    string `json:"UE-Proxy,omitempty"`
	Username string `json:"UE-Username,omitempty"`
	Password string `json:"UE-Password,omitempty"`
	Data     string `json:"Data,omitempty"`

	encrypted bool `json:"-"`
}

// Request - базовая структура запроса с API secret_api
type Request struct {
	Vendor     string `json:"Vendor"`
	Password   string `json:"RequestPassword"`
	APIVersion string `json:"APIVersion"`
	Category   string `json:"RequestCategory"`
	Type       string `json:"RequestType"`
	MAC        string `json:"UE-MAC,omitempty"`
	UserIP     string `json:"UE-IP,omitempty"`

	method   string
	destAddr string `json:"-"`
}

// newRequest - конструктор базоваого запроса
func newRequest() *Request {
	return &Request{
		Vendor:     "Secret",
		APIVersion: "1.0",
		Category:   "UserOnlineControl",
	}
}

// loginRequest - конструктор login запроса
func ip2macRequest(remoteAddr, password, userIP string) (req *Request) {
	req = newRequest()
	req.Password = password
	req.destAddr = remoteAddr
	req.Type = requestTypeStatus

	req.UserIP = userIP
	req.method = "ip_to_mac"
	return
}

func mac2ipRequest(remoteAddr, password, mac string) (req *Request) {
	req = newRequest()
	req.Password = password
	req.destAddr = remoteAddr
	req.Type = requestTypeStatus

	req.MAC = mac
	req.method = "mac_to_ip"
	return
}

// send - отправляет запрос на API secret_api
func (req *Request) send(paths ...string) (body []byte, err error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		err = fmt.Errorf("unable to marshal request: %s", err)
		return
	}

	client := &http.Client{
		Timeout: options.Timeout,
	}
	paths = append([]string{req.destAddr}, paths...)
	sURL := path.Join(paths...)
	request, err := http.NewRequest("POST", sURL, bytes.NewBuffer(requestBody))
	if err != nil {
		err = fmt.Errorf("unable to create new HTTP POST request (%s): %s", sURL, err)
		return
	}
	request.Header.Set("Content-type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("unable to make HTTP POST request  (%s): %s", sURL, err)
		return
	}
	if response.Body != nil {
		defer response.Body.Close()
		if body, err = ioutil.ReadAll(response.Body); err != nil {
			err = fmt.Errorf("unable to read body HTTP POST response (%s): %s", sURL, err)
			return
		}
	}

	return
}

func errorByCode(code int, msg string) (err grpc.Error) {
	switch code {
	case http.StatusOK, http.StatusSwitchingProtocols:
	case http.StatusMultipleChoices:
		err = grpc.NewError(rpccodes.NotFound, fmt.Errorf("lookup failed for given MAC or IP address: %s", msg))
	case http.StatusMovedPermanently:
		err = grpc.NewError(rpccodes.FailedPrecondition, fmt.Errorf("login failed: %s", msg))
	case http.StatusBadRequest:
		err = grpc.NewError(rpccodes.Internal, fmt.Errorf("internal remote server error: %s", msg))
	case http.StatusUnauthorized:
		err = grpc.NewError(rpccodes.Unimplemented, fmt.Errorf("remote server authentication connection error occured or the connection request times out: %s", msg))
	case http.StatusFound:
		err = grpc.NewError(rpccodes.InvalidArgument, fmt.Errorf("JSON request is not well-formed: %s", msg))
	case http.StatusSeeOther:
		err = grpc.NewError(rpccodes.Unimplemented, fmt.Errorf("version mismatch: %s", msg))
	case http.StatusNotModified:
		err = grpc.NewError(rpccodes.InvalidArgument, fmt.Errorf("request type is not supported: %s", msg))
	case http.StatusUseProxy:
		err = grpc.NewError(rpccodes.InvalidArgument, fmt.Errorf("request category not supported: %s", msg))
	case 306:
		err = grpc.NewError(rpccodes.Unauthenticated, fmt.Errorf("the request password is mismatched: %s", msg))
	default:
		err = grpc.NewError(rpccodes.Unknown, fmt.Errorf("Secret request returning with uknown code: %d: %s", code, msg))
	}
	return
}
