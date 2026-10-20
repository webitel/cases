package app

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/webitel/cases/api/cases"
	ftspb "github.com/webitel/cases/api/fts"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/util"
)

const (
	ftsExportPageSize = 1000
	ftsExportMaxPages = 50
)

// ftsExportObjectNames lists the FTS object scopes searched during export.
var ftsExportObjectNames = []string{"cases"}

func (c *CaseService) resolveFtsIdsForExport(ctx context.Context, req *cases.ExportCasesRequest) (bool, error) {
	ftsFilters, rest := util.PartitionFilter(req.GetFilters(), "fts")
	if len(ftsFilters) == 0 {
		return false, nil
	}

	req.Filters = rest

	if len(req.GetIds()) > 0 {
		return true, nil
	}

	query := strings.TrimSpace(ftsFilters[0].Value)
	if query == "" {
		req.Ids = nil
		return true, nil
	}

	outCtx := ctx
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		outCtx = metadata.NewOutgoingContext(ctx, md)
	}

	var (
		ids        []string
		seen       = make(map[string]struct{})
		page int32 = 1
	)
	for ; page <= ftsExportMaxPages; page++ {
		resp, err := c.app.ftsSearchClient.Search(outCtx, &ftspb.SearchRequest{
			Q:          query,
			ObjectName: ftsExportObjectNames,
			Size:       ftsExportPageSize,
			Page:       page,
		})
		if err != nil {
			return true, errors.New("fts search failed for export", errors.WithCause(err))
		}
		for _, item := range resp.GetItems() {
			id := item.GetId()
			if id == "" {
				continue
			}
			if _, dup := seen[id]; dup {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
		if !resp.GetNext() {
			break
		}
	}

	req.Ids = ids

	return true, nil
}
