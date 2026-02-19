package app

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/options"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/util"
)

const pageSize = 5000


// exportCSV streams CSV data page by page — each page is sent as a separate gRPC chunk.
func (c *CaseService) exportCSV(
	ctx context.Context,
	req *cases.ExportCasesRequest,
	fields []string,
	stream cases.Cases_ExportCasesServer,
) error {
	page := 1

	for {
		pageOpts, err := c.buildExportPageOptions(ctx, req, fields, page, pageSize)
		if err != nil {
			return errors.Internal(fmt.Sprintf("failed to build search options: %v", err))
		}

		list, err := c.app.Store.Case().List(pageOpts)
		if err != nil {
			return errors.Internal(fmt.Sprintf("failed to list cases: %v", err))
		}

		casesPage := list.GetItems()
		if len(casesPage) == 0 {
			break
		}

		rows, err := casesToRows(casesPage, fields)
		if err != nil {
			return errors.Internal(fmt.Sprintf("failed to convert cases to rows: %v", err))
		}

		chunkData, err := generateCSVChunk(fields, rows, page)
		if err != nil {
			return errors.Internal(fmt.Sprintf("failed to generate CSV chunk: %v", err))
		}

		if len(chunkData) > 0 {
			if err := stream.Send(&cases.ExportCasesResponse{Data: chunkData}); err != nil {
				return errors.Internal(fmt.Sprintf("failed to send chunk: %v", err))
			}
		}

		if !list.Next {
			break
		}
		page++
	}

	return nil
}

func (c *CaseService) exportXLSX(
	ctx context.Context,
	req *cases.ExportCasesRequest,
	fields []string,
	stream cases.Cases_ExportCasesServer,
) error {
	var allRows [][]string
	page := 1

	// 1. Collect all rows
	for {
		pageOpts, err := c.buildExportPageOptions(ctx, req, fields, page, pageSize)
		if err != nil {
			return errors.Internal(fmt.Sprintf("failed to build search options: %v", err))
		}

		list, err := c.app.Store.Case().List(pageOpts)
		if err != nil {
			return errors.Internal(fmt.Sprintf("failed to list cases: %v", err))
		}

		casesPage := list.GetItems()
		if len(casesPage) == 0 {
			break
		}

		rows, err := casesToRows(casesPage, fields)
		if err != nil {
			return errors.Internal(fmt.Sprintf("failed to convert cases to rows: %v", err))
		}

		allRows = append(allRows, rows...)

		if !list.Next {
			break
		}
		page++
	}

	if len(allRows) == 0 {
		return errors.InvalidArgument("no cases to export")
	}

	// 2. Generate complete XLSX using StreamWriter
	xlsxData, err := generateXLSXStreamWriter(fields, allRows)
	if err != nil {
		return errors.Internal(fmt.Sprintf("failed to generate XLSX: %v", err))
	}

	// 3. Send XLSX in chunks
	const maxChunkSize = 1024 * 1024 // 1 MB
	for i := 0; i < len(xlsxData); i += maxChunkSize {
		end := i + maxChunkSize
		if end > len(xlsxData) {
			end = len(xlsxData)
		}
		if err := stream.Send(&cases.ExportCasesResponse{Data: xlsxData[i:end]}); err != nil {
			return errors.Internal(fmt.Sprintf("failed to send XLSX chunk: %v", err))
		}
	}

	return nil
}

// buildExportPageOptions builds search options for a single page of export data.
func (c *CaseService) buildExportPageOptions(
	ctx context.Context,
	req *cases.ExportCasesRequest,
	fields []string,
	page, pageSize int,
) (*options.SearchOptions, error) {
	opts, err := options.NewSearchOptions(
		ctx,
		options.WithSearch(&cases.SearchCasesRequest{
			Page:      int32(page),
			Size:      int32(pageSize),
			Q:         req.GetQ(),
			Ids:       req.GetIds(),
			Sort:      req.GetSort(),
			Fields:    fields,
			Filters:   req.GetFilters(),
			ContactId: req.GetContactId(),
			Qin:       req.GetQin(),
			FiltersV1: req.GetFiltersV1(),
		}),
		options.WithPagination(&cases.SearchCasesRequest{
			Page: int32(page),
			Size: int32(pageSize),
		}),
		options.WithFields(
			&cases.SearchCasesRequest{Fields: fields},
			CaseMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
			util.EnsureCustomField,
		),
		options.WithFiltersV1(c.filtrationEnv, req.GetFiltersV1()),
		options.WithFilters(req.GetFilters()),
		options.WithSort(&cases.SearchCasesRequest{Sort: req.GetSort()}),
		options.WithQin(req.GetQin()),
	)
	if err != nil {
		return nil, err
	}
	opts.AddCustomContext("export_mode", true)
	return opts, nil
}

func getDefaultExportHeaders() []string {
	return []string{
		"id",
		"subject",
		"description",
		"status",
		"priority",
		"assignee",
		"created_at",
		"updated_at",
		"service",
		"source",
		"group",
		"rating",
	}
}

