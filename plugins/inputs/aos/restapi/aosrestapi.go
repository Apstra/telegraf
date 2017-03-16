// Copyright 2014-present, Apstra, Inc. All rights reserved.
//
// This source code is licensed under End User License Agreement found in the
// LICENSE file at http://www.apstra.com/community/eula

package aosrestapi

import "fmt"
import "time"
import "bytes"
import "encoding/json"
import "net/http"
import "io/ioutil"
import "errors"

var client = &http.Client{Timeout: 10 * time.Second}

type AosToken struct {
	Token string
}

type aosSystemList struct {
		Items				[]aosSystem 	`json:"items"`
}

type aosBlueprintList struct {
		Items				[]aosBlueprintSummary 	`json:"items"`
}

func (bplist *aosBlueprintList) GetSystemFromBp( bpId string, systemId string ) *aosBlueprintSystemNode {

    // Search for Node based on Id on a specific Blueprint
    for _, bp := range bplist.Items {
      if bp.Id == bpId {
        // fmt.Printf("Blueprint Find : %s \n", bp.Name)

        for _, node := range bp.Content.System.Nodes {
          // fmt.Printf("  Node : %s \n", node.Name)

          if node.Id == systemId {
            return &node
          }
        }
      }
    }

    return nil

}

type aosSystem struct {
		DeviceKey		string 							`json:"device_key"`
		Facts				aosSystemFacts 			`json:"facts"`
		Id 					string		 					`json:"id"`
		Status			aosSystemStatus			`json:"status"`
		UserConfig	aosSystemUserConfig `json:"user_config"`
    Blueprint   aosBlueprintSystemNode
}

type aosSystemFacts struct {
		AosHclModel 	string 	  `json:"aos_hcl_model"`
		AosServer			string 		`json:"aos_server"`
		AosVersion		string 		`json:"aos_version"`
		HwModel		 		string		`json:"hw_model"`
		HwVersion		 	string		`json:"hw_version"`
		MgmtIfName		string		`json:"mgmt_ifname"`
		MgmtIpAddr		string		`json:"mgmt_ipaddr"`
		MgmtMacAddr	 	string		`json:"mgmt_macaddr"`
		OsArch				string		`json:"os_arch"`
		OsFamily			string		`json:"os_family"`
		OsVersion			string		`json:"os_version"`
		// OsVersionInfo aosSystemVersionInfo	`json:"os_version_info"`
		SerialNumber  string		`json: "serial_number"`
		Vendor			  string		`json:"vendor"`
}

type aosSystemStatus struct {
		AgentStartTime 	string		`json:"agent_start_time"`
		CommState				string		`json:"comm_state"`
    BlueprintActive bool      `json:"blueprint_active"`
    BlueprintId     string    `json:"blueprint_id"`
		DeviceStartTime	string		`json:"device_start_time"`
		DomainName 			string		`json:"domain_name"`
		ErrorMessage		string		`json:"error_message"`
		Fqdn 						string		`json:"fqdn"`
		Hostname				string		`json:"hostname"`
		IsAcknowledged	bool			`json:"is_acknowledged"`
		PoolId					string		`json:"pool_id"`
		State						string		`json:"state"`
}

type aosSystemUserConfig struct {
		AdminState		string 	`json:"admin_state"`
		AosHclModel		string 	`json:"aos_hcl_model"`
		Location			string 	`json:"location"`
}

// type aosSystemVersionInfo struct {
// 		Build			string 	`json:"build"`
// 		Major			int 		`json:"major"`
// 		Minor			int			`json:"minor"`
// }

type aosBlueprintSummary struct {
		Id		       string   `json:"id"`
		ReferenceArchitecture	string `json:"reference_architecture"`
    Name	       string   `json:"display_name"`
    Content      aosBlueprintDetail
}

type aosBlueprintDetail struct {
		Id		        string  `json:"id"`
		ReferenceArchitecture	string `json:"reference_architecture"`
    Name	        string   `json:"display_name"`
    System        aosBlueprintSystem  `json:"system"`
}

type aosBlueprintSystem struct {
    Nodes         []aosBlueprintSystemNode
    EdgeIpConnectivity  aosBlueprintSystemEdgeIp `json:"l3_edge_ip_connectivity_groups"`
}


type aosBlueprintSystemEdgeIp struct {
    ExternalSubnets []interface{}   `json:"external_subnets"`
    InternalSubnets []interface{}   `json:"internal_subnets"`
}


type aosBlueprintSystemNode struct {
  Hostname        string  `json:"hostname"`
  Name            string  `json:"display_name"`
  LoopbackId      string  `json:"loopback_ip"`
  Id              string  `json:"id"`
  HclId           string  `json:"hcl_id"`
  Role            string  `json:"role"`
  DeviceConfigRendering   string  `json:"device_config_rendering"`
  DeployMode      string  `json:"deploy_mode"`
  Position        int  `json:"position"`
  Asn             int     `json:"asn"`
}

type AosServerApi  struct {
  Address     string
  Port        int
  User        string
  Password    string

  Token       string
  Blueprints  map[string]aosBlueprintSummary
  Systems     map[string]aosSystem
  StreamingSessions []string

}

type apiResponseId struct{
    Id string
}

