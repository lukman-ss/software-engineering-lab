package featureflags

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// Flag mendeskripsikan model feature flag
type Flag struct {
	Key               string
	Enabled           bool
	RolloutPercentage int      // 0 - 100
	InternalUsers     []string // Daftar UserID internal
}

// FeatureService mengevaluasi feature flags
type FeatureService struct {
	flags map[string]Flag
}

// NewFeatureService membuat instance FeatureService baru
func NewFeatureService(flags []Flag) *FeatureService {
	fm := make(map[string]Flag)
	for _, f := range flags {
		fm[f.Key] = f
	}
	return &FeatureService{flags: fm}
}

// UpdateFlag mengubah konfigurasi flag secara runtime
func (s *FeatureService) UpdateFlag(f Flag) {
	s.flags[f.Key] = f
}

// IsEnabled mengevaluasi apakah fitur ON untuk userID tertentu
func (s *FeatureService) IsEnabled(featureKey, userID string) bool {
	flag, exists := s.flags[featureKey]
	if !exists {
		return false // Default OFF jika flag tidak ditemukan (safe fallback)
	}

	if !flag.Enabled {
		return false // Global OFF / Kill Switch
	}

	// 1. Cek Internal Rollout
	for _, internalUser := range flag.InternalUsers {
		if userID == internalUser {
			return true
		}
	}

	// 2. Cek Percentage Rollout (Deterministic bucketing)
	// Hash (user_id + feature_key) % 100
	hashInput := fmt.Sprintf("%s:%s", userID, featureKey)
	hash := sha256.Sum256([]byte(hashInput))
	
	// Ambil 8 byte pertama untuk dijadikan integer
	bucketVal := binary.BigEndian.Uint64(hash[:8])
	bucket := int(bucketVal % 100)

	// Jika bucket di bawah persentase rollout, fitur ON. (Contoh: 10% = bucket 0-9)
	return bucket < flag.RolloutPercentage
}
