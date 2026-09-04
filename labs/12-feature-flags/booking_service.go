package featureflags

// User mendeskripsikan user yang melakukan booking
type User struct {
	ID         string
	IsInternal bool
}

// BookingRequest payload untuk request booking
type BookingRequest struct {
	UserID     string
	BranchID   string
	MechanicID string
	ServiceType string
}

// BookingResponse payload hasil booking
type BookingResponse struct {
	Flow    string `json:"flow"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// BookingService mengelola alur booking dengan safe fallback
type BookingService struct {
	featureService *FeatureService
	metrics        *Metrics
	simulateFail   bool // flag internal untuk exercise testing kill switch
}

// NewBookingService membuat BookingService baru
func NewBookingService(fs *FeatureService, metrics *Metrics) *BookingService {
	return &BookingService{featureService: fs, metrics: metrics}
}

// SetSimulateFail memungkinkan trigger failure rate secara artificial 
func (s *BookingService) SetSimulateFail(fail bool) {
	s.simulateFail = fail
}

// CreateBooking mengeksekusi booking berdasarkan evaluasi feature flag
func (s *BookingService) CreateBooking(req BookingRequest) BookingResponse {
	// Feature Flag boundary
	if s.featureService.IsEnabled("online_booking", req.UserID) {
		return s.newBookingFlow(req)
	}
	return s.legacyBookingFlow(req)
}

func (s *BookingService) newBookingFlow(req BookingRequest) BookingResponse {
	if s.simulateFail {
		return BookingResponse{
			Flow:    "online_booking",
			Message: "Vendor Error (Simulated)",
			Success: false,
		}
	}
	return BookingResponse{
		Flow:    "online_booking",
		Message: "Online booking enabled",
		Success: true,
	}
}

func (s *BookingService) legacyBookingFlow(req BookingRequest) BookingResponse {
	return BookingResponse{
		Flow:    "legacy",
		Message: "Online booking is not available",
		Success: true, // Legacy stable
	}
}
