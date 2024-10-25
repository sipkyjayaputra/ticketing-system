package helpers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HrsvClient struct {
	IsAuthenticated bool
	URL             string
	Session         *http.Client
	Email       	string
	Password        string
	token           string
}

func NewHrsvClient(hostname string, email string, password string) *HrsvClient {
	url := fmt.Sprintf("https://%s/api", hostname)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable certificate verification
		},
	}

	session := &http.Client{
		Transport: tr,
	}
	
	return &HrsvClient{
		IsAuthenticated: false,
		URL:             url,
		Session:         session,
		Email:        	 email,
		Password:        password,
		token:			 "",
	}
}

func (c *HrsvClient) Authenticate() error {
	// Prepare the payload for authentication
	payload := map[string]string{
		"email":    c.Email,
		"password": c.Password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.URL+"/auth/login", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.Session.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: status code %d", resp.StatusCode)
	}

	// Read and unmarshal the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for access token in the response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("data field not found or not a valid map")
	}
	
	token, ok := data["access_token"].(string)
	if !ok {
		return fmt.Errorf("access token not found in response")
	}
	
	// Set the authentication status and token
	c.IsAuthenticated = true
	c.token = "Bearer " + token
	return nil
}

func (c *HrsvClient) AuthenticateWithPayload(payload map[string]string) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.URL+"/auth/login", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.Session.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: status code %d", resp.StatusCode)
	}

	// // Read and unmarshal the response body
	// respBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return fmt.Errorf("failed to read response body: %w", err)
	// }

	// var result map[string]interface{}
	// if err := json.Unmarshal(respBody, &result); err != nil {
	// 	return fmt.Errorf("failed to unmarshal response: %w", err)
	// }

	// // Check for success and access token in the response
	// statusCode, ok := result["status_code"].(int)
	// if !ok || !statusCode ! {
	// 	return fmt.Errorf("authentication failed: unexpected response format or unsuccessful login")
	// }

	return nil
}

func (c *HrsvClient) Mutate(endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.URL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setCommonHeaders(req)
	resp, err := c.Session.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	return c.parseResponse(resp.Body)
}


func (c *HrsvClient) Query(endpoint string) (map[string]interface{}, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.URL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setCommonHeaders(req)
	resp, err := c.Session.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	return c.parseResponse(resp.Body)
}


func (c *HrsvClient) ensureAuthenticated() error {
	if !c.IsAuthenticated {
		return c.Authenticate()
	}
	return nil
}

func (c *HrsvClient) setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)
}

func (c *HrsvClient) checkResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *HrsvClient) parseResponse(body io.Reader) (map[string]interface{}, error) {
	responseBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}


// func (c *HrsvClient) Mutate(endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
// 	if !c.IsAuthenticated {
// 		if err := c.Authenticate(); err != nil {
// 			return nil, err
// 		}
// 	}
// 	body, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", c.URL+endpoint, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Add("Accept", "application/json")
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Authorization", c.token)
// 	resp, err := c.Session.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result map[string]interface{}
// 	err = json.Unmarshal(responseBody, &result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func (c *HrsvClient) Query(endpoint string) (map[string]interface{}, error) {
// 	if !c.IsAuthenticated {
// 		if err := c.Authenticate(); err != nil {
// 			return nil, err
// 		}
// 	}
// 	req, err := http.NewRequest("GET", c.URL+endpoint, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("Accept", "application/json")
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Authorization", c.token)
// 	resp, err := c.Session.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result := make(map[string]interface{})
// 	err = json.Unmarshal(responseBody, &result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

