package compat

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrContractViolation = errors.New("cannot apply contract: legacy traffic still active")
	ErrLegacyUnavailable = errors.New("legacy field has been retired (contracted)")
)

type Service struct {
	store    *MemoryStore
	flags    *FeatureFlags
	obs      *Observability
	backfill *BackfillWorker
}

func NewService(store *MemoryStore, flags *FeatureFlags, obs *Observability) *Service {
	if store == nil {
		store = NewMemoryStore()
	}
	if flags == nil {
		flags = NewFeatureFlags()
	}
	if obs == nil {
		obs = NewObservability()
	}
	return &Service{
		store:    store,
		flags:    flags,
		obs:      obs,
		backfill: NewBackfillWorker(store, obs, 50),
	}
}

func (s *Service) Store() *MemoryStore {
	return s.store
}

func (s *Service) Flags() *FeatureFlags {
	return s.flags
}

func (s *Service) Obs() *Observability {
	return s.obs
}

func (s *Service) BackfillWorker() *BackfillWorker {
	return s.backfill
}

// CreateUser handles writes according to current WriteMode
func (s *Service) CreateUser(name, phone string, extraPhones ...string) (UserResponse, error) {
	writeMode := s.flags.GetWriteMode()

	switch writeMode {
	case WriteLegacyOnly:
		u, err := s.store.CreateLegacy(name, phone)
		if err != nil {
			return UserResponse{}, err
		}
		var ph string
		if u.Phone != nil {
			ph = *u.Phone
		}
		return UserResponse{ID: u.ID, Name: u.Name, Phone: ph}, nil

	case WriteDual:
		u, err := s.store.CreateDual(name, phone)
		if err != nil {
			s.obs.DualWriteErrors.Add(1)
			return UserResponse{}, err
		}
		s.obs.DualWriteCount.Add(1)
		// If extra phones supplied in dual-write, also write them to user_phones
		for _, extra := range extraPhones {
			_, _ = s.store.SavePhoneEntry(u.ID, extra, false)
		}
		return s.GetUser(u.ID)

	case WriteNewOnly:
		allPhones := append([]string{phone}, extraPhones...)
		u, err := s.store.CreateModern(name, allPhones)
		if err != nil {
			return UserResponse{}, err
		}
		return s.GetUser(u.ID)

	default:
		return UserResponse{}, errors.New("unknown write mode")
	}
}

// GetUser returns an enriched response with both old and new fields
func (s *Service) GetUser(id int) (UserResponse, error) {
	u, err := s.store.GetUser(id)
	if err != nil {
		return UserResponse{}, err
	}

	phones, _ := s.store.GetPhones(id)
	var phoneVal string

	readMode := s.flags.GetReadMode()

	switch readMode {
	case ReadLegacyOnly:
		if u.Phone != nil {
			phoneVal = *u.Phone
		}
	case ReadFallback:
		if len(phones) > 0 {
			// Find primary or first phone
			phoneVal = phones[0].Number
			for _, p := range phones {
				if p.IsPrimary {
					phoneVal = p.Number
					break
				}
			}
		} else if u.Phone != nil {
			// Fallback read from legacy column
			phoneVal = *u.Phone
			// Lazy backfill: write to new schema during fallback read
			_, _ = s.store.SavePhoneEntry(id, phoneVal, true)
		}
	case ReadNewOnly:
		if len(phones) > 0 {
			phoneVal = phones[0].Number
			for _, p := range phones {
				if p.IsPrimary {
					phoneVal = p.Number
					break
				}
			}
		}
	}

	return UserResponse{
		ID:     u.ID,
		Name:   u.Name,
		Phone:  phoneVal,
		Phones: phones,
	}, nil
}

// GetLegacyUser simulates a legacy v1 client reading the endpoint
func (s *Service) GetLegacyUser(id int) (LegacyConsumerDTO, error) {
	s.obs.LegacyReadHits.Add(1)

	if s.flags.IsContractApplied() {
		return LegacyConsumerDTO{}, ErrLegacyUnavailable
	}

	resp, err := s.GetUser(id)
	if err != nil {
		return LegacyConsumerDTO{}, err
	}

	return LegacyConsumerDTO{
		ID:    resp.ID,
		Name:  resp.Name,
		Phone: resp.Phone,
	}, nil
}

// GetModernUser simulates an upgraded v2 client reading the endpoint
func (s *Service) GetModernUser(id int) (ModernConsumerDTO, error) {
	s.obs.NewReadHits.Add(1)

	resp, err := s.GetUser(id)
	if err != nil {
		return ModernConsumerDTO{}, err
	}

	return ModernConsumerDTO{
		ID:     resp.ID,
		Name:   resp.Name,
		Phones: resp.Phones,
	}, nil
}

// ReconcileData compares legacy users.phone against user_phones primary number to detect data drift
func (s *Service) ReconcileData() (int, error) {
	if s.flags.IsContractApplied() {
		return 0, nil
	}

	ids := s.store.GetUserIDs(0, s.store.TotalUsers()+10)
	drifts := 0

	for _, id := range ids {
		user, err := s.store.GetUser(id)
		if err != nil {
			continue
		}
		if user.Phone == nil {
			continue
		}

		phones, err := s.store.GetPhones(id)
		if err != nil || len(phones) == 0 {
			drifts++
			s.obs.DriftDetected.Add(1)
			continue
		}

		var primaryNumber string
		for _, p := range phones {
			if p.IsPrimary {
				primaryNumber = p.Number
				break
			}
		}

		if primaryNumber != *user.Phone {
			drifts++
			s.obs.DriftDetected.Add(1)
		}
	}

	return drifts, nil
}

// ApplyContract verifies conditions and drops legacy schema
func (s *Service) ApplyContract(force bool) error {
	// Guard: Ensure legacy read hits are zero since last reset, unless forced
	if !force && s.obs.LegacyReadHits.Load() > 0 {
		return fmt.Errorf("%w: recorded %d legacy reads", ErrContractViolation, s.obs.LegacyReadHits.Load())
	}

	// 1. Switch write mode to NewOnly
	s.flags.SetWriteMode(WriteNewOnly)
	// 2. Switch read mode to NewOnly
	s.flags.SetReadMode(ReadNewOnly)
	// 3. Mark contract applied
	s.flags.SetContractApplied(true)
	// 4. Drop legacy column from storage
	s.store.ApplyContractDropLegacyColumn()

	return nil
}

func SerializeResponse(v any) ([]byte, error) {
	return json.Marshal(v)
}
