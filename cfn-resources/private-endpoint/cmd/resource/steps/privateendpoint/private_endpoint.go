package privateendpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/constants"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progressevent"
	"go.mongodb.org/atlas/mongodbatlas"
	"net/http"
	"reflect"
)

const (
	ProviderName            = "AWS"
	StatusPendingAcceptance = "PENDING_ACCEPTANCE"
	StatusPending           = "PENDING"
	StatusAvailable         = "AVAILABLE"
	StatusInitiating        = "INITIATING"
	StatusDeleting          = "DELETING"
	StatusUpdateCompleted   = "UPDATE_COMPLETED"
)

// Todo: im not convinced about this resource, maybe we can find another way
type privateEndpointCreationCallBackContext struct {
	StateName        constants.EventStatus
	ID               string
	PrivateEndpoints []AtlasPrivateEndpointCallBack
}

type AtlasPrivateEndpointCallBack struct {
	VpcId               string
	SubnetId            string
	InterfaceEndpointId string
	Status              string
}

type AtlasPrivateEndpointInput struct {
	VpcId               string
	SubnetId            string
	InterfaceEndpointId string
	Status              *string
}

func (s *privateEndpointCreationCallBackContext) FillStruct(m map[string]interface{}) error {
	s.ID = fmt.Sprint(m["ID"])
	eventStatusParam := fmt.Sprint(m["StateName"])
	eventStatus, err := constants.ParseEventStatus(eventStatusParam)
	if err != nil {
		return err
	}
	s.StateName = eventStatus

	privateEndpointI := m["PrivateEndpoints"]
	a := reflect.ValueOf(privateEndpointI)
	if a.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]AtlasPrivateEndpointCallBack, a.Len())

	for i := 0; i < a.Len(); i++ {
		v := reflect.ValueOf(a.Index(i).Interface().(interface{}))
		peCallback := AtlasPrivateEndpointCallBack{}

		//TODO: this can be done with Reflection

		for _, key := range v.MapKeys() {
			valStr := fmt.Sprint(v.MapIndex(key).Interface())
			switch key.String() {
			case "VpcId":
				peCallback.VpcId = valStr
			case "SubnetId":
				peCallback.SubnetId = valStr
			case "InterfaceEndpointId":
				peCallback.InterfaceEndpointId = valStr
			case "Status":
				peCallback.Status = valStr
			}
		}
		ret[i] = peCallback
	}

	s.PrivateEndpoints = ret

	return nil
}

func GetCallback(privateEndpointInput []AtlasPrivateEndpointInput, endpointServiceID string, state constants.EventStatus) (map[string]interface{}, error) {
	endpointCallBacks := make([]AtlasPrivateEndpointCallBack, len(privateEndpointInput))

	for i, pe := range privateEndpointInput {
		callBack := AtlasPrivateEndpointCallBack{
			VpcId:               pe.VpcId,
			SubnetId:            pe.SubnetId,
			InterfaceEndpointId: pe.InterfaceEndpointId,
		}

		if pe.Status != nil {
			callBack.Status = *pe.Status
		}

		endpointCallBacks[i] = callBack
	}

	callBackContext := privateEndpointCreationCallBackContext{
		StateName:        state,
		ID:               endpointServiceID,
		PrivateEndpoints: endpointCallBacks,
	}

	return getMapFromCallBackContext(callBackContext)
}

func getMapFromCallBackContext(callBackContext privateEndpointCreationCallBackContext) (map[string]interface{}, error) {
	var callBackMap map[string]interface{}
	data, _ := json.Marshal(callBackContext)
	err := json.Unmarshal(data, &callBackMap)

	return callBackMap, err
}

func Create(mongodbClient *mongodbatlas.Client, groupID string, privateEndpointInput []AtlasPrivateEndpointInput, endpointServiceID string) handler.ProgressEvent {
	for _, endpoint := range privateEndpointInput {
		interfaceEndpointRequest := &mongodbatlas.InterfaceEndpointConnection{
			ID: endpoint.InterfaceEndpointId,
		}

		_, response, err := mongodbClient.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(),
			groupID,
			ProviderName,
			endpointServiceID,
			interfaceEndpointRequest)
		if err != nil {
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
				response.Response)
		}
	}

	for i := range privateEndpointInput {
		status := StatusInitiating
		privateEndpointInput[i].Status = &status
	}

	callBackMap, err := GetCallback(privateEndpointInput, endpointServiceID, constants.CreatingPrivateEndpoint)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error Unmarshalling callback map : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest)
	}

	return progress_events.GetInProgressProgressEvent("Adding private endpoint", callBackMap, nil, 20)
}

