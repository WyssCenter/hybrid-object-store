package config

import (
	"crypto/sha1"
	"fmt"
)

const (
	DuplexSyncType  = "duplex"
	SimplexSyncType = "simplex"
)

//================ Sync Configuration ======================
// SyncConfiguration describes a sync configuration between two namespaces, as understood by a single core service
type SyncConfiguration struct {
	SyncType string `json:"sync_type"`

	SourceCoreService string            `json:"source_core_service"`
	SourceNamespace   string            `json:"source_namespace"`
	SourcePolicies    map[string]string `json:"source_policies"`

	TargetCoreService string `json:"target_core_service"`
	TargetNamespace   string `json:"target_namespace"`
}

// Hash is a unique string for a SyncConfiguration and can be used as a key or way to compare SyncConfigurations
func (sc *SyncConfiguration) Hash() string {
	h := sha1.New()
	h.Write([]byte(sc.SourceCoreService))
	h.Write([]byte(sc.SourceNamespace))
	// ??? do we need to sort this?
	for key, val := range sc.SourcePolicies {
		h.Write([]byte(key))
		h.Write([]byte(val))
	}
	h.Write([]byte(sc.TargetCoreService))
	h.Write([]byte(sc.TargetNamespace))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash
}

func (sc *SyncConfiguration) String() string {
	dir := "->"
	if sc.SyncType == DuplexSyncType {
		dir = "<->"
	}

	return fmt.Sprintf(
		"<SyncConfiguration: %s:%s %s %s:%s>",
		sc.SourceCoreService,
		sc.SourceNamespace,
		dir,
		sc.TargetCoreService,
		sc.TargetNamespace,
	)
}
