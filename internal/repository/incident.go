package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/squadcast_assignment/internal/infrastructure/database"
	"github.com/squadcast_assignment/internal/model/domain"
	"github.com/squadcast_assignment/internal/model/entity"
)

// IncidentRepository exposes functionalities related to the application
type IncidentRepository interface {
	CreateIncident(incident domain.Incident) (int64, error)
	UpdateIncident(id int64, incident domain.Incident) (int64, error)
	GetIncident(id int64) (domain.Incident, error)
	DeleteIncident(id int64) error
}

type incidentRepository struct {
	db database.DBClient
}

// InitIncidentRepository returns an instance of a repository
func InitIncidentRepository(db database.DBClient) *incidentRepository {
	ir := incidentRepository{db: db}
	return &ir
}

func (ir incidentRepository) CreateIncident(incident domain.Incident) (int64, error) {
	lastInsertID, err := ir.db.Execute(
		fmt.Sprintf(`INSERT INTO incidents(message, status, ack) VALUES ("%s", "unresolved", "no");`, incident.Message), "CREATE")
	if err != nil {
		return lastInsertID, err
	}
	return lastInsertID, nil
}

func (ir incidentRepository) GetIncident(id int64) (domain.Incident, error) {
	var incidentEntity = entity.Incidents{}
	result, err := ir.db.Query(fmt.Sprintf(`SELECT id, message, status, ack FROM incidents WHERE id=%d;`, id))
	if err != nil {
		return domain.Incident{}, err
	}
	if result.Next() {
		result.Scan(&incidentEntity.ID, &incidentEntity.Message, &incidentEntity.Status, &incidentEntity.Ack)
	}
	result, err = ir.db.Query(fmt.Sprintf(`SELECT comment FROM comments WHERE incident_id=%d;`, id))
	if err != nil {
		return domain.Incident{}, err
	}

	var comment sql.NullString
	for result.Next() {
		result.Scan(&comment)
		incidentEntity.Comments = append(incidentEntity.Comments, comment)
	}
	incident := domain.Incident{
		ID:      strconv.Itoa(int(incidentEntity.ID.Int64)),
		Message: incidentEntity.Message.String,
		Status:  incidentEntity.Status.String,
		Ack:     incidentEntity.Ack.String,
	}
	for _, v := range incidentEntity.Comments {
		incident.Comment = append(incident.Comment, v.String)
	}
	return incident, nil
}

func (ir incidentRepository) UpdateIncident(id int64, incident domain.Incident) (int64, error) {
	var rowsAffected int64
	var err error

	if incident.Message != "" {
		rowsAffected, err = ir.db.Execute(
			fmt.Sprintf(`UPDATE incidents SET message=("%s") WHERE id=%d;`, incident.Message, id), "UPDATE")
		if err != nil {
			return rowsAffected, err
		}
	}

	if incident.Ack != "" {
		rowsAffected, err = ir.db.Execute(
			fmt.Sprintf(`UPDATE incidents SET ack=("%s") WHERE id=%d;`, incident.Ack, id), "UPDATE")
		if err != nil {
			return rowsAffected, err
		}
	}

	if incident.Status != "" {
		rowsAffected, err = ir.db.Execute(
			fmt.Sprintf(`UPDATE incidents SET status=("%s") WHERE id=%d;`, incident.Status, id), "UPDATE")
		if err != nil {
			return rowsAffected, err
		}
	}

	if incident.Comment != nil {
		if len(incident.Comment) == 0 {
			return -1, errors.New("comment is empty")
		}
		rowsAffected, err = ir.db.Execute(
			fmt.Sprintf(`INSERT INTO comments(incident_id, comment) VALUES(%d, "%s");`, id, incident.Comment[0]), "CREATE")
		if err != nil {
			return rowsAffected, err
		}
	}

	return rowsAffected, nil
}

func (ir incidentRepository) DeleteIncident(id int64) error {
	_, err := ir.db.Execute(fmt.Sprintf("DELETE FROM incidents WHERE id=%d;", id), "DELETE")
	if err != nil {
		return err
	}
	return nil
}
