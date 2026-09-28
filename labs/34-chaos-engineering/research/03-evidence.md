## Evidence 1

Claim: Chaos Engineering adalah disiplin eksperimentasi pada sistem untuk membangun kepercayaan pada kemampuan sistem menahan kondisi turbulen di production.

Evidence: "Chaos Engineering is the discipline of experimenting on a system in order to build confidence in the system’s capability to withstand turbulent conditions in production."

Source: Principles of Chaos Engineering

URL: https://principlesofchaos.org/

Confidence: HIGH

Corroborated By: AWS Well-Architected Reliability Pillar

Notes: Definisi diterima secara luas di industri sebagai standar keandalan sistem terdistribusi.

## Evidence 2

Claim: Penentuan Steady State dilakukan dengan mengamati metrik output perilaku sistem (misal throughput, error rates, latensi) daripada kondisi internal komponen.

Evidence: "Focus on the systemic output of the system rather than internal conditions... Map the measurable output to a standard behavior that defines normal operation: steady state."

Source: Principles of Chaos Engineering

URL: https://principlesofchaos.org/

Confidence: HIGH

Corroborated By: Netflix TechBlog

Notes: Mengukur output sistem mencegah alarm palsu akibat fluktuasi internal yang tidak berdampak pada pengguna.

## Evidence 3

Claim: Hipotesis eksperimen chaos harus dibangun berdasarkan korelasi antara gangguan yang disuntikkan (fault injection) dan metrik steady state yang tetap terjaga.

Evidence: "Propose a hypothesis: 'We expect the steady state to continue despite the injection of X fault.'"

Source: Principles of Chaos Engineering

URL: https://principlesofchaos.org/

Confidence: HIGH

Corroborated By: AWS Well-Architected Reliability Pillar

Notes: Menguji asumsi arsitektur sebelum insiden nyata terjadi.

## Evidence 4

Claim: Blast radius harus dikendalikan dengan memulai eksperimen dari lingkungan kecil (canary/staging) ke subset trafik production, serta memiliki mekanisme abort otomatis.

Evidence: "Contain the blast radius by starting small, using canary analysis, and having an automatic abort mechanism if safety thresholds are breached."

Source: AWS Well-Architected Reliability Pillar

URL: https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/chaos-engineering.html

Confidence: HIGH

Corroborated By: Gremlin Reliability Documentation

Notes: Mencegah eksperimen merusak seluruh pengguna aktif saat uji coba gagal.
