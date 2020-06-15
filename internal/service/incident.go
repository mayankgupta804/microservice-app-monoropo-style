package service

import (
	"github.com/squadcast_assignment/internal/model/domain"
	"github.com/squadcast_assignment/internal/repository"
	"github.com/squadcast_assignment/internal/serializer"
)

type IncidentService interface {
	GetIncident(id int64) (*domain.Incident, *domain.Error)
	DeleteIncident(id int64) *domain.Error
	CreateIncident(reqData serializer.CreateIncidentRequest) (int64, *domain.Error)
	UpdateIncident(id int64, reqData serializer.UpdateIncidentRequest) *domain.Error
}

type incidentServiceImpl struct {
	ir repository.IncidentRepository
}

func NewIncidentService(repo repository.IncidentRepository) incidentServiceImpl {
	is := incidentServiceImpl{ir: repo}
	return is
}

func (is incidentServiceImpl) CreateIncident(reqData serializer.CreateIncidentRequest) (int64, *domain.Error) {
	incident := domain.Incident{
		Message: reqData.Message,
	}
	lastInsertID, creationErr := is.ir.CreateIncident(incident)
	if creationErr != nil {
		return -1, domain.NewError(creationErr.Error())
	}
	if lastInsertID <= 0 {
		return -1, domain.NewError("error inserting data in db")
	}
	return lastInsertID, nil
}

func (is incidentServiceImpl) GetIncident(id int64) (*domain.Incident, *domain.Error) {
	incident, getErr := is.ir.GetIncident(id)
	if getErr != nil {
		return nil, domain.NewError(getErr.Error())
	}
	return &incident, nil
}

func (is incidentServiceImpl) UpdateIncident(id int64, reqData serializer.UpdateIncidentRequest) *domain.Error {
	incident := domain.Incident{
		Message: reqData.Message,
		Status:  reqData.ResolutionStatus,
		Ack:     reqData.Ack,
		Comment: []string{reqData.Comment},
	}

	_, updateErr := is.ir.UpdateIncident(id, incident)
	if updateErr != nil {
		return domain.NewError(updateErr.Error())
	}
	return nil
}

func (is incidentServiceImpl) DeleteIncident(id int64) *domain.Error {
	deleteErr := is.ir.DeleteIncident(id)
	if deleteErr != nil {
		return domain.NewError(deleteErr.Error())
	}
	return nil
}
