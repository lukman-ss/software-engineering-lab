package compat

import (
	"sync/atomic"
)

type WriteMode int
type ReadMode int

const (
	WriteLegacyOnly WriteMode = iota
	WriteDual
	WriteNewOnly
)

const (
	ReadLegacyOnly ReadMode = iota
	ReadFallback
	ReadNewOnly
)

type MigrationPhase struct {
	Write WriteMode
	Read  ReadMode
}

type FeatureFlags struct {
	writeMode       atomic.Value // WriteMode
	readMode        atomic.Value // ReadMode
	contractApplied atomic.Bool  // true if legacy column/endpoints dropped
}

func NewFeatureFlags() *FeatureFlags {
	f := &FeatureFlags{}
	f.writeMode.Store(WriteLegacyOnly)
	f.readMode.Store(ReadLegacyOnly)
	return f
}

func (f *FeatureFlags) GetWriteMode() WriteMode {
	return f.writeMode.Load().(WriteMode)
}

func (f *FeatureFlags) SetWriteMode(m WriteMode) {
	f.writeMode.Store(m)
}

func (f *FeatureFlags) GetReadMode() ReadMode {
	return f.readMode.Load().(ReadMode)
}

func (f *FeatureFlags) SetReadMode(m ReadMode) {
	f.readMode.Store(m)
}

func (f *FeatureFlags) IsContractApplied() bool {
	return f.contractApplied.Load()
}

func (f *FeatureFlags) SetContractApplied(b bool) {
	f.contractApplied.Store(b)
}