func NewAosServerApi (address string, port int, user string, password string) *AosServerApi  {

    api := AosServerApi { Address: address, Port: port,User: user,Password: password}

    // initialize Maps
    api.Blueprints = make(map[string]aosBlueprintSummary, 20)
    api.Systems = make(map[string]aosSystem, 100)

    return &api
}

func (api *AosServerApi ) Login() (err error) {

  url := fmt.Sprintf("http://%v:%v/api/user/login", api.Address, api.Port)
  // fmt.Println("Login() URL:>", url)

  var jsonStr = []byte(fmt.Sprintf(`{ "username": "%v", "password": "%v" }`, api.User, api.Password))

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/json")

  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != 201 {
    panic(resp.Status)
  }

  token := AosToken{}
  json.NewDecoder(resp.Body).Decode(&token)

  api.Token = token.Token
  // fmt.Printf("Login() Token: %+v \n", token.Token)

  return nil
}

func (api *AosServerApi ) httpRequest(httpType string, address string, v interface{}, expectedCode int) error {

  url := fmt.Sprintf("http://%v:%v/api/%v", api.Address, api.Port, address)

  req, err := http.NewRequest(httpType, url, nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("AUTHTOKEN", api.Token)

  resp, err := client.Do(req)
  if err != nil {
      return err
  }
  defer resp.Body.Close()

  // Check if Return code is what we are expecting
  if resp.StatusCode != expectedCode {
     return errors.New(fmt.Sprintf("Status Code is not %v got %v", expectedCode, resp.Status))
  }

  // Read Body and Unmarshal JSON
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil { return err  }

  err = json.Unmarshal(body, v)
  if err != nil {  return err }

  return nil
}

func (api *AosServerApi) StopStreaming() error {

  for _, id := range api.StreamingSessions {

    url := fmt.Sprintf("http://%v:%v/api/streaming-config/%v", api.Address, api.Port, id)

    req, err := http.NewRequest("DELETE", url,  nil)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("AUTHTOKEN", api.Token)

    resp, err := client.Do(req)
    if err != nil {  return err  }

    defer resp.Body.Close()

    if resp.StatusCode != 202 {
      return errors.New(resp.Status)
    }
  }

  return nil
}

func (api *AosServerApi) StartStreaming(streamingType string, address string, port int) error {

  url := fmt.Sprintf("http://%v:%v/api/streaming-config", api.Address, api.Port)

  var jsonStr = []byte(fmt.Sprintf(`{
        "streaming_type": "%v",
        "hostname": "%v",
        "protocol": "protoBufOverTcp",
        "port": %v
      }`, streamingType, address, port))

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("AUTHTOKEN", api.Token)

  resp, err := client.Do(req)
  if err != nil {  return err  }

  defer resp.Body.Close()

  if resp.StatusCode != 201 {
    panic(resp.Status)
  }

  stramingResp := apiResponseId{}

  json.NewDecoder(resp.Body).Decode(&stramingResp)

  api.StreamingSessions = append(api.StreamingSessions, stramingResp.Id)
  // fmt.Printf("Created Streaming Session : %v\n", stramingResp.Id)

  return nil
}

func (api *AosServerApi ) GetBlueprints() error {

  blueprintList := aosBlueprintList{}
  err := api.httpRequest("GET", "blueprints", &blueprintList, 200)
  if err != nil { return err }

  // Save all items in the API Object by ID
  for i := 0; i < len(blueprintList.Items); i++ {

    id := blueprintList.Items[i].Id
    // fmt.Printf("GetBlueprints() - Id %v \n", id )

    api.Blueprints[id] = blueprintList.Items[i]

    // fmt.Printf("Id: %s - %s\n", api.Blueprints[id].Id, api.Blueprints[id].Name )

    // Not possible to assign a field inside a struct directly in if part of a map
    // Must same stuct temporarely
    tmp := api.Blueprints[id]

    err = api.httpRequest("GET", "blueprints/" + api.Blueprints[id].Id, &tmp.Content , 200)
    if err != nil { return err }

    api.Blueprints[id] = tmp
  }

  return nil
}

func (api *AosServerApi ) GetSystems() error {

  systemList := aosSystemList{}

  err := api.httpRequest("GET", "systems", &systemList, 200)
  if err != nil { return err }

  for _, system := range systemList.Items {

    id := system.Id

    api.Systems[id] = system
    s := api.Systems[id]

    // If Blueprint is defined, search for the System.Node information and copy it inside system
    if system.Status.BlueprintId != "" {

      if bp, ok := api.Blueprints[system.Status.BlueprintId]; ok {
        var found bool

        for _, node := range bp.Content.System.Nodes {

          // fmt.Printf("GetSystems() node - Id %v \n", node.Id )
          //
          if node.Id == id  {
            s.Blueprint = node
            api.Systems[id] = s
            found = true
          }
        }
        if found == false {
          return errors.New(fmt.Sprintf("System %v has Blueprint ID defined (%v) but was not able to find it in System.Nodes", id, system.Status.BlueprintId))
        }
      } else {
        return errors.New(fmt.Sprintf("System %v has Blueprint ID defined (%v) but not blueprint with this is exist in map", id, system.Status.BlueprintId))
      }
    }
  }
  return nil
}

func (api *AosServerApi ) GetSystemByKey( deviceKey string ) *aosSystem {

    for _, system := range api.Systems {

      if system.DeviceKey == deviceKey {
        return &system
      }
    }

    return nil
}
