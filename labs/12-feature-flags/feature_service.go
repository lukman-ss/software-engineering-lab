package featureflags

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sync"
)

// Flag represents the configuration for a single feature flag.
type Flag struct {
	Key               string
	Enabled           bool
	RolloutPercentage int      // 0 - 100
	InternalUsers     []string // User IDs allowed for internal targeting
}

// FeatureService handles runtime evaluation and mutation of feature flags.
type FeatureService struct {
	mu    sync.RWMutex
	flags map[string]Flag
}

// NewFeatureService creates a FeatureService instance.
func NewFeatureService(flags []Flag) *FeatureService {
	fm := make(map[string]Flag)
	for _, f := range flags {
		fm[f.Key] = f
	}
	return &FeatureService{flags: fm}
}

// SetFlag dynamically updates or creates a flag at runtime without restarting.
func (s *FeatureService) SetFlag(f Flag) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags[f.Key] = f
}

// GetFlag retrieves current configuration of a flag.
func (s *FeatureService) GetFlag(key string) (Flag, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, exists := s.flags[key]
	return f, exists
}

// IsEnabled evaluates whether a feature is active for a given userID.
func (s *FeatureService) IsEnabled(featureKey, userID string) bool {
	s.mu.RLock()
	flag, exists := s.flags[featureKey]
	s.mu.RUnlock()

	// 1. Safe fallback: unknown flags are always OFF
	if !exists {
		return false
	}

	// 2. Global OFF / Kill Switch: immediate OFF for all callers
	if !flag.Enabled {
		return false
	}

	// 3. Internal Targetting: explicitly whitelisted users get the feature ON
	for _, internalUser := range flag.InternalUsers {
		if userID == internalUser {
			return true
		}
	}

	// 4. Percentage Rollout: 0% means OFF for normal users
	if flag.RolloutPercentage <= 0 {
		return false
	}
	if flag.RolloutPercentage >= 100 {
		return true
	}

	// 5. Deterministic Bucketing: hash(userID + featureKey) % 100
	// Ensures a consistent bucket (0-99) per user across requests.
	hashInput := fmt.Sprintf("%s:%s", userID, featureKey)
	hash := sha256.Sum256([]byte(hashInput))
	bucketVal := binary.BigEndian.Uint64(hash[:8])
	bucket := int(bucketVal % 100)

	// User is eligible if bucket < RolloutPercentage (e.g. 10% -> 0-9)
	return bucket < flag.RolloutPercentage
}
