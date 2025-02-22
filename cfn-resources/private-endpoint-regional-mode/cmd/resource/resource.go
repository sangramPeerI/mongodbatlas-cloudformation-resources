package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/logger"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progressevent"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	"go.mongodb.org/atlas/mongodbatlas"
)

var CreateRequiredFields = []string{constants.ProjectID, constants.PubKey, constants.PvtKey}
var ReadRequiredFields = []string{constants.ProjectID, constants.PubKey, constants.PvtKey}
var UpdateRequiredFields []string
var DeleteRequiredFields = []string{constants.ProjectID, constants.PubKey, constants.PvtKey}
var ListRequiredFields = []string{constants.ProjectID, constants.PubKey, constants.PvtKey}

func setup() {
	util.SetupLogger("mongodb-atlas-private-endpoint-regional-mode")
}

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	if errEvent := validator.ValidateModel(CreateRequiredFields, currentModel); errEvent != nil {
		_, _ = logger.Warnf("Validation Error")
		return *errEvent, nil
	}
	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}
	if isRegModeSettingExists(currentModel, mongodbClient) {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Regionalized Setting for Private Endpoint already enabled for : %s", *currentModel.ProjectId),
			cloudformation.HandlerErrorCodeAlreadyExists), nil
	}
	enabled := true
	currentModel.Enabled = &enabled
	// API call to Add Regional Mode for Private Endpoint
	return resourcePrivateEndpointRegionalModeUpdate(req, prevModel, currentModel, mongodbClient)
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	if errEvent := validator.ValidateModel(ReadRequiredFields, currentModel); errEvent != nil {
		_, _ = logger.Warnf("Validation Error")
		return *errEvent, nil
	}

	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}
	regPrivateEndpointSetting, response, err := mongodbClient.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), *currentModel.ProjectId)
	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error reading  : %s", err.Error()),
			response.Response), nil
	}
	enabled := regPrivateEndpointSetting.Enabled
	if !enabled {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Regionalized Setting for Private Endpoint not found for Project : %s", *currentModel.ProjectId),
			cloudformation.HandlerErrorCodeNotFound), nil
	}
	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "READ Complete",
		ResourceModel:   regionalPrivateEndpointToModel(*currentModel, regPrivateEndpointSetting),
	}, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("not implemented: Update")
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	if errEvent := validator.ValidateModel(DeleteRequiredFields, currentModel); errEvent != nil {
		_, _ = logger.Warnf("Validation Error")
		return *errEvent, nil
	}
	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}
	if isRegModeSettingExists(currentModel, mongodbClient) {
		enabled := false
		currentModel.Enabled = &enabled
		events, err := resourcePrivateEndpointRegionalModeUpdate(req, prevModel, currentModel, mongodbClient)
		if err != nil {
			return progress_events.GetFailedEventByCode(fmt.Sprintf("Error in disabling regionalized mode for private endpoint for Project : %s", *currentModel.ProjectId),
				events.HandlerErrorCode), nil
		}
		// Response
		return handler.ProgressEvent{
			OperationStatus: handler.Success,
			Message:         "Delete Complete",
		}, nil
	}
	return progress_events.GetFailedEventByCode(fmt.Sprintf("Error in disabling regionalized mode for private endpoint for Project : %s", *currentModel.ProjectId),
		cloudformation.HandlerErrorCodeNotFound), nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("not implemented: List")
}

func resourcePrivateEndpointRegionalModeUpdate(req handler.Request, prevModel *Model, currentModel *Model, client *mongodbatlas.Client) (handler.ProgressEvent, error) {
	regionalizedPrivateEndpointSetting, response, err := client.PrivateEndpoints.UpdateRegionalizedPrivateEndpointSetting(context.Background(), *currentModel.ProjectId, *currentModel.Enabled)
	if err != nil {
		return progress_events.GetFailedEventByResponse(
			fmt.Sprintf("Error in enabling regionalized settings : %s", err.Error()),
			response.Response), nil
	}
	currentModel.Enabled = &regionalizedPrivateEndpointSetting.Enabled
	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Create Complete",
		ResourceModel:   regionalPrivateEndpointToModel(*currentModel, regionalizedPrivateEndpointSetting),
	}, nil
}

func isRegModeSettingExists(currentModel *Model, client *mongodbatlas.Client) bool {
	var isExists bool
	regModeSetting, _, err := client.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), *currentModel.ProjectId)
	if err != nil {
		return isExists
	}
	if regModeSetting.Enabled {
		isExists = true
	}
	return isExists
}

func regionalPrivateEndpointToModel(currentModel Model, regPrivateMode *mongodbatlas.RegionalizedPrivateEndpointSetting) *Model {
	out := &Model{
		ApiKeys:   currentModel.ApiKeys,
		Enabled:   &regPrivateMode.Enabled,
		ProjectId: currentModel.ProjectId,
	}
	return out
}
