package component

import "github.com/untemi/carshift/internal/db/sqlc"

func ProfilePicture(u *sqlc.User) string {
	if u.PfpName != "" {
		return "/pictures/pfp/" + u.PfpName
	} else {
		return "https://ui-avatars.com/api/?name=" + u.Firstname + " " + u.Lastname
	}
}
