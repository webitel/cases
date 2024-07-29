package lookup

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	_gen "buf.build/gen/go/webitel/general/protocolbuffers/go"
	db "github.com/webitel/cases/internal/db"
	"github.com/webitel/cases/model"
	"log"
	"time"
)

type StatusLookup struct {
	storage db.DB
}

func (s StatusLookup) Create(rpc *model.CreateOptions, add *_go.StatusLookup) (*_go.StatusLookup, error) {
	query, args, err := s.buildCreateGroupQuery(rpc.Session.GetDomainId(), rpc.Session.GetUserId(), rpc.Time, add)
	d, dbErr := s.storage.Database()

	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	var createdByLookup, updatedByLookup _gen.Lookup

	err = d.QueryRowContext(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &rpc.Time, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&rpc.Time, &updatedByLookup.Id, &updatedByLookup.Name,
	)

	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	//When we create a new lookup - CREATED/UPDATED_AT are the same
	t := rpc.Time.Unix()

	return &_go.StatusLookup{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		CreatedAt:   t,
		UpdatedAt:   t,
		CreatedBy:   &createdByLookup,
		UpdatedBy:   &updatedByLookup,
	}, nil
}

func (s StatusLookup) Search(rpc *model.SearchOptions, ids []string) ([]*_go.StatusLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) Delete(rpc *model.DeleteOptions, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) Update(rpc *model.UpdateOptions, lookup *_go.StatusLookup) (*_go.StatusLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) buildCreateGroupQuery(domainID int64, createdBy int64, t time.Time, lookup *_go.StatusLookup) (string,
	[]interface{}, error) {
	query := `
with ins as (
    INSERT INTO contacts.group (name, dc, created_at, description, created_by, updated_at, 
updated_by) //TODO CREATE TABLE FOR CASES
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    returning *
)
select ins.id,
    ins.name,
    ins.created_at,
    ins.description,
    ins.created_by created_by_id,
    coalesce(c.name::text, c.username) created_by_name,
    ins.updated_at,
    ins.updated_by updated_by_id,
    coalesce(u.name::text, u.username) updated_by_name
from ins
  left join directory.wbt_user u on u.id = ins.updated_by
  left join directory.wbt_user c on c.id = ins.created_by;
`
	args := []interface{}{lookup.Name, domainID, t, lookup.Description, createdBy, t, createdBy}
	return query, args, nil
}

func NewStatusLookupStore(store db.DB) (db.StatusLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_status_lookup.check.bad_arguments",
			"error creating config interface to the status_lookup table, main store is nil")
	}
	return &StatusLookup{storage: store}, nil
}
