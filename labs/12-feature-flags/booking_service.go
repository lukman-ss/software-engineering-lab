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
}

// BookingService mengelola alur booking dengan safe fallback
type BookingService struct {
	featureService *FeatureService
}

// NewBookingService membuat BookingService baru
func NewBookingService(fs *FeatureService) *BookingService {
	return &BookingService{featureService: fs}
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
	return BookingResponse{
		Flow:    "online_booking",
		Message: "Online booking enabled",
	}
}

func (s *BookingService) legacyBookingFlow(req BookingRequest) BookingResponse {
	return BookingResponse{
		Flow:    "legacy",
		Message: "Online booking is not available",
	}
}