// casesToRows converts a list of cases to rows suitable for CSV/XLSX export
func casesToRows(casesList []*cases.Case, headers []string) ([][]string, error) {
	var rows [][]string

	for _, caseItem := range casesList {
		row := caseToRow(caseItem, headers)
		rows = append(rows, row)
	}

	return rows, nil
}

func caseToRow(caseItem *cases.Case, headers []string) []string {
	var row []string

	var customMap map[string]interface{}
	if caseItem.Custom != nil {
		customMap = caseItem.Custom.AsMap()
	}

	for _, header := range headers {
		value := getFieldValueForExport(caseItem, header, customMap)
		row = append(row, value)
	}

	return row
}

// getFieldValueForExport extracts a field value from a case and converts it to string
func getFieldValueForExport(caseItem *cases.Case, fieldName string, customMap map[string]interface{}) string {
	if caseItem == nil {
		return ""
	}

	switch fieldName {
	case "id":
		return strconv.FormatInt(caseItem.Id, 10)
	case "etag":
		return caseItem.Etag
	case "ver":
		return strconv.FormatInt(int64(caseItem.Ver), 10)
	case "subject":
		return caseItem.Subject
	case "description":
		return caseItem.Description
	case "contact_info":
		return caseItem.ContactInfo
	case "status":
		if caseItem.Status != nil {
			return caseItem.Status.Name
		}
		return ""
	case "status_id":
		if caseItem.Status != nil {
			return strconv.FormatInt(caseItem.Status.Id, 10)
		}
		return ""
	case "priority":
		if caseItem.Priority != nil {
			return caseItem.Priority.Name
		}
		return ""
	case "priority_id":
		if caseItem.Priority != nil {
			return strconv.FormatInt(caseItem.Priority.Id, 10)
		}
		return ""
	case "assignee":
		if caseItem.Assignee != nil {
			return caseItem.Assignee.Name
		}
		return ""
	case "assignee_id":
		if caseItem.Assignee != nil {
			return strconv.FormatInt(caseItem.Assignee.Id, 10)
		}
		return ""
	case "reporter":
		if caseItem.Reporter != nil {
			return caseItem.Reporter.Name
		}
		return ""
	case "reporter_id":
		if caseItem.Reporter != nil {
			return strconv.FormatInt(caseItem.Reporter.Id, 10)
		}
		return ""
	case "impacted":
		if caseItem.Impacted != nil {
			return caseItem.Impacted.Name
		}
		return ""
	case "impacted_id":
		if caseItem.Impacted != nil {
			return strconv.FormatInt(caseItem.Impacted.Id, 10)
		}
		return ""
	case "created_by":
		if caseItem.CreatedBy != nil {
			return caseItem.CreatedBy.Name
		}
		return ""
	case "created_by_id":
		if caseItem.CreatedBy != nil {
			return strconv.FormatInt(caseItem.CreatedBy.Id, 10)
		}
		return ""
	case "created_at":
		return formatTimeForExport(caseItem.CreatedAt)
	case "updated_by":
		if caseItem.UpdatedBy != nil {
			return caseItem.UpdatedBy.Name
		}
		return ""
	case "updated_by_id":
		if caseItem.UpdatedBy != nil {
			return strconv.FormatInt(caseItem.UpdatedBy.Id, 10)
		}
		return ""
	case "updated_at":
		return formatTimeForExport(caseItem.UpdatedAt)
	case "service":
		if caseItem.Service != nil {
			return caseItem.Service.Name
		}
		return ""
	case "service_id":
		if caseItem.Service != nil {
			return strconv.FormatInt(caseItem.Service.Id, 10)
		}
		return ""
	case "source":
		if caseItem.Source != nil {
			return caseItem.Source.Name
		}
		return ""
	case "source_id":
		if caseItem.Source != nil {
			return strconv.FormatInt(caseItem.Source.Id, 10)
		}
		return ""
	case "group":
		if caseItem.Group != nil {
			return caseItem.Group.Name
		}
		return ""
	case "group_id":
		if caseItem.Group != nil {
			return strconv.FormatInt(caseItem.Group.Id, 10)
		}
		return ""
	case "close_reason_group":
		if caseItem.CloseReasonGroup != nil {
			return caseItem.CloseReasonGroup.Name
		}
		return ""
	case "close_reason_group_id":
		if caseItem.CloseReasonGroup != nil {
			return strconv.FormatInt(caseItem.CloseReasonGroup.Id, 10)
		}
		return ""
	case "close_reason":
		if caseItem.CloseReason != nil {
			return caseItem.CloseReason.Name
		}
		return ""
	case "close_reason_id":
		if caseItem.CloseReason != nil {
			return strconv.FormatInt(caseItem.CloseReason.Id, 10)
		}
		return ""
	case "close_result":
		return caseItem.CloseResult
	case "rating":
		return strconv.FormatInt(caseItem.Rating, 10)
	case "rating_comment":
		return caseItem.RatingComment
	case "sla":
		if caseItem.Sla != nil {
			return caseItem.Sla.Name
		}
		return ""
	case "sla_id":
		if caseItem.Sla != nil {
			return strconv.FormatInt(caseItem.Sla.Id, 10)
		}
		return ""
	case "sla_condition":
		if caseItem.SlaCondition != nil {
			return caseItem.SlaCondition.Name
		}
		return ""
	case "sla_condition_id":
		if caseItem.SlaCondition != nil {
			return strconv.FormatInt(caseItem.SlaCondition.Id, 10)
		}
		return ""
	case "status_condition":
		if caseItem.StatusCondition != nil {
			return caseItem.StatusCondition.Name
		}
		return ""
	case "status_condition_id":
		if caseItem.StatusCondition != nil {
			return strconv.FormatInt(caseItem.StatusCondition.Id, 10)
		}
		return ""
	case "planned_reaction_at":
		return formatTimeForExport(caseItem.PlannedReactionAt)
	case "planned_resolve_at":
		return formatTimeForExport(caseItem.PlannedResolveAt)
	case "reacted_at":
		return formatTimeForExport(caseItem.ReactedAt)
	case "resolved_at":
		return formatTimeForExport(caseItem.ResolvedAt)
	case "difference_in_reaction":
		return strconv.FormatInt(caseItem.DifferenceInReaction, 10)
	case "difference_in_resolve":
		return strconv.FormatInt(caseItem.DifferenceInResolve, 10)
	case "author":
		if caseItem.Author != nil {
			return caseItem.Author.Name
		}
		return ""
	case "author_id":
		if caseItem.Author != nil {
			return strconv.FormatInt(caseItem.Author.Id, 10)
		}
		return ""
	case "custom":
		if len(customMap) > 0 {
			data, err := json.Marshal(customMap)
			if err != nil {
				return ""
			}
			return string(data)
		}
		return ""
	default:
		if v, ok := customMap[fieldName]; ok {
			return formatCustomFieldValue(v)
		}
		return ""
	}
}