func ValidateCreationCompletion(mongodbClient *mongodbatlas.Client, groupID string, req handler.Request) (*ValidationResponse, *handler.ProgressEvent) {
	callBackContext := privateEndpointCreationCallBackContext{}

	err := callBackContext.FillStruct(req.CallbackContext)
	if err != nil {
		pe := progress_events.GetFailedEventByCode(fmt.Sprintf("Error parsing PrivateEndpointCallBackContext : %s", err.Error()), cloudformation.HandlerErrorCodeServiceInternalError)
		return nil, &pe
	}

	completed := true
	ids := make([]string, len(callBackContext.PrivateEndpoints))

	for i := range callBackContext.PrivateEndpoints {
		ids[i] = callBackContext.PrivateEndpoints[i].InterfaceEndpointId
		if callBackContext.PrivateEndpoints[i].Status != StatusAvailable {
			privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(),
				groupID,
				ProviderName,
				callBackContext.ID,
				callBackContext.PrivateEndpoints[i].InterfaceEndpointId)
			if err != nil {
				pe := progress_events.GetFailedEventByResponse(fmt.Sprintf("Error validating private endpoint create : %s", err.Error()),
					response.Response)
				return nil, &pe
			}

			callBackContext.PrivateEndpoints[i].Status = privateEndpointResponse.AWSConnectionStatus

			switch privateEndpointResponse.AWSConnectionStatus {
			case StatusPendingAcceptance, StatusPending:
				completed = false
			case StatusAvailable:
				continue
			default:
				pe := progress_events.GetFailedEventByCode(fmt.Sprintf("Resource is in status : %s", privateEndpointResponse.AWSConnectionStatus),
					cloudformation.HandlerErrorCodeInternalFailure)
				return nil, &pe
			}
		}
	}

	if completed {
		endpoints := make([]AtlasPrivateEndpointCallBack, len(callBackContext.PrivateEndpoints))
		for i, v := range callBackContext.PrivateEndpoints {
			endpoints[i] = AtlasPrivateEndpointCallBack{
				VpcId:               v.VpcId,
				SubnetId:            v.SubnetId,
				InterfaceEndpointId: v.InterfaceEndpointId,
				Status:              v.Status,
			}
		}
		vr := ValidationResponse{
			ID:        callBackContext.ID,
			Endpoints: endpoints,
		}
		return &vr, nil
	} else {
		pe := progress_events.GetInProgressProgressEvent("Adding private endpoint in progress",
			req.CallbackContext, nil, 20)
		return nil, &pe
	}
}

func Update(mongodbClient *mongodbatlas.Client, groupID string, endpointServiceID string, toAdd []AtlasPrivateEndpointInput, toDelete []AtlasPrivateEndpointInput) handler.ProgressEvent {
	if len(toAdd) == 0 && len(toDelete) == 0 {
		return handler.ProgressEvent{
			OperationStatus: handler.Success,
			Message:         "no private endpoints to delete or create to update",
		}
	}

	if len(toAdd) > 0 {
		addProgressEvent := Create(mongodbClient, groupID, toAdd, endpointServiceID)
		if addProgressEvent.OperationStatus == handler.Failed {
			return addProgressEvent
		}

		for i := range toAdd {
			status := StatusInitiating
			toAdd[i].Status = &status
		}
	}

	if len(toDelete) > 0 {
		deleteInput := make([]string, 0, len(toDelete))

		for i := range toDelete {
			status := StatusDeleting
			toDelete[i].Status = &status
			deleteInput = append(deleteInput, toDelete[i].InterfaceEndpointId)
		}

		deleteProgressEvent := Delete(mongodbClient, groupID, endpointServiceID, deleteInput)
		if deleteProgressEvent.OperationStatus == handler.Failed {
			return *deleteProgressEvent
		}
	}

	privateEndpointInput := append(toDelete, toAdd...)
	callBackMap, err := GetCallback(privateEndpointInput, endpointServiceID, constants.UpdatingPrivateEndpoint)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error Unmarshalling callback map : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest)
	}

	return progress_events.GetInProgressProgressEvent("Adding private endpoint", callBackMap, nil, 20)
}

