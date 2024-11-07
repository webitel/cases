package postgres

import (
	"context"
	"time"

	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	general "buf.build/gen/go/webitel/general/protocolbuffers/go"

	sq "github.com/Masterminds/squirrel"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseComment struct {
	storage store.Store
}

// Publish implements store.CommentCaseStore for publishing a single comment.
func (c *CaseComment) Publish(
	ctx context.Context,
	rpc *model.CreateOptions,
	add *_go.CaseComment,
) (*_go.CaseComment, error) {
	// Establish database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.comment.create.database_connection_error", dbErr)
	}

	// Build the insert and select query with RETURNING clause
	query, args, err := c.buildCreateCommentsQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.comment.create.query_build_error", err)
	}

	var (
		createdBy, updatedBy general.Lookup
		createdAt, updatedAt time.Time
	)

	// Dynamically build scan arguments based on the requested fields
	scanArgs := []interface{}{}
	for _, field := range rpc.Fields {
		switch field {
		case "id":
			scanArgs = append(scanArgs, &add.Id)
		case "case_id":
			scanArgs = append(scanArgs, &add.CaseId)
		case "created_at":
			scanArgs = append(scanArgs, &createdAt)
		case "comment":
			scanArgs = append(scanArgs, &add.Text)
		case "created_by":
			scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
			add.CreatedBy = &createdBy
		case "updated_by":
			scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
			add.UpdatedBy = &updatedBy
		case "updated_at":
			scanArgs = append(scanArgs, &updatedAt)
		}
	}

	// Execute the query and scan the result directly into add
	if err := d.QueryRow(ctx, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.comment.create.scan_error", err)
	}

	// Set createdAt and updatedAt fields on add after scanning
	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)

	return add, nil
}

// buildCreateCommentsQuery builds a single query that inserts one comment
// and returns only the specified fields from the inserted row using a CTE.
func (c *CaseComment) buildCreateCommentsQuery(
	rpc *model.CreateOptions,
	comment *_go.CaseComment, // Single comment
) (string, []interface{}, error) {
	// Ensure "id" is in the fields list
	rpc.Fields = util.EnsureIdField(rpc.Fields)
	// Start building the insert part of the query using Squirrel
	insertBuilder := sq.
		Insert("cases.case_comment").
		Columns("dc", "case_id", "created_at", "created_by", "updated_at", "updated_by", "comment").
		Values(
			rpc.Session.GetDomainId(), // dc
			rpc.ID,                    // case_id
			rpc.CurrentTime(),         // created_at (and updated_at)
			rpc.Session.GetUserId(),   // created_by (and updated_by)
			rpc.CurrentTime(),         // updated_at
			rpc.Session.GetUserId(),   // updated_by
			comment.Text,              // comment text
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	// Generate SQL for the INSERT with the RETURNING clause
	insertSQL, insertArgs, err := insertBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Wrap the insert in a CTE to allow selecting specific fields from the inserted rows
	insertSQL = "WITH ins AS (" + insertSQL + ")"

	// Dynamically build the SELECT query for retrieving only the specified fields
	selectBuilder := sq.Select().From("ins")

	// Add only the fields specified in rpc.Fields to the SELECT clause
	for _, field := range rpc.Fields {
		switch field {
		case "id":
			selectBuilder = selectBuilder.Column("ins.id")
		case "case_id":
			selectBuilder = selectBuilder.Column("ins.case_id")
		case "created_at":
			selectBuilder = selectBuilder.Column("ins.created_at")
		case "comment":
			selectBuilder = selectBuilder.Column("ins.comment")
		case "created_by":
			selectBuilder = selectBuilder.
				Column("ins.created_by AS created_by_id").
				Column("COALESCE(c.name::text, c.username) AS created_by_name").
				LeftJoin("directory.wbt_user c ON c.id = ins.created_by")
		case "updated_by":
			selectBuilder = selectBuilder.
				Column("ins.updated_by AS updated_by_id").
				Column("COALESCE(u.name::text, u.username) AS updated_by_name").
				LeftJoin("directory.wbt_user u ON u.id = ins.updated_by")
		case "updated_at":
			selectBuilder = selectBuilder.Column("ins.updated_at")
		}
	}

	// Generate SQL and arguments for the final combined query
	selectSQL, selectArgs, err := selectBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Combine the CTE insert statement with the select statement
	finalSQL := insertSQL + " " + selectSQL
	finalArgs := append(insertArgs, selectArgs...)

	return finalSQL, finalArgs, nil
}

// Delete implements store.CommentCaseStore.
func (c *CaseComment) Delete(
	ctx context.Context,
	req *model.DeleteOptions,
) (*_go.CaseComment, error) {
	panic("unimplemented")
}

// List implements store.CommentCaseStore.
func (c *CaseComment) List(
	ctx context.Context,
	rpc *model.SearchOptions,
) (*_go.CaseCommentList, error) {
	panic("unimplemented")
}

// Update implements store.CommentCaseStore.
func (c *CaseComment) Update(
	ctx context.Context,
	req *model.UpdateOptions,
	upd *_go.CaseComment,
) (*_go.CaseComment, error) {
	panic("unimplemented")
}

func NewCaseCommentStore(store store.Store) (store.CaseCommentStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_comment.check.bad_arguments",
			"error creating comment case interface to the case_comment table, main store is nil")
	}
	return &CaseComment{storage: store}, nil
}
