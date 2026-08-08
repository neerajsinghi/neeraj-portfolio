package auth

import (
	"testing"

	"neeraj-portfolio/backend/internal/blog"
)

func TestClaimRolesIgnoresUnknownGroups(t *testing.T) {
	roles := claimRoles([]any{"editor", "billing", "admin"})
	if !roles[blog.RoleEditor] || !roles[blog.RoleAdmin] {
		t.Fatalf("roles = %#v", roles)
	}
	if len(roles) != 2 {
		t.Fatalf("unexpected roles = %#v", roles)
	}
}
