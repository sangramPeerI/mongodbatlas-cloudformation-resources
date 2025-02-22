package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	log "github.com/mongodb/mongodbatlas-cloudformation-resources/util/logger"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progressevent"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	mongodbatlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	CallBackSeconds = 30
)

var CreateRequiredFields = []string{constants.ProjID, constants.PvtKey, constants.PubKey, constants.Name}
var ReadRequiredFields = []string{constants.ProjID, constants.Name, constants.PvtKey, constants.PubKey}
var UpdateRequiredFields = []string{constants.PvtKey, constants.PubKey}
var DeleteRequiredFields = []string{constants.ProjID, constants.Name, constants.PvtKey, constants.PubKey}
var ListRequiredFields = []string{constants.PvtKey, constants.PubKey}

func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-ServerlessInstance")
}

func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	_, _ = log.Debugf("Create() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(CreateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		_, _ = log.Warnf("Create - error: %+v", err)
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	// Callback
	if stateName, ok := req.CallbackContext[constants.StateName]; ok {
		_, _ = log.Debugf("Callback state: %s", stateName)
		return serverlessCallback(client, currentModel, constants.IdleState)
	}

	serverlessInstanceRequest := &mongodbatlas.ServerlessCreateRequestParams{
		Name:                         *currentModel.Name,
		ProviderSettings:             setProviderSettings(currentModel),
		ServerlessBackupOptions:      setBackupOptions(currentModel),
		TerminationProtectionEnabled: currentModel.TerminationProtectionEnabled,
	}

	serverless, res, err := client.ServerlessInstances.Create(context.Background(), *currentModel.ProjectID, serverlessInstanceRequest)
	if err != nil {
		_, _ = log.Warnf("Serverless - Create() - error: %+v", err)
		if res.Response.StatusCode == 400 && strings.Contains(err.Error(), constants.Duplicate) {
			return handler.ProgressEvent{
				OperationStatus:  handler.Failed,
				Message:          err.Error(),
				HandlerErrorCode: cloudformation.HandlerErrorCodeAlreadyExists}, nil
		}
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              fmt.Sprintf("Create ServerlessInstance `%s`", serverless.StateName),
		ResourceModel:        currentModel,
		CallbackDelaySeconds: CallBackSeconds,
		CallbackContext: map[string]interface{}{
			constants.StateName: serverless.StateName,
		},
	}, nil
}

