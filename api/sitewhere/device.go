package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.sdkeji.top/share/sqmslib/api"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DeviceModule struct {
	api *SiteWhereAPI
}

type CreateNewDeviceRequest struct {
	Comments        string            `json:"comments,omitempty"`        //备注名
	DeviceTypeToken string            `json:"deviceTypeToken"` //类型token
	Metadata        map[string]string `json:"metadata,omitempty"`
	Token           string            `json:"token"` //Token
}

type Device struct {
	ID                    string              `json:"id"`
	Comments              string              `json:"comments"`    //备注名
	CreatedDate           string              `json:"createdDate"` //创建时间
	DeviceElementMappings []map[string]string `json:"deviceElementMappings"`
	DeviceTypeID          string              `json:"deviceTypeId"`
	DeviceAssignmentID    string              `json:"deviceAssignmentId"`
	Metadata              map[string]string   `json:"metadata"`
	Token                 string              `json:"token"`
}

type DeviceType struct {
	ContainerPolice string            `json:"containerPolice"`
	CreatedDate     string            `json:"createdDate"` //创建时间
	Description     string            `json:"description"` //描述
	ID              string            `json:"id"`          //ID
	ImageURL        string            `json:"imageUrl"`    //图像链接
	Metadata        map[string]string `json:"metadata"`
	Name            string            `json:"name"`        //类型名称
	Token           string            `json:"token"`       //类型token
	UpdatedDate     string            `json:"updatedDate"` //更新时间
}

type CreateDeviceTypeRequest struct {
	BackgroundColor string            `json:"backgroundColor"`
	BorderColor     string            `json:"borderColor"`
	ForegroundColor string            `json:"foregroundColor"`
	Icon            string            `json:"icon"`
	ContainerPolicy string            `json:"containerPolicy"`
	Description     string            `json:"description,omitempty"`
	ImageURL        string            `json:"imageUrl,omitempty"`
	Name            string            `json:"name"`
	Token           string            `json:"token"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type DeviceListResponse struct {
	NumResults int `json:"numResults"`
	Results    []struct {
		Device
	} `json:"results"`
}

type DeviceTypeListResponse struct {
	NumResults int `json:"numResults"`
	Results    []struct {
		DeviceType
	} `json:"results"`
}

type QueryDeviceTypesCond struct {
	IncludeAsset bool  `json:"includeAsset"`
	Page         int32 `json:"page"`
	PageSize     int32 `json:"pageSize"`
}

type QueryDevicesCond struct {
	DeviceType        string `json:"deviceType"`
	EndDate           string `json:"end_date"`
	ExcludeAssigned   bool   `json:"exclude_assigned"` //是否排除已分配
	IncludeDeviceType bool   `json:"includeDeviceType"`
	IncludeAssignment bool   `json:"includeAssignment"`
	Page              int32  `json:"page"`
	PageSize          int32  `json:"pageSize"`
}

//获取单个Device类型
func (m DeviceModule) GetDeviceTypeByToken(tokens string) (res DeviceType, err error) {
	URL, err := url.Parse(m.api.createURL("/sitewhere/api/devicetypes/" + tokens))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	reqs, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(reqs, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		//m.api.Warn("received error response.", "response", string(b))
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//获取Device类型列表
func (m DeviceModule) GetDeviceTypeList(req QueryDeviceTypesCond) (deviceTypes DeviceTypeListResponse, err error) {
	URL, err := url.Parse(m.api.createURL("/sitewhere/api/devicetypes"))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	q, err := query.Values(req)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	URL.RawQuery = q.Encode()

	reqs, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(reqs, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		//m.api.Warn("received error response.", "response", string(b))
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}
	err = json.Unmarshal(b, &deviceTypes)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//创建Device类型
func (m DeviceModule) CreateDeviceType(request CreateDeviceTypeRequest) (types DeviceType, err error) {

	b, err := json.Marshal(request)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	URL := m.api.createURL("/sitewhere/api/devicetypes")
	m.api.Debug("create new device type.", "url", URL)
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(b))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(req, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}

	err = json.Unmarshal(b, &types)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//删除Device类型
func (m DeviceModule) DeleteDeviceType(typeToken string) (types DeviceType, err error) {
	URL, err := url.Parse(m.api.createURL("/sitewhere/api/devicetypes/" + typeToken))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	reqs, err := http.NewRequest(http.MethodDelete, URL.String(), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(reqs, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		//m.api.Warn("received error response.", "response", string(b))
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}
	err = json.Unmarshal(b, &types)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//更新Device类型
func (m DeviceModule) UpdateDeviceType(typeToken string, request CreateDeviceTypeRequest) (types DeviceType, err error) {

	b, err := json.Marshal(request)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	URL := m.api.createURL("/sitewhere/api/devicetypes/" + typeToken)
	m.api.Debug("update device types.", "url", URL)
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(b))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(req, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}

	err = json.Unmarshal(b, &types)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//创建Device
func (m DeviceModule) CreateDevice(request CreateNewDeviceRequest) (res Device, err error) {

	b, err := json.Marshal(request)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	URL := m.api.createURL("/sitewhere/api/devices")
	m.api.Debug("create new device.", "url", URL)
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(b))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(req, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}

	err = json.Unmarshal(b, &res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//删除Device
func (m DeviceModule) DeleteDevice(typeToken string) (res Device, err error) {
	URL, err := url.Parse(m.api.createURL("/sitewhere/api/devices/" + typeToken))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	reqs, err := http.NewRequest(http.MethodDelete, URL.String(), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(reqs, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		//m.api.Warn("received error response.", "response", string(b))
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//获取Device列表
func (m DeviceModule) GetDeviceList(req QueryDevicesCond) (res DeviceListResponse, err error) {
	URL, err := url.Parse(m.api.createURL("/sitewhere/api/devices"))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	q, err := query.Values(req)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	URL.RawQuery = q.Encode()
	reqs, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(reqs, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		//m.api.Warn("received error response.", "response", string(b))
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//更新Device
func (m DeviceModule) UpdateDevice(typeToken string, request CreateNewDeviceRequest) (res Device, err error) {

	b, err := json.Marshal(request)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	URL := m.api.createURL("/sitewhere/api/devices/" + typeToken)
	m.api.Debug("update device.", "url", URL)
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(b))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(req, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}

	err = json.Unmarshal(b, &res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//获取单个Device
func (m DeviceModule) GetDevice(typeToken string) (res Device, err error) {
	URL, err := url.Parse(m.api.createURL("/sitewhere/api/devices/" + typeToken))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	reqs, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	token, err := m.api.auth.Authorization("admin", "password")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	key := Bearer + " " + token

	resp, err := m.api.do(reqs, key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if resp.StatusCode != 200 {
		//m.api.Warn("received error response.", "response", string(b))
		var result api.APIError
		err = json.Unmarshal(b, &result)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		err = result.WithStatus(resp.StatusCode)
		return
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
