package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/webitel/cases/auth"
	customrel "github.com/webitel/custom/reflect"
	customreg "github.com/webitel/custom/registry"
	cstore "github.com/webitel/custom/store"
	custompgx "github.com/webitel/custom/store/postgres"
)

func (s *Store) Custom() cstore.Catalog {
	if s.customStore == nil {
		cs := custompgx.NewCatalog(s.conn)
		// if err != nil {
		// 	return nil
		// }
		s.customStore = cs
	}
	return s.customStore
}

const (
	// base type name to lookup extension for ...
	customTypeCases = "cases"
	customFieldName = "custom"
	customCtxState  = "__cx__"
)

func (s *Store) GetExtension(ctx context.Context, dc int64, pkg string) customrel.ExtensionDescriptor {
	custom, err := customreg.GetExtension(ctx, dc, pkg)
	if err != nil {
		slog.Warn("[custom]: extensions/cases", "error", err)
		return nil
	}
	return custom
}

func (s *Store) Extension(as customrel.ExtensionDescriptor) custompgx.ExtensionQueryBuilder {
	if as == nil {
		return nil
	}
	store := s.Custom()
	if kind, is := store.(*custompgx.Catalog); is {
		impl, err := kind.Extension(as)
		if err != nil {
			slog.Warn(fmt.Sprintf("[custom]: %s", as.Path()), "error", err)
			// not available !
			return nil
		}
		return impl
	}
	// not available !
	return nil
}

// customCtx query state
type customCtx struct {
	ok     bool                            // initialized ?
	typof  customrel.ExtensionDescriptor   // dataset typeof.Fields descriptor
	refer  custompgx.ExtensionQueryBuilder // dataset queries builder
	table  string                          // [optional] common table ; [default] refer.Table()
	fields []string                        // query custom{field(s)..} to return
	// query  sq.Sqlizer
	// params custompgx.Parameters
}

// prepare custom "extensions/cases" querier context
func (c *CaseStore) custom(ctx context.Context) (custom *customCtx) {
	// opts, is := ctx.(options.SearchOptions)
	if opts, _ := ctx.(interface {
		GetFilter(string) any
	}); opts != nil {
		// try to extract already prepared context
		custom, _ = opts.GetFilter(customCtxState).(*customCtx)
	}

	if custom == nil {
		custom = &customCtx{}
	} else if custom.typof != nil {
		if custom.typof.Dictionary().Name() != customTypeCases {
			custom = &customCtx{} // [re]init ;
		}
	}
	if !custom.ok {
		var dc int64 // zero ; invalid
		if opts, _ := ctx.(interface {
			GetAuthOpts() auth.Auther
		}); opts != nil {
			// detect customer's domain component id
			dc = opts.GetAuthOpts().GetDomainId()
		}
		if dc < 1 {
			// failed ; cannot determine domain component id !
			return nil
		}
		custom.typof = c.storage.GetExtension(
			ctx, dc, customTypeCases,
		)
		custom.refer = c.storage.Extension(
			custom.typof,
		)
		custom.ok = true // once: prepared !
	}
	return custom
}
