package caching

import "fmt"

// DashboardKey mengembangkan key versioning untuk handle schema migration.
// Key format: cmms:dashboard:v1:branch:{id}:date:{YYYY-MM-DD}
type DashboardKey struct {
	namespace string
	entity    string
	version   int
	date      string
	branchID  int64
}

// NewDashboardKey membuat key dengan versioning support
func NewDashboardKey(branchID int64) *DashboardKey {
	return &DashboardKey{
		namespace: "cmms",
		entity:    "dashboard",
		version:   1,
		date:      Today(),
		branchID:  branchID,
	}
}

// String mengembalikan cache key string
func (k *DashboardKey) String() string {
	return fmt.Sprintf("%s:%s:v%d:branch:%d:%s",
		k.namespace, k.entity, k.version, k.branchID, k.date)
}

// WithVersion mengubah versi (untuk migration)
func (k *DashboardKey) WithVersion(v int) *DashboardKey {
	nk := *k
	nk.version = v
	return &nk
}

// InvalidateKey untuk delete (menggunakan key yang sama)
func (k *DashboardKey) InvalidateKey() string {
	return k.String()
}