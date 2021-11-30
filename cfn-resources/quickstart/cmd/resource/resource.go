package resource

import (
	"context"
	//"encoding/base64"
	"errors"
	//"fmt"
	log "github.com/sirupsen/logrus"
	//"strings"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"go.mongodb.org/atlas/mongodbatlas"
	
)
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Create complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	//svc := secretsmanager.New(req.Session)

	/*
	1 - fetch AWS Secret for Atlas API Key
	2 - validate API KEY
	3 - create or ensure Project
	4 - create Cluster 
	5 - create IPAccessList
	6 - create AWS IAM Role 
	7 - create Atlas MOngoDB User with new IAM Role

	*/

	secretName := *currentModel.AtlasApiKeySecret
	log.Infof("Read for Cluster Id/SecretName:%s", secretName)
	key, err := util.GetApiKeyFromDeploymentSecret(&req, secretName)
	if err != nil {
		log.Infof("Read - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}
	log.Debugf("key:%+v", key)

	client, err := util.CreateMongoDBClient(key.PublicKey, key.PrivateKey)
	if err != nil {
		log.Infof("Read - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
	}

	project, _, err := client.Projects.Create(context.Background(), &mongodbatlas.Project{
		Name:  *currentModel.ClusterName,
		//OrgID: *key.OrgId,
	})
	if err != nil {
		//return handler.ProgressEvent{}, fmt.Errorf("error creating project: %s", err)
		log.Debugf("Create - error: %+v", err)
		// TODO- Should detect and return HandlerErrorCodeAlreadyExists
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource Not Found",
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}

	// setup currentModel here dude!!

	log.Debugf("Create - no error>>> woequicktyoeewwwww: %+v", project)
	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: Create")
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Read complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: Read")
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{}, errors.New("Not implemented: Update is not supported. Please visit Atlas to edit your Quickstart deployment.")
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Delete complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: Delete")
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "List complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: List")
}
