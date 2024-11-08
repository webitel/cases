package postgres

import (
	"strconv"
	"strings"
	"time"

	_go "github.com/webitel/cases/api/cases"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

type CaseComment struct {
	storage store.Store
}

// Publish implements store.CommentCaseStore for publishing a single comment.
func (c *CaseComment) Publish(
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
		createdBy, updatedBy _go.Lookup
		createdAt, updatedAt time.Time
		id                   int64 // Not present in model | need for eta encoding
	)

	// Dynamically build scan arguments based on the requested fields
	scanArgs := []interface{}{}
	for _, field := range rpc.Fields {
		switch field {
		case "id":
			scanArgs = append(scanArgs, &id)
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
		case "ver":
			scanArgs = append(scanArgs, &add.Ver)
		}
	}

	// Execute the query and scan the result directly into add
	if err := d.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.comment.create.scan_error", err)
	}
	//  Encode etag from the comment ID and version
	e := etag.EncodeEtag(etag.EtagCaseComment, id, add.Ver)
	add.Id = e

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
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)
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
		case "ver":
			selectBuilder = selectBuilder.Column("ins.ver")
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
	rpc *model.DeleteOptions,
) error {
	d, err := c.storage.Database()
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.comment.delete.database_connection_error", err)
	}

	query, args, dbErr := c.buildDeleteCaseCommentQuery(rpc)
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.cases.comment.delete.query_build_error", dbErr)
	}

	res, execErr := d.Exec(rpc.Context, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("postgres.cases.comment.delete.exec_error", execErr)
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBNoRowsError("postgres.cases.comment.delete.not_found")
	}

	return nil
}

func (c CaseComment) buildDeleteCaseCommentQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deleteCaseCommentQuery
	args := []interface{}{pq.Array(ids), rpc.Session.GetDomainId()}
	return query, args, nil
}

var deleteCaseCommentQuery = store.CompactSQL(`
	DELETE FROM cases.case_comment
	WHERE id = ANY($1) AND dc = $2
`)

// List implements store.CommentCaseStore.
func (c *CaseComment) List(rpc *model.SearchOptions) (*_go.CaseCommentList, error) {
	// Connect to the database
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.list.database_connection_error", dbErr)
	}

	// Build the query using BuildListCaseCommentsQuery
	query, args, err := c.BuildListCaseCommentsQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.list.query_build_error", err)
	}

	// Execute the query
	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.list.execution_error", err)
	}
	defer rows.Close()

	var commentList []*_go.CaseComment
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1 // Fetch all records if size is -1

	for rows.Next() {
		// Respect size limit unless fetching all records
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		comment := &_go.CaseComment{}
		var (
			createdBy, updatedBy         _go.Lookup
			tempCreatedAt, tempUpdatedAt time.Time
			scanArgs                     []interface{}
		)

		// Scan fields dynamically based on requested fields
		for _, field := range rpc.Fields {
			switch field {
			case "id":
				scanArgs = append(scanArgs, &comment.Id)
			case "comment":
				scanArgs = append(scanArgs, &comment.Text)
			case "case_id":
				scanArgs = append(scanArgs, &comment.CaseId)
			case "ver":
				scanArgs = append(scanArgs, &comment.Ver)
			case "created_at":
				scanArgs = append(scanArgs, &tempCreatedAt)
			case "updated_at":
				scanArgs = append(scanArgs, &tempUpdatedAt)
			case "created_by":
				scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
			case "updated_by":
				scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
			}
		}

		// Execute the scan
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.cases.case_comment.list.row_scan_error", err)
		}

		id, err := strconv.Atoi(comment.Id)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.cases.case_comment.list.row_scan_error", err)
		}

		//  Encode etag from the comment ID and version
		e := etag.EncodeEtag(etag.EtagCaseComment, int64(id), comment.Ver)
		comment.Id = e

		// Set optional fields if present
		if util.ContainsField(rpc.Fields, "created_by") {
			comment.CreatedBy = &createdBy
		}
		if util.ContainsField(rpc.Fields, "updated_by") {
			comment.UpdatedBy = &updatedBy
		}
		if util.ContainsField(rpc.Fields, "created_at") {
			comment.CreatedAt = util.Timestamp(tempCreatedAt)
		}
		if util.ContainsField(rpc.Fields, "updated_at") {
			comment.UpdatedAt = util.Timestamp(tempUpdatedAt)
		}
		// Check if "edited" is requested and set it if created and updated times differ
		if util.ContainsField(rpc.Fields, "edited") {
			comment.Edited = !tempCreatedAt.Equal(tempUpdatedAt)
		}

		// Append to the comment list
		commentList = append(commentList, comment)
		lCount++
	}

	// Return results in CaseCommentList format
	return &_go.CaseCommentList{
		Page:  int64(rpc.Page),
		Next:  next,
		Items: commentList,
	}, nil
}

