package routers

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/handler"
	gf_cluster "github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster"
	"github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/authorization"
	"github.com/galaxy-future/BridgX/config"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	var router *gin.Engine
	if config.GlobalConfig.DebugMode {
		gin.SetMode(gin.DebugMode)
		router = gin.Default()
		//visit http://0.0.0.0:9090/debug/pprof/
		pprof.Register(router)
	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.Default()
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello world!")
	})

	router.GET("/api/v1/galaxy_cloud/cluster/config", cluster.HandleGetClusterConfigInfoByName)

	user := router.Group("/user")
	{
		user.POST("/login", handler.Login)
		user.Use(authorization.RefreshTokenConditionCheck()).POST("/refresh_token", handler.RefreshToken)
	}

	v1Api := router.Group("/api/v1/")
	v1Api.Use(authorization.CheckTokenAuth())

	{
		cloudAccountPath := v1Api.Group("cloud_account/")
		{
			cloudAccountPath.POST("create", handler.CreateCloudAccount)
			cloudAccountPath.GET("list", handler.ListCloudAccounts)
			cloudAccountPath.GET("info", handler.GetAccountInfo)
			cloudAccountPath.POST("edit", handler.EditCloudAccount)
			cloudAccountPath.DELETE("delete/:ids", handler.DeleteCloudAccount)
		}

		clusterPath := v1Api.Group("cluster/")
		{
			clusterPath.GET("id/:id", handler.GetClusterById)
			clusterPath.GET("name/:name", handler.GetClusterByName)
			clusterPath.POST("create", handler.CreateCluster)
			clusterPath.POST("create_custom_public", handler.CreateCustomPublic)
			clusterPath.POST("create_custom_private", handler.CreateCustomPrivate)
			clusterPath.POST("edit", handler.EditCluster)
			clusterPath.DELETE("delete/:ids", handler.DeleteClusters)
			clusterPath.GET("num", handler.GetClusterCount)
			clusterPath.GET("instance_stat", handler.GetInstanceStat)
			clusterPath.GET("describe_all", handler.ListClusters)
			clusterPath.GET("custom/detail", handler.CustomClusterDetail)
			clusterPath.GET("auth", handler.GetClusterAuthByName)

			clusterPath.POST("list_by_tags", handler.ListClustersByTags)
			clusterPath.GET("get_tags", handler.GetClusterTags)
			clusterPath.POST("add_tags", handler.AddClusterTags)
			clusterPath.POST("edit_tags", handler.EditClusterTags)
			clusterPath.DELETE("delete_tags", handler.DeleteClusterTags)

			clusterPath.POST("expand", handler.ExpandCluster)
			clusterPath.POST("shrink", handler.ShrinkCluster)
			clusterPath.POST("shrink_all", handler.ShrinkAllInstances)

			clusterPath.POST("instance/check", handler.CheckInstanceConnectable)
		}
		vpcPath := v1Api.Group("vpc/")
		{
			vpcPath.GET("info/:id", handler.GetVpcById)
			vpcPath.POST("create", handler.CreateVpc)
			vpcPath.GET("describe", handler.DescribeVpc)
		}
		subnetPath := v1Api.Group("subnet/")
		{
			subnetPath.GET("info/:id", handler.GetSwitchById)
			subnetPath.POST("create", handler.CreateSwitch)
			subnetPath.GET("describe", handler.DescribeSwitch)
		}
		groupPath := v1Api.Group("security_group/")
		{
			groupPath.POST("create", handler.CreateSecurityGroup)
			groupPath.GET("describe", handler.DescribeSecurityGroup)
			groupPath.POST("rule/add", handler.AddSecurityGroupRule)
			groupPath.POST("create_with_rule", handler.CreateSecurityGroupWithRules)
			groupPath.GET(":id/rules", handler.GetSecurityGroupWithRules)
		}
		networkPath := v1Api.Group("network_config/")
		{
			networkPath.POST("create", handler.CreateNetworkConfig)
			networkPath.POST("sync", handler.SyncNetworkConfig)
			networkPath.GET("template", handler.GetNetCfgTemplate)
		}
		regionPath := v1Api.Group("region/")
		{
			regionPath.GET("list", handler.ListRegions)
		}
		zonePath := v1Api.Group("zone/")
		{
			zonePath.GET("list", handler.ListZones)
		}
		instanceTypePath := v1Api.Group("instance_type/")
		{
			instanceTypePath.GET("list", handler.ListInstanceType)
		}
		instancePath := v1Api.Group("instance/")
		{
			instancePath.GET("num", handler.GetInstanceCount)
			instancePath.GET("id/describe", handler.GetInstance)
			instancePath.GET("describe_all", handler.GetInstanceList)
			instancePath.GET("list_custom", handler.GetCustomInstanceList)
			instancePath.GET("usage_total", handler.GetInstanceUsageTotal)
			instancePath.GET("usage_statistics", handler.GetInstanceUsageStatistics)
			instancePath.POST("sync_expire_time", handler.SyncInstanceExpireTime)
		}
		taskPath := v1Api.Group("task/")
		{
			taskPath.GET("num", handler.GetTaskCount)
			taskPath.GET("list", handler.GetTaskList)
			taskPath.GET("describe", handler.GetTaskDescribe)
			taskPath.GET("describe_all", handler.GetTaskDescribeAll)
			taskPath.GET("instances", handler.GetTaskInstances)
		}
		userPath := v1Api.Group("user/")
		{
			userPath.GET("info", handler.GetUserInfo)
			userPath.POST("create_user", handler.CreateUser)
			userPath.POST("modify_password", handler.ModifyAdminPassword)
			userPath.POST("modify_username", handler.ModifyUsername)
			userPath.POST("modify_user_type", handler.ModifyUsertype)
			userPath.POST("enable_user", handler.EnableUser)
			userPath.GET("list", handler.ListUsers)
		}
		orgPath := v1Api.Group("org/")
		{
			orgPath.POST("create", handler.CreateOrg)
			orgPath.POST("edit", handler.EditOrg)
			orgPath.GET("list", handler.ListOrgs)
			orgPath.GET("id/:id", handler.GetOrgById)
		}
		keyPairPath := v1Api.Group("key_pair/")
		{
			keyPairPath.POST("create", handler.CreateKeyPair)
			keyPairPath.POST("import", handler.ImportKeyPair)
			keyPairPath.GET("info/:id", handler.GetKeyPair)
			keyPairPath.GET("list", handler.ListKeyPairs)
		}
		imagePath := v1Api.Group("image/")
		{
			imagePath.GET("list", handler.GetImageList)
		}
		logPath := v1Api.Group("log/")
		{
			logPath.GET("extract", handler.ExtractLog)
		}

		gfCluster := v1Api.Group("galaxy_cloud")
		gf_cluster.RegisterHandler(gfCluster)
	}
	return router
}
