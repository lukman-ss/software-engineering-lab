package api_versioning

import "encoding/json"

// Invoice adalah model internal domain untuk invoicing sistem bengkel CMMS.
// Ini adalah model single source of truth yang tidak terhubung langsung dengan HTTP contract.
type Invoice struct {
	ID       int      // Invoice ID
	Customer Customer // Data customer (relasi)
	Total    int64    // Total tagihan dalam IDR
	Status   string   // PAID, PENDING, OVERDUE
}

// Customer merepresentasikan data pelanggan di sistem.
type Customer struct {
	ID    int    // Customer ID
	Name  string // Nama lengkap
	Phone string // Nomor telepon
}

// LegacyInvoice adalah model yang dipakai oleh client Android lama (versi 1.0).
// Client lama mencoba meng-unmarshal JSON langsung ke struct ini.
// Perubahan tipe `customer` dari string ke object adalah BREAKING CHANGE.
type LegacyInvoice struct {
	ID       int    `json:"id"`
	Customer string `json:"customer"` // Harus string, bukan object!
	Total    int64  `json:"total"`
	Status   string `json:"status"`
}

// ParseLegacyInvoice mensimulasikan cara legacy client mendecode respons server.
// Gunakan helper ini di test agar perilaku consumer lama benar-benar terlihat.
func ParseLegacyInvoice(body []byte) (LegacyInvoice, error) {
	var inv LegacyInvoice
	err := json.Unmarshal(body, &inv)
	return inv, err
}

// InvoiceRepository mendefinisikan interface untuk mengakses data invoice.
type InvoiceRepository interface {
	GetInvoiceByID(ctx any, id int) (Invoice, error)
}