func (c *CaseComment) BuildListCaseCommentsQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	// Begin building the query with base selections
	queryBuilder := sq.Select().
		From("cases.case_comment AS cc").
		Where(sq.Eq{"cc.dc": rpc.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	if rpc.Id != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": rpc.Id})
	}

	// Ensure "id" is in the fields list
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)

	// Ensure that if "edited" is requested, both "updated_at" and "created_at" are included
	if util.ContainsField(rpc.Fields, "edited") {
		rpc.Fields = util.EnsureFields(rpc.Fields, "updated_at", "created_at")
	}

	// Add columns based on selected fields
	for _, field := range rpc.Fields {
		switch field {
		case "id", "comment", "created_at", "updated_at", "ver", "case_id":
			queryBuilder = queryBuilder.Column("cc." + field)
		case "created_by":
			queryBuilder = queryBuilder.
				Column("COALESCE(created_by.id, 0) AS cbi").    // created_by_id with NULL handled as 0
				Column("COALESCE(created_by.name, '') AS cbn"). // created_by_name with NULL handled as ''
				LeftJoin("directory.wbt_auth AS created_by ON cc.created_by = created_by.id")
		case "updated_by":
			queryBuilder = queryBuilder.
				Column("COALESCE(updated_by.id, 0) AS ubi").    // updated_by_id with NULL handled as 0
				Column("COALESCE(updated_by.name, '') AS ubn"). // updated_by_name with NULL handled as ''
				LeftJoin("directory.wbt_auth AS updated_by ON cc.updated_by = updated_by.id")
		}
	}

	// Filter by Qin
	if len(rpc.IDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.id": rpc.IDs})
	}

	// Apply filters for specific fields, like filtering by case ID or partial text match
	if caseID, ok := rpc.Filter["case_id"].(string); ok && caseID != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": caseID})
	}

	if text, ok := rpc.Filter["text"].(string); ok && len(text) > 0 {
		substr := util.Substring(text)
		combinedLike := strings.Join(substr, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"cc.text": combinedLike})
	}

	// Sort based on rpc.Sort with specified columns and sort order
	parsedFields := util.FieldsFunc(rpc.Sort, util.InlineFields)
	var sortFields []string

	for _, sortField := range parsedFields {
		desc := false
		if strings.HasPrefix(sortField, "!") {
			desc = true
			sortField = strings.TrimPrefix(sortField, "!")
		}

		column := "cc." + sortField
		if desc {
			column += " DESC"
		} else {
			column += " ASC"
		}
		sortFields = append(sortFields, column)
	}

	queryBuilder = queryBuilder.OrderBy(sortFields...)

	// Pagination: Apply offset and limit based on page number and size
	size := rpc.GetSize()
	page := rpc.Page

	if page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1))
	}

	// Generate the final SQL and arguments
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.cases.case_comments.query_build.sql_generation_error", err)
	}

	return store.CompactSQL(query), args, nil
}

// Update implements store.CommentCaseStore.
func (c *CaseComment) Update(
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
