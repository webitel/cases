package util

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

func getUserNameColumn(usersTableAlias string, columnAlias string) string {
	return fmt.Sprintf("COALESCE(%s.name, %[1]s.username) %s_name", usersTableAlias, columnAlias)
}

func SetUserColumn(base squirrel.SelectBuilder, mainTableAlias string, userTableAlias string, columnAlias string) squirrel.SelectBuilder {
	base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.updated_by = %[1]s.id", userTableAlias, mainTableAlias))
	base = base.Column(Ident(userTableAlias, fmt.Sprintf("id %s_id", columnAlias)))
	base = base.Column(getUserNameColumn(userTableAlias, columnAlias))
	return base
}
