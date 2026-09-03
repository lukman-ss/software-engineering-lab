# Git Auto Commit Skill

## Trigger
- Dipanggil otomatis setelah menyelesaikan instruksi yang mengedit atau menambah file.
- Atau saat user meminta commit.

## Workflow
1. Cek status repositori dengan `git status`.
2. Jika ada perubahan, stage semua perubahan: `git add .`
3. Buat pesan commit yang sangat singkat dan deskriptif berdasarkan perubahan (format: `feat/fix/chore: deskripsi`).
4. Eksekusi commit: `git commit -m "pesan"`
5. Laporkan ke user bahwa perubahan telah di-commit.
