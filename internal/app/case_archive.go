package app

import (
	"context"
	"slices"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
)

const PermissionManageArchivedCases = "manage_archived_cases"

var (
	ErrCaseArchived = errors.Forbidden("case is archived", errors.WithID("app.case.archived"))

	ErrArchivedCaseStatusOnly = errors.Forbidden(
		"only the status of an archived case can be changed",
		errors.WithID("app.case.archived.status_only"),
	)

	archivedCaseUpdateFields = []string{
		"etag", "id", "ver", "userID", "status", "status_condition", "close_reason", "close_result", "close_article",
	}
)

type archiveGuardSkipKey struct{}

func withoutArchiveGuard(ctx context.Context) context.Context {
	return context.WithValue(ctx, archiveGuardSkipKey{}, true)
}

func archiveGuardSkipped(ctx context.Context) bool {
	skip, ok := ctx.Value(archiveGuardSkipKey{}).(bool)

	return ok && skip
}

func (a *App) ensureCaseNotArchived(ctx context.Context, domainID int64, caseIDs ...int64) error {
	for _, caseID := range caseIDs {
		archived, err := a.Store.Case().IsArchived(ctx, domainID, caseID)
		if err != nil {
			return err
		}

		if archived {
			return ErrCaseArchived
		}
	}

	return nil
}

func (a *App) ensureCommentsNotArchived(ctx context.Context, domainID int64, commentIDs []int64) error {
	archived, err := a.Store.CaseComment().IsArchived(ctx, domainID, commentIDs)
	if err != nil {
		return err
	}

	if archived {
		return ErrCaseArchived
	}

	return nil
}

func (a *App) ensureRelationsNotArchived(ctx context.Context, domainID int64, relationIDs []int64) error {
	archived, err := a.Store.RelatedCase().IsArchived(ctx, domainID, relationIDs)
	if err != nil {
		return err
	}

	if archived {
		return ErrCaseArchived
	}

	return nil
}

func (a *App) ensureArchivedCaseUpdateAllowed(ctx context.Context, session auth.Auther, caseID int64, mask []string) error {
	if archiveGuardSkipped(ctx) {
		return nil
	}

	archived, err := a.Store.Case().IsArchived(ctx, session.GetDomainId(), caseID)
	if err != nil {
		return err
	}

	if !archived {
		return nil
	}

	return archivedCaseUpdateAllowed(mask, session.HasPermission(PermissionManageArchivedCases))
}

func archivedCaseUpdateAllowed(mask []string, canManage bool) error {
	if !canManage {
		return ErrCaseArchived
	}

	for _, field := range mask {
		if !slices.Contains(archivedCaseUpdateFields, field) {
			return ErrArchivedCaseStatusOnly
		}
	}

	return nil
}
