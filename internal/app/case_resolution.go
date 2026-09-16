package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/webitel/cases/internal/model"
)

// casesExchange is the broker exchange of every cases event.
const casesExchange = "cases"

// Close article event for the knowledge base.
const (
	caseResolutionEventType = "case.resolution_set"
	caseResolutionSchema    = 1
	caseResolutionEvent     = "resolution_set"
	caseResolutionTimeout   = 5 * time.Second
)

// publishCaseResolution sends the close article event of a written case; a failure is only logged.
func (a *App) publishCaseResolution(ctx context.Context, domainID, caseID, articleID int64) {
	body, err := json.Marshal(model.CaseResolutionAMQPMessage{
		Type:       caseResolutionEventType,
		Schema:     caseResolutionSchema,
		OccurredAt: time.Now().UTC(),
		DomainID:   domainID,
		CaseID:     caseID,
		ArticleID:  articleID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "could not encode case resolution event", slog.Any("error", err))

		return
	}

	pubCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), caseResolutionTimeout)
	defer cancel()

	routingKey := fmt.Sprintf("%s.%s.%s.%d", casesExchange, model.ScopeCases, caseResolutionEvent, domainID)

	err = a.rabbitPublisher.Publish(pubCtx, casesExchange, routingKey, body, nil)
	if err != nil {
		slog.ErrorContext(ctx, "could not publish case resolution event",
			slog.Int64("case_id", caseID), slog.Any("error", err))
	}
}