func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	_, _ = log.Debugf("Read() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(ReadRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		_, _ = log.Warnf("Read - error: %+v", err)
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	cluster, res, err := client.ServerlessInstances.Get(context.Background(), *currentModel.ProjectID, *currentModel.Name)
	if err != nil {
		_, _ = log.Warnf("Read - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	// Read Instance
	model := readServerlessInstance(cluster)
	model.ApiKeys = &ApiKeyDefinition{
		PrivateKey: currentModel.ApiKeys.PrivateKey,
		PublicKey:  currentModel.ApiKeys.PublicKey,
	}
	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   model,
	}, nil
}

func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	_, _ = log.Debugf("Update() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(UpdateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		_, _ = log.Warnf("Update - error: %+v", err)
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	// Callback
	if _, ok := req.CallbackContext[constants.StateName]; ok {
		return serverlessCallback(client, currentModel, constants.IdleState)
	}

	// CFN TEST : currently Update is throwing 500 Error instead of 404 if resource not exists
	_, res, err := client.ServerlessInstances.Get(context.Background(), *currentModel.ProjectID, *currentModel.Name)
	if err != nil {
		_, _ = log.Warnf("Read in Update - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	serverlessInstanceRequest := &mongodbatlas.ServerlessUpdateRequestParams{
		ServerlessBackupOptions:      setBackupOptions(currentModel),
		TerminationProtectionEnabled: currentModel.TerminationProtectionEnabled,
	}

	serverless, res, err := client.ServerlessInstances.Update(context.Background(), *currentModel.ProjectID, *currentModel.Name, serverlessInstanceRequest)
	if err != nil {
		_, _ = log.Warnf("Update - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	// Response
	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              fmt.Sprintf("Create ServerlessInstance `%s`", serverless.StateName),
		ResourceModel:        currentModel,
		CallbackDelaySeconds: CallBackSeconds,
		CallbackContext: map[string]interface{}{
			constants.StateName: serverless.StateName,
			constants.ID:        serverless.ID,
		},
	}, nil
}

func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	_, _ = log.Debugf("Delete() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(DeleteRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		_, _ = log.Warnf("Delete - error: %+v", err)
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}
	// Callback
	if _, ok := req.CallbackContext[constants.StateName]; ok {
		return serverlessCallback(client, currentModel, constants.DeletedState)
	}

	res, err := client.ServerlessInstances.Delete(context.Background(), *currentModel.ProjectID, *currentModel.Name)
	if err != nil {
		_, _ = log.Warnf("Delete - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	// Response
	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              "Deleting ServerlessInstance",
		ResourceModel:        currentModel,
		CallbackDelaySeconds: CallBackSeconds,
		CallbackContext: map[string]interface{}{
			constants.StateName: constants.DeletingState,
		},
	}, nil
}

func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	_, _ = log.Debugf("List() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(ListRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		_, _ = log.Warnf("List - error: %+v", err)
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}
	listOptions := &mongodbatlas.ListOptions{
		PageNum:      0,
		ItemsPerPage: 1000,
	}
	clustersResp, res, err := client.ServerlessInstances.List(context.Background(), *currentModel.ProjectID, listOptions)
	if err != nil {
		_, _ = log.Warnf("List - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	instances := []interface{}{} // cfn test needs empty array instead nil, when items entries found
	for i := range clustersResp.Results {
		cluster := readServerlessInstance(clustersResp.Results[i])

		instances = append(instances, cluster)
	}
	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   instances,
	}, nil
}

func setBackupOptions(currentModel *Model) (serverlessBackupOptions *mongodbatlas.ServerlessBackupOptions) {
	if currentModel.ContinuousBackupEnabled == nil {
		return nil
	}
	serverlessBackupOptions = &mongodbatlas.ServerlessBackupOptions{
		ServerlessContinuousBackupEnabled: currentModel.ContinuousBackupEnabled,
	}
	return serverlessBackupOptions
}

func setProviderSettings(currentModel *Model) (serverlessProviderSettings *mongodbatlas.ServerlessProviderSettings) {
	if currentModel.ProviderSettings == nil {
		return &mongodbatlas.ServerlessProviderSettings{
			ProviderName:        constants.Serverless,
			BackingProviderName: constants.AWS,
		}
	}
	serverlessProviderSettings = &mongodbatlas.ServerlessProviderSettings{
		BackingProviderName: constants.AWS,
	}
	if currentModel.ProviderSettings.ProviderName != nil {
		serverlessProviderSettings.ProviderName = *currentModel.ProviderSettings.ProviderName
	}
	if currentModel.ProviderSettings.RegionName != nil {
		serverlessProviderSettings.RegionName = *currentModel.ProviderSettings.RegionName
	}
	return serverlessProviderSettings
}

func readServerlessInstance(cluster *mongodbatlas.Cluster) (serverless *Model) {
	serverless = &Model{}
	serverless.Name = &cluster.Name
	serverless.Id = &cluster.ID
	serverless.ProjectID = &cluster.GroupID
	if cluster.ProviderSettings != nil {
		serverless.ProviderSettings = &ServerlessInstanceProviderSettings{
			ProviderName: &cluster.ProviderSettings.ProviderName,
			RegionName:   &cluster.ProviderSettings.RegionName,
		}
	}

	if cluster.ServerlessBackupOptions != nil {
		serverless.ContinuousBackupEnabled = cluster.ServerlessBackupOptions.ServerlessContinuousBackupEnabled
	}

	if cluster.ConnectionStrings != nil {
		serverless.ConnectionStrings = &ServerlessInstanceConnectionStrings{
			StandardSrv:     &cluster.ConnectionStrings.StandardSrv,
			PrivateEndpoint: readPrivateEndpoint(cluster.ConnectionStrings.PrivateEndpoint),
		}
	}
	serverless.CreateDate = &cluster.CreateDate
	serverless.MongoDBVersion = &cluster.MongoDBVersion
	serverless.Links = readLinks(cluster.Links)
	serverless.TerminationProtectionEnabled = cluster.TerminationProtectionEnabled
	serverless.StateName = &cluster.StateName
	return serverless
}

func readPrivateEndpoint(privateEPs []mongodbatlas.PrivateEndpoint) (endPoints []ServerlessInstancePrivateEndpoint) {
	for i := range privateEPs {
		var pep = ServerlessInstancePrivateEndpoint{}
		pep.Endpoints = readPrivateEndpointEndpoints(privateEPs[i].Endpoints)
		pep.Type = &privateEPs[i].Type
		pep.SrvConnectionString = &privateEPs[i].SRVConnectionString
		endPoints = append(endPoints, pep)
	}
	return
}

func readPrivateEndpointEndpoints(peEndpoints []mongodbatlas.Endpoint) (epEndpoints []ServerlessInstancePrivateEndpointEndpoint) {
	for i := range peEndpoints {
		epEndpoints = append(epEndpoints, ServerlessInstancePrivateEndpointEndpoint{
			EndpointId:   &peEndpoints[i].EndpointID,
			ProviderName: &peEndpoints[i].ProviderName,
			Region:       &peEndpoints[i].Region,
		})
	}
	return
}

func readLinks(clsLinks []*mongodbatlas.Link) (links []Link) {
	for i := range clsLinks {
		links = append(links, Link{
			Href: &clsLinks[i].Href,
			Rel:  &clsLinks[i].Rel,
		})
	}
	return
}

func serverlessCallback(client *mongodbatlas.Client, currentModel *Model, targtStatus string) (progressEvent handler.ProgressEvent, err error) {
	serverless, resp, err := client.ServerlessInstances.Get(context.Background(), *currentModel.ProjectID, *currentModel.Name)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 && targtStatus == constants.DeletedState {
			_, _ = log.Debugf("404:No instance found")
			return handler.ProgressEvent{
				OperationStatus: handler.Success,
				Message:         "Deleted ServerlessInstance",
				ResourceModel:   nil,
			}, nil
		}
		_, _ = log.Warnf("Read - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), resp.Response), nil
	}

	_, _ = log.Debugf("stateName : %s", serverless.StateName)
	currentModel.Id = &serverless.ID
	if serverless.StateName != constants.IdleState {
		return handler.ProgressEvent{
			OperationStatus:      handler.InProgress,
			Message:              fmt.Sprintf("Create ServerlessInstance `%s`", serverless.StateName),
			ResourceModel:        currentModel,
			CallbackDelaySeconds: CallBackSeconds,
			CallbackContext: map[string]interface{}{
				constants.StateName: serverless.StateName,
			},
		}, nil
	}
	_, _ = log.Debugf("Response : %+v", serverless)

	model := readServerlessInstance(serverless)
	model.ApiKeys = &ApiKeyDefinition{
		PrivateKey: currentModel.ApiKeys.PrivateKey,
		PublicKey:  currentModel.ApiKeys.PublicKey,
	}
	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         fmt.Sprintf("Create ServerlessInstance `%s`", serverless.StateName),
		ResourceModel:   model,
	}, nil
}