// Objects with "name" are simplified to just the name.
// Arrays of objects are joined with ", ".
func formatCustomFieldValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case map[string]interface{}:
		if name, ok := val["name"]; ok {
			return fmt.Sprintf("%v", name)
		}
		data, _ := json.Marshal(val)
		return string(data)
	case []interface{}:
		var names []string
		for _, item := range val {
			if m, ok := item.(map[string]interface{}); ok {
				if name, ok := m["name"]; ok {
					names = append(names, fmt.Sprintf("%v", name))
					continue
				}
			}
			names = append(names, fmt.Sprintf("%v", item))
		}
		return strings.Join(names, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatTimeForExport formats a timestamp for export
func formatTimeForExport(ts int64) string {
	if ts == 0 {
		return ""
	}

	t := time.Unix(ts/1000, (ts%1000)*1000000)
	return t.Format("2006-01-02 15:04:05")
}

func generateCSVChunk(headers []string, rows [][]string, page int) ([]byte, error) {
	buf := &bytes.Buffer{}

	// Write UTF-8 BOM on the first page so that Excel correctly recognizes the encoding
	if page == 1 {
		buf.Write([]byte{0xEF, 0xBB, 0xBF})
	}

	writer := csv.NewWriter(buf)

	// Only write headers on the first page
	if page == 1 {
		if err := writer.Write(headers); err != nil {
			return nil, fmt.Errorf("failed to write CSV headers: %w", err)
		}
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

func generateXLSXStreamWriter(headers []string, rows [][]string) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := f.GetSheetName(0)
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream writer: %w", err)
	}

	for i := range headers {
		if err := sw.SetColWidth(i+1, i+1, 25); err != nil {
			return nil, fmt.Errorf("failed to set column width: %w", err)
		}
	}

	currentRow := 1

	// Write header row with bold style
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create style: %w", err)
	}

	headerRow := make([]any, len(headers))
	for i, h := range headers {
		headerRow[i] = excelize.Cell{StyleID: style, Value: h}
	}
	headerCell, _ := excelize.CoordinatesToCellName(1, currentRow)
	if err := sw.SetRow(headerCell, headerRow); err != nil {
		return nil, fmt.Errorf("failed to write header row: %w", err)
	}
	currentRow++

	// Write data rows
	for rowIdx, row := range rows {
		rowData := make([]any, len(row))
		for i, v := range row {
			rowData[i] = v
		}
		cell, err := excelize.CoordinatesToCellName(1, currentRow+rowIdx)
		if err != nil {
			return nil, fmt.Errorf("failed to get cell name: %w", err)
		}
		if err := sw.SetRow(cell, rowData); err != nil {
			return nil, fmt.Errorf("failed to write row %d: %w", rowIdx, err)
		}
	}

	if err := sw.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush stream writer: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("failed to write XLSX: %w", err)
	}

	return buf.Bytes(), nil
}
