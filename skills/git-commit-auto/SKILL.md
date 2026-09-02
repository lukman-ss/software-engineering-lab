---
name: git-commit-auto
description: Otomatis commit perubahan git setelah mengedit/membuat file atau ketika diminta commit. Gunakan saat pengguna ingin commit perubahan kode secara manual atau otomatis.
---

# Git Commit Auto

Skill untuk melakukan git commit setelah ada perubahan file.

## Kapan Digunakan

- Diminta commit perubahan git.
- Setelah mengedit atau membuat file baru dalam repositori git.

## Langkah Eksekusi

1. Cek status git:
   ```bash
   git status --porcelain
   ```
2. Jika ada perubahan (staged/unstaged/untracked):
   ```bash
   git add -A
   git commit -m "Auto-commit: update repository files"
   ```
3. Jika tidak ada perubahan, lewati commit.

## Catatan
- Gunakan rtk jika tersedia: `rtk git status`, `rtk git add -A`, `rtk git commit -m "..."`.
- Pastikan tidak membuat commit kosong.
