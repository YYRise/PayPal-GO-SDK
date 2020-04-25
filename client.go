package paypalsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	kGetAccessTokenAPI = "/v1/oauth2/token"
)

// NewClient returns new Client struct
// APIBase is a base API URL, for testing you can use paypalsdk.APIBaseSandBox
func NewClient(clientID string, secret string, APIBase string) (*Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

// GetAccessToken returns struct of TokenResponse
// No need to call SetAccessToken to apply new access token for current Client
// Endpoint: POST /v1/oauth2/token
func (c *Client) GetAccessToken() (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	t := TokenResponse{}
	err = c.SendWithBasicAuth(req, &t)

	// Set Token fur current Client
	if t.Token != "" {
		c.Token = &t
		c.tokenExpiresAt = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	}

	return &t, err
}

// SetHTTPClient sets *http.Client to current client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.Client = client
}

// SetAccessToken sets saved token to current client
func (c *Client) SetAccessToken(token string) {
	c.Token = &TokenResponse{
		Token: token,
	}
	c.tokenExpiresAt = time.Time{}
}

// SetLog will set/change the output destination.
// If log file is set paypalsdk will log all requests and responses to this Writer
func (c *Client) SetLog(log io.Writer) {
	c.Log = log
}
func (c *Client) Send(req *http.Request, result interface{}) error {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")

	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/json")
	}

	var (
		err  error
		rsp  *http.Response
		data []byte
	)

	rsp, err = c.Client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	data, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if req.URL.Path != kGetAccessTokenAPI {
		var buf = &bytes.Buffer{}
		buf.WriteString("\n=========== Begin ============")
		buf.WriteString("\n【请求信息】")
		buf.WriteString(fmt.Sprintf("\n%s %d %s", req.Method, rsp.StatusCode, req.URL.String()))
		for key := range req.Header {
			buf.WriteString(fmt.Sprintf("\n%s: %s", key, req.Header.Get(key)))
		}
		buf.WriteString("\n【返回信息】")
		for key := range rsp.Header {
			buf.WriteString(fmt.Sprintf("\n%s: %s", key, rsp.Header.Get(key)))
		}
		buf.WriteString(fmt.Sprintf("\n%s", string(data)))
		buf.WriteString("\n===========  End  ============")

		logger.Println(buf.String())
	}

	switch rsp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if result != nil {
			if err = json.Unmarshal(data, result); err != nil {
				if err.Error() == "json: cannot unmarshal number into Go value of type string" {
					return nil
				}
				return err
			}
		}
		return nil
	case http.StatusUnauthorized:
		var e = &IdentityError{}
		e.Response = rsp
		if len(data) > 0 {
			if err = json.Unmarshal(data, e); err != nil {
				return err
			}
		}
		return e
	case http.StatusNoContent:
		//if req.Method == http.MethodDelete {
		//	return nil
		//}
		//fallthrough
	default:
		var e = &ResponseError{}
		e.Response = rsp
		if len(data) > 0 {
			if err = json.Unmarshal(data, e); err != nil {
				return err
			}
		}
		return e
	}

	return err
}

/*
// Send makes a request to the API, the response body will be
// unmarshaled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) error {
	var (
		err  error
		resp *http.Response
		data []byte
	)

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}
	resp, err = c.Client.Do(req)
	c.log(req, resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{Response: resp}
		data, err = ioutil.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			json.Unmarshal(data, errResp)
		}

		return errResp
	}

	if v == nil {
		return nil
	}

	if w, ok := v.(io.Writer); ok {
		io.Copy(w, resp.Body)
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
*/
// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired or already expired, it will try to get a new one before
// making the main request
// client.Token will be updated when changed
func (c *Client) SendWithAuth(req *http.Request, v interface{}) error {
	if c.Token != nil {
		if !c.tokenExpiresAt.IsZero() && c.tokenExpiresAt.Sub(time.Now()) < RequestNewTokenBeforeExpiresIn {
			// c.Token will be updated in GetAccessToken call
			if _, err := c.GetAccessToken(); err != nil {
				return err
			}
		}

		req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	}

	return c.Send(req, v)
}

// SendWithBasicAuth makes a request to the API using clientID:secret basic auth
func (c *Client) SendWithBasicAuth(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.ClientID, c.Secret)

	return c.Send(req, v)
}

// NewRequest constructs a request
// Convert payload to a JSON
func (c *Client) NewRequest(method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		var b []byte
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	request, err := http.NewRequest(method, url, buf)
	if err != nil {
		logrus.WithField("request", fmt.Sprintf("%+v", request)).WithError(err).Error("NewRequest:error")
	}
	logrus.WithField("request", fmt.Sprintf("%+v", request)).Info("NewRequest")
	return request, err
}

// log will dump request and response to the log file
func (c *Client) log(r *http.Request, resp *http.Response) {
	if c.Log != nil {
		var (
			reqDump  string
			respDump []byte
		)

		if r != nil {
			reqDump = fmt.Sprintf("%s %s. Data: %s", r.Method, r.URL.String(), r.Form.Encode())
		}
		if resp != nil {
			respDump, _ = httputil.DumpResponse(resp, true)
		}

		c.Log.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", reqDump, string(respDump))))
	}
}

func (c *Client) StdLog(req *http.Request, rsp *http.Response) {
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logger.Println("ioutil.ReadAll:error--", err)
		return
	}
	//if req.URL.Path != kGetAccessTokenAPI {
	var buf = &bytes.Buffer{}
	buf.WriteString("\n=========== Begin ============")
	buf.WriteString("\n【请求信息】")
	buf.WriteString(fmt.Sprintf("\n%s %d %s", req.Method, rsp.StatusCode, req.URL.String()))
	for key := range req.Header {
		buf.WriteString(fmt.Sprintf("\n%s: %s", key, req.Header.Get(key)))
	}
	buf.WriteString("\n【返回信息】")
	for key := range rsp.Header {
		buf.WriteString(fmt.Sprintf("\n%s: %s", key, rsp.Header.Get(key)))
	}
	buf.WriteString(fmt.Sprintf("\n%s", string(data)))
	buf.WriteString("\n===========  End  ============")

	logger.Println(buf.String())
	//}
}
