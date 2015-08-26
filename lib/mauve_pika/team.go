package mauve_pika

import (
	"strings"

	"github.com/synful/kudos/lib/purple_unicorn"
)

// This file handles the Kudos notion of a Team. The key use-case envisioned
// here is for projects that allow for multiple explicit collaborators. We
// store per-assignment team-related information in the database.
// TODO: this only hands the storing and processing of team grades. Enforcing
// aspects of Teams (their size, when they have to be registered by) are not
// handled by this package/file

type Team struct {
	Members []purple_unicorn.User
}

type TeamManifest struct {
	//assignment -> team
	Teams map[string]Team
}

func (t *Team) Name() string {
	strs := make([]string, 0, len(t.Members))
	for _, u := range t.Members {
		strs = append(strs, string(u))
	}
	return strings.Join(strs, ":")
}
