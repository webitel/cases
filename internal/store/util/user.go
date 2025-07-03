package util

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

func getUserNameColumn(usersTableAlias string, columnAlias string) string {
	return fmt.Sprintf("COALESCE(%s.name, %[1]s.username) %s_name", usersTableAlias, columnAlias)
}

func SetUserColumn(base squirrel.SelectBuilder, from string, userTable string, columnAlias string) squirrel.SelectBuilder {
	base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s = %[1]s.id", userTable, Ident(from, columnAlias)))
	base = base.Column(Ident(userTable, fmt.Sprintf("id %s_id", columnAlias)))
	base = base.Column(getUserNameColumn(userTable, columnAlias))
	return base
}

func getContactNameColumn(contactTable string, columnAlias string) string {
	return fmt.Sprintf("COALESCE(%s.common_name, %[1]s.username) %s_name", contactTable, columnAlias)
}

func SetContactColumn(base squirrel.SelectBuilder, from string, contactTable string, columnAlias string) squirrel.SelectBuilder {
	base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %s = %[1]s.id", contactTable, Ident(from, columnAlias)))
	base = base.Column(Ident(contactTable, fmt.Sprintf("id %s_id", columnAlias)))
	base = base.Column(getContactNameColumn(contactTable, columnAlias))
	return base
}
