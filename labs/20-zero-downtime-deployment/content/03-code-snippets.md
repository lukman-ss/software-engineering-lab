# Code Snippets

## Snippet 1 — Graceful Server Shutdown with PreStop Delay

Source File: `internal/server/server.go`
Purpose: Mengatur penghentian server HTTP secara bertahap, mencakup perubahan status kesiapan, penundaan simulasi routing table, dan penutupan listener secara halus.

```go
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Server received shutdown request")

	s.SetReady(false)
	log.Println("Server marked unready, detached from load balancer")

	if s.preStop > 0 {
		log.Printf("Executing preStop sleep for %v to allow routing table updates...", s.preStop)
		time.Sleep(s.preStop)
	}

	log.Println("Initiating graceful shutdown of HTTP listeners...")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	s.wg.Wait()
	log.Println("All in-flight requests completed. Server stopped gracefully.")
	return nil
}
```

Explanation:
Method `Shutdown` pertama-tama menandai instance tidak siap (`s.SetReady(false)`) agar dikeluarkan dari rotasi load balancer. Lalu ia melakukan jeda eksekusi selama `preStop` untuk memberikan waktu pada jaringan mendeteksi pencabutan pod/instance, sebelum akhirnya memanggil `srv.Shutdown(ctx)` dan menunggu seluruh goroutine in-flight selesai via `s.wg.Wait()`.

---

## Snippet 2 — Expand and Contract Database Fallback Reading

Source File: `internal/db/db.go`
Purpose: Memberikan pembacaan kompatibel (backward compatibility) terhadap entitas user yang disimpan dalam format skema lama maupun baru.

```go
func (s *UserStore) GetUser(id string) (UserRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.records[id]
	if !ok {
		return UserRecord{}, ErrNotFound
	}

	if rec.FirstName == "" && rec.LastName == "" && rec.Name != "" {
		parts := strings.SplitN(rec.Name, " ", 2)
		rec.FirstName = parts[0]
		if len(parts) > 1 {
			rec.LastName = parts[1]
		}
	} else if rec.Name == "" {
		rec.Name = strings.TrimSpace(rec.FirstName + " " + rec.LastName)
	}

	return rec, nil
}
```

Explanation:
Fungsi `GetUser` memeriksa keberadaan nilai kolom. Bila kolom baru (`FirstName`, `LastName`) kosong namun kolom warisan (`Name`) ada, kode akan memecah string `Name` secara on-the-fly. Sebaliknya, bila kolom baru ada namun format lama kosong, fungsi menggabungkannya. Ini menjaga kompatibilitas aplikasi versi lama maupun baru saat beroperasi bersamaan.

---

## Snippet 3 — Worker Loop with Context Termination

Source File: `internal/worker/worker.go`
Purpose: Mengelola siklus hidup goroutine background worker yang memproses pekerjaan dan berhenti saat menerima instruksi shutdown.

```go
func (w *Worker) Start(concurrency int) {
	for i := 0; i < concurrency; i++ {
		w.wg.Add(1)
		go func(workerID int) {
			defer w.wg.Done()
			for {
				select {
				case <-w.ctx.Done():
					return
				case job, ok := <-w.jobChan:
					if !ok {
						return
					}
					log.Printf("Worker %d starting job %s", workerID, job.ID)
					time.Sleep(job.Duration)
					w.completedMu.Lock()
					w.completed = append(w.completed, job.ID)
					w.completedMu.Unlock()
					log.Printf("Worker %d finished job %s", workerID, job.ID)
				}
			}
		}(i)
	}
}
```

Explanation:
Goroutine pekerja mendengarkan channel tugas `w.jobChan` atau pembatalan konteks `w.ctx.Done()`. Selama sebuah tugas sedang diproses (`time.Sleep`), goroutine tidak langsung terputus secara mendadak saat sinyal berhenti dikirim, melainkan menyelesaikan tugas tersebut hingga tuntas sebelum membaca instruksi berhenti pada iterasi berikutnya.

---

## Snippet 4 — Simulating Orchestrator Termination in Demo

Source File: `cmd/demo/main.go`
Purpose: Menunjukkan rangkaian orkestrasi: inisialisasi, status siap, trafik in-flight, dan shutdown tertib berbasis sinyal OS.

```go
	// Simulate an orchestrator sending SIGTERM during a deployment
	time.Sleep(500 * time.Millisecond)
	log.Println("Simulating SIGTERM from orchestrator (e.g. Kubernetes)")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// In a real app we wait for <-sig, here we mock sending it to ourselves
	go func() {
		sig <- syscall.SIGTERM
	}()

	<-sig
	log.Println("SIGTERM received, initiating graceful shutdown procedures")

	// Shutdown server gracefully (includes preStop hook and connection draining)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Stop background workers gracefully
	w.Stop()
```

Explanation:
Aplikasi demo merekam sinyal sistem operasi (`SIGINT`, `SIGTERM`), menyimulasikan penerimaan `SIGTERM`, kemudian secara berurutan memicu penutupan server (dengan batas waktu timeout 10 detik via context) dan penutupan worker antrean.
