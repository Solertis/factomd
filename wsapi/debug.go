// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wsapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/web"
)

func HandleDebug(ctx *web.Context) {
	ServersMutex.Lock()
	state := ctx.Server.Env["state"].(interfaces.IState)
	ServersMutex.Unlock()

	if err := checkAuthHeader(state, ctx.Request); err != nil {
		remoteIP := ""
		remoteIP += strings.Split(ctx.Request.RemoteAddr, ":")[0]
		fmt.Printf("Unauthorized V2 API client connection attempt from %s\n", remoteIP)
		ctx.ResponseWriter.Header().Add("WWW-Authenticate", `Basic realm="factomd RPC"`)
		http.Error(ctx.ResponseWriter, "401 Unauthorized.", http.StatusUnauthorized)

		return
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		HandleV2Error(ctx, nil, NewInvalidRequestError())
		return
	}

	j, err := primitives.ParseJSON2Request(string(body))
	if err != nil {
		HandleV2Error(ctx, nil, NewInvalidRequestError())
		return
	}

	jsonResp, jsonError := HandleDebugRequest(state, j)

	if jsonError != nil {
		HandleV2Error(ctx, j, jsonError)
		return
	}

	ctx.Write([]byte(jsonResp.String()))
}

func HandleDebugRequest(state interfaces.IState, j *primitives.JSON2Request) (*primitives.JSON2Response, *primitives.JSONError) {
	var resp interface{}
	var jsonError *primitives.JSONError
	params := j.Params
	switch j.Method {
	case "audit-servers":
		resp, jsonError = HandleAuditServers(state, params)
		break
	case "configuration":
		resp, jsonError = HandleConfig(state, params)
		break
	case "current-minute":
		resp, jsonError = HandleCurrentMinute(state, params)
		break
	case "drop-rate":
		resp, jsonError = HandleDropRate(state, params)
		break
	case "federated-servers":
		resp, jsonError = HandleFedServers(state, params)
		break
	case "holding-queue":
		resp, jsonError = HandleHoldingQueue(state, params)
		break
	case "network-info":
		resp, jsonError = HandleNetworkInfo(state, params)
		break
	case "node-status":
		resp, jsonError = HandleNodeStatus(state, params)
		break
	case "predictive-fer":
		resp, jsonError = HandlePredictiveFER(state, params)
		break
	default:
		jsonError = NewMethodNotFoundError()
		break
	}
	if jsonError != nil {
		return nil, jsonError
	}

	fmt.Printf("API V2 method: <%v>  paramaters: %v\n", j.Method, params)

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp
	return jsonResp, nil
}

func HandleAuditServers(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	type ret struct {
		AuditServers []interfaces.IFctServer
	}
	r := new(ret)

	r.AuditServers = state.GetAuditServers(state.GetLeaderHeight())
	return r, nil
}

func HandleConfig(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	return state.GetCfg(), nil
}

func HandleCurrentMinute(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	return state.GetCurrentMinute(), nil
}

func HandleDropRate(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	return state.GetDropRate(), nil
}

func HandleFedServers(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	type ret struct {
		FederatedServers []interfaces.IFctServer
	}
	r := new(ret)

	r.FederatedServers = state.GetFedServers(state.GetLeaderHeight())
	return r, nil
}

func HandleHoldingQueue(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	type ret struct {
		Messages []interfaces.IMsg
	}
	r := new(ret)

	for _, v := range state.LoadHoldingMap() {
		r.Messages = append(r.Messages, v)
	}
	return r, nil
}

func HandleNetworkInfo(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	type ret struct {
		NetworkNumber int
		NetworkName   string
		NetworkID     uint32
	}
	r := new(ret)
	r.NetworkNumber = state.GetNetworkNumber()
	r.NetworkName = state.GetNetworkName()
	r.NetworkID = state.GetNetworkID()
	return r, nil
}

func HandleNodeStatus(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	type ret struct {
		Status []string
	}
	r := new(ret)
	r.Status = state.GetStatus()

	return r, nil
}

func HandlePredictiveFER(
	state interfaces.IState,
	params interface{},
) (
	interface{},
	*primitives.JSONError,
) {
	type ret struct {
		PredictiveFER uint64
	}
	r := new(ret)
	r.PredictiveFER = state.GetPredictiveFER()
	return r, nil
}
