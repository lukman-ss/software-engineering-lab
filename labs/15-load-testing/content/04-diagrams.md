## Load Test Diagram
```text
+-------------------------------------------------------------+
|                     Load Test Runner                        |
|                                                             |
|  [VU 1] -----> HTTP POST /booking                           |
|  [VU 2] -----> HTTP POST /booking                           |
|  ...                                                        |
|  [VU N] -----> HTTP POST /booking                           |
|                                                             |
|  (Tiap VU mencatat latensi ke memory slice independen)     |
+------------------------------+------------------------------+
                               |
                               v
+-------------------------------------------------------------+
|                  HTTP Server (/booking)                     |
|                                                             |
|        +-------------------------------------------+        |
|        | Semafor Penampung Koneksi (Max = 5 Slot)  |        |
|        +-------------------------------------------+        |
|            | Slot 1 | Slot 2 | Slot 3 | Slot 4 | Slot 5     |
|                                                             |
|   (Request ke-6 dan seterusnya tertahan mengantre)          |
|   (Pemrosesan simulasi query memakan durasi 20ms)           |
+-------------------------------------------------------------+
```