func ValidateUpdateCompletion(mongodbClient *mongodbatlas.Client, groupID string, req handler.Request, endpointServiceID string) (*ValidationResponse, *handler.ProgressEvent) {
	callBackContext := privateEndpointCreationCallBackContext{}

	err := callBackContext.FillStruct(req.CallbackContext)
	if err != nil {
		pe := progress_events.GetFailedEventByCode(fmt.Sprintf("Error parsing Filling Struct of Callback Context : %s", err.Error()), cloudformation.HandlerErrorCodeServiceInternalError)
		return nil, &pe
	}

	completed := true

	for i := range callBackContext.PrivateEndpoints {
		if callBackContext.PrivateEndpoints[i].Status != StatusUpdateCompleted {
			privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(),
				groupID,
				ProviderName,
				callBackContext.ID,
				callBackContext.PrivateEndpoints[i].InterfaceEndpointId)

			if callBackContext.PrivateEndpoints[i].Status == StatusDeleting {
				if err != nil {
					if response.StatusCode == http.StatusNotFound {
						callBackContext.PrivateEndpoints = remove(callBackContext.PrivateEndpoints, i)
						continue
					}
					pe := progress_events.GetFailedEventByResponse(fmt.Sprintf("Error validating private endpoint create : %s", err.Error()),
						response.Response)
					return nil, &pe
				}
				completed = false
			}

			switch privateEndpointResponse.AWSConnectionStatus {
			case StatusPendingAcceptance, StatusPending:
				completed = false
			case StatusAvailable:
				callBackContext.PrivateEndpoints[i].Status = StatusUpdateCompleted
				continue
			default:
				pe := progress_events.GetFailedEventByCode(fmt.Sprintf("Resource is in status : %s", privateEndpointResponse.AWSConnectionStatus),
					cloudformation.HandlerErrorCodeInternalFailure)
				return nil, &pe
			}
		}
	}

	if completed {
		endpoints := make([]AtlasPrivateEndpointCallBack, len(callBackContext.PrivateEndpoints))
		for i, v := range callBackContext.PrivateEndpoints {
			endpoints[i] = AtlasPrivateEndpointCallBack{
				VpcId:               v.VpcId,
				SubnetId:            v.SubnetId,
				InterfaceEndpointId: v.InterfaceEndpointId,
				Status:              v.Status,
			}
		}
		vr := ValidationResponse{
			ID:        callBackContext.ID,
			Endpoints: endpoints,
		}
		return &vr, nil
	} else {
		mapCallback, err := getMapFromCallBackContext(callBackContext)
		if err != nil {
			pe := progress_events.GetFailedEventByCode(fmt.Sprintf("Error parsing Updating PrivateEndpointCallBackContext : %s", err.Error()), cloudformation.HandlerErrorCodeServiceInternalError)
			return nil, &pe
		}
		pe := progress_events.GetInProgressProgressEvent("Update in progress",
			mapCallback, nil, 20)
		return nil, &pe
	}
}

func remove(slice []AtlasPrivateEndpointCallBack, s int) []AtlasPrivateEndpointCallBack {
	return append(slice[:s], slice[s+1:]...)
}

func sliceDifference(current, previous []AtlasPrivateEndpointInput) []AtlasPrivateEndpointInput {
	pSlice := make(map[string]string, len(previous))
	for _, p := range previous {
		pSlice[p.ToString()] = ""
	}
	var diff []AtlasPrivateEndpointInput
	for _, x := range current {
		if _, found := pSlice[x.ToString()]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func (i AtlasPrivateEndpointInput) ToString() string {
	return fmt.Sprintf("%s%s", i.VpcId, i.SubnetId)
}

func Delete(mongodbClient *mongodbatlas.Client, groupID string, endpointServiceID string, interfaceEndpoints []string) *handler.ProgressEvent {
	for _, intEndpoints := range interfaceEndpoints {
		response, err := mongodbClient.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(),
			groupID,
			ProviderName,
			endpointServiceID,
			intEndpoints)
		if err != nil {
			pe := progress_events.GetFailedEventByResponse(fmt.Sprintf("Error deleting private endpoint : %s",
				err.Error()),
				response.Response)
			return &pe
		}
	}

	return nil
}

type ValidationResponse struct {
	ID        string
	Endpoints []AtlasPrivateEndpointCallBack
}
