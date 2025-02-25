package handler

import (
	"net/http"
	"strconv"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/gin-gonic/gin"
)

func CreateNetworkConfig(ctx *gin.Context) {
	req := request.CreateNetworkRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	logs.Logger.Infof("req is:%v ", req)
	resp, err := service.CreateNetwork(ctx, &service.CreateNetworkRequest{
		Provider:          req.Provider,
		RegionId:          req.RegionId,
		CidrBlock:         req.CidrBlock,
		VpcName:           req.VpcName,
		ZoneId:            req.ZoneId,
		SwitchCidrBlock:   req.SwitchCidrBlock,
		GatewayIp:         req.GatewayIp,
		SwitchName:        req.SwitchName,
		SecurityGroupName: req.SecurityGroupName,
		SecurityGroupType: req.SecurityGroupType,
		AK:                req.AK,
		Rules:             req.Rules,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func SyncNetworkConfig(ctx *gin.Context) {
	req := request.SyncNetworkRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, response.ParamInvalid, nil)
		return
	}

	err = service.SyncNetwork(ctx, service.SyncNetworkRequest{
		Provider:   req.Provider,
		RegionId:   req.RegionId,
		AccountKey: req.AccountKey,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func GetNetCfgTemplate(ctx *gin.Context) {
	provider := ctx.Query("provider")
	netCfg := request.CreateNetworkRequest{
		Provider:          provider,
		CidrBlock:         "10.0.0.0/16",
		VpcName:           "默认的vpc",
		SwitchCidrBlock:   "10.0.0.0/24",
		GatewayIp:         "10.0.0.254",
		SwitchName:        "默认的子网",
		SecurityGroupName: "默认的安全组",
		SecurityGroupType: "normal",
		Rules: []service.GroupRule{
			{
				Protocol:  cloud.ProtocolAll,
				Direction: cloud.SecGroupRuleIn,
				CidrIp:    "0.0.0.0/0",
			},
			{
				Protocol:  cloud.ProtocolAll,
				Direction: cloud.SecGroupRuleOut,
				CidrIp:    "0.0.0.0/0",
			},
		},
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, netCfg)
	return
}

func CreateVpc(ctx *gin.Context) {
	req := request.CreateVpcRequest{}
	err := ctx.Bind(&req)
	if err != nil || !req.Check() {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	resp, err := service.CreateVPC(ctx, service.CreateVPCRequest{
		Provider:  req.Provider,
		RegionId:  req.RegionId,
		VpcName:   req.VpcName,
		CidrBlock: req.CidrBlock,
		AK:        req.AK,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetVpcById(ctx *gin.Context) {
	id := ctx.Param("id")
	resp, err := service.GetVpcById(ctx, id)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func DescribeVpc(ctx *gin.Context) {
	req := request.DescribeVpcRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, response.PermissionDenied, nil)
		return
	}
	pageNumber, pageSize := getPager(ctx)
	resp, err := service.GetVPC(ctx, service.GetVPCRequest{
		Provider:   req.Provider,
		RegionId:   req.RegionId,
		VpcName:    req.VpcName,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		AccountKey: req.AccountKey,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func CreateSwitch(ctx *gin.Context) {
	req := request.CreateSwitchRequest{}
	err := ctx.Bind(&req)
	if err != nil || !req.Check() {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	resp, err := service.CreateSwitch(ctx, service.CreateSwitchRequest{
		AK:         req.AK,
		SwitchName: req.SwitchName,
		ZoneId:     req.ZoneId,
		VpcId:      req.VpcId,
		CidrBlock:  req.CidrBlock,
		GatewayIp:  req.GatewayIp,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetSwitchById(ctx *gin.Context) {
	switchId := ctx.Param("id")
	vpcId := ctx.Query("vpc_id")
	resp, err := service.GetSwitchById(ctx, vpcId, switchId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func DescribeSwitch(ctx *gin.Context) {
	vpcId := ctx.Query("vpc_id")
	switchName := ctx.Query("switch_name")
	zoneId := ctx.Query("zone_id")
	pageNumber, pageSize := getPager(ctx)

	resp, err := service.GetSwitch(ctx, service.GetSwitchRequest{
		SwitchName: switchName,
		VpcId:      vpcId,
		ZoneId:     zoneId,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func CreateSecurityGroup(ctx *gin.Context) {
	req := request.CreateSecurityGroupRequest{}
	err := ctx.Bind(&req)
	if err != nil || !req.Check() {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	logs.Logger.Infof("req is:%v ", req)

	resp, err := service.CreateSecurityGroup(ctx, service.CreateSecurityGroupRequest{
		AK:                req.AK,
		VpcId:             req.VpcId,
		SecurityGroupName: req.SecurityGroupName,
		SecurityGroupType: req.SecurityGroupType,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func DescribeSecurityGroup(ctx *gin.Context) {
	ak := ctx.Query("account_key")
	vpcId := ctx.Query("vpc_id")
	securityGroupName := ctx.Query("security_group_name")
	pageNumber, pageSize := getPager(ctx)

	resp, err := service.GetSecurityGroup(ctx, service.GetSecurityGroupRequest{
		AK:                ak,
		SecurityGroupName: securityGroupName,
		VpcId:             vpcId,
		PageNumber:        pageNumber,
		PageSize:          pageSize,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func AddSecurityGroupRule(ctx *gin.Context) {
	req := request.AddSecurityGroupRuleRequest{}
	err := ctx.Bind(&req)
	if err != nil || !req.Check() {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	logs.Logger.Infof("req is:%v ", req)

	resp, err := service.AddSecurityGroupRule(ctx, service.AddSecurityGroupRuleRequest{
		AK:              req.AK,
		RegionId:        req.RegionId,
		VpcId:           req.VpcId,
		SecurityGroupId: req.SecurityGroupId,
		Rules:           req.Rules,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func CreateSecurityGroupWithRules(ctx *gin.Context) {
	req := request.CreateSecurityGroupWithRuleRequest{}
	err := ctx.Bind(&req)
	if err != nil || !req.Check() {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	groupId, err := service.CreateSecurityGroup(ctx, service.CreateSecurityGroupRequest{
		AK:                req.AK,
		VpcId:             req.VpcId,
		SecurityGroupName: req.SecurityGroupName,
		SecurityGroupType: req.SecurityGroupType,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if len(req.Rules) == 0 || req.Rules[0].Protocol == "" {
		response.MkResponse(ctx, http.StatusOK, response.Success, groupId)
		return
	}
	_, err = service.AddSecurityGroupRule(ctx, service.AddSecurityGroupRuleRequest{
		AK:              req.AK,
		RegionId:        req.RegionId,
		VpcId:           req.VpcId,
		SecurityGroupId: groupId,
		Rules:           req.Rules,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, groupId)
}

func GetSecurityGroupWithRules(ctx *gin.Context) {
	secGrpId := ctx.Param("id")
	resp, err := service.GetSecurityGroupWithRules(ctx, secGrpId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func getPager(ctx *gin.Context) (pageNumber int, pageSize int) {
	pageNumber, _ = strconv.Atoi(ctx.Query("page_number"))
	if pageNumber < 1 {
		pageNumber = 1
	}
	pageSize, _ = strconv.Atoi(ctx.Query("page_size"))
	if pageSize < 1 || pageSize > constants.DefaultPageSize {
		pageSize = constants.DefaultPageSize
	}
	return pageNumber, pageSize
}

func GetOrgKeys(ctx *gin.Context) (*types.OrgKeys, error) {
	user := helper.GetUserClaims(ctx)
	return service.GetAccountsByOrgId(user.GetOrgIdForTest())
}
