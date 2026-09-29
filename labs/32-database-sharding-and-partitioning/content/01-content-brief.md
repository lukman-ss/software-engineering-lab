# Content Brief

Topic: Database Sharding and Table Partitioning
Target Reader: Senior Backend Engineers, Distributed Systems Engineers, Database Architects
Problem: Skalabilitas database monolitik terbentur batasan I/O, kapasitas memori, dan write contention pada single instance, sementara transisi ke horizontal scaling memperkenalkan kompleksitas routing, write hotspots, query scatter-gather, collision ID terdistribusi, dan resharding data migration.
Core Mental Model: Logical Partitioning membagi tabel secara internal dalam 1 instance fisik (optimal untuk lifecycle & pruning), sedangkan Physical Sharding membagi dataset ke beberapa independent database nodes dengan router hashing/key-ranges; Consistent Hashing meminimalkan migrasi data saat cluster scale-out ($O(1/N)$ vs hash modulo $\approx (N-1)/N$).
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts:
- Logical Table Partitioning vs Physical Database Sharding
- Sharding Key Cardinality & Monotonic Write Hotspots
- Routing Algorithms: Hash Modulo ($K \pmod M$) vs Consistent Hash Ring dengan Virtual Nodes
- Non-Sharded Query Optimization: Scatter-Gather (Broadcast) vs Global Secondary Index (Lookup Vindex)
- Distributed Unique ID Generation: RFC 9562 UUIDv7 (Time-Ordered) & Sequence Block Allocation (Vitess-style)
- Cross-Shard Complexity: Two-Phase Commit (TwoPC) latency & fractured read isolation trade-offs
Verified Behaviors:
- Range-based logical partitioning berhasil melakukan query pruning (hanya scan 1 dari 4 partisi untuk range tertentu).
- Drop partition instan membuang data tanpa overhead VACUUM table besar.
- Sharding key berbasis nilai monotonik (date/autoincrement) menyebabkan 100% write hotspot pada 1 shard aktif.
- High-cardinality shard key (`user_id` / `tenant_id`) mendistribusikan write secara seimbang ke seluruh node.
- Scale-out node dari 4 ke 5 node pada Consistent Hash hanya memindahkan $\approx 12.00\% - 16.00\%$ keys, sementara Hash Modulo memindahkan $\approx 75.12\% - 79.84\%$ keys.
- Query berdasarkan non-shard key (`email`) memicu broadcast 4/4 shard pada scatter-gather, sedangkan GSI routing melakukan 1 direct point-lookup ke shard yang tepat.
- UUIDv7 terbukti time-ordered secara leksikografis ($u_1 < u_2$) dan sequence block allocator menghasilkan ID unik sekuensial tanpa tabrakan.
Available Case Studies: E-commerce multi-tenant transaction scaling, time-series data lifecycle management, user account lookup by email vs tenant ID.
Warnings:
- Consistent Hashing membutuhkan virtual nodes (misal 50–150 vnodes) untuk mencegah skew/ketimpangan data antar node.
- Global Secondary Index (Lookup Vindex) menambahkan overhead double-write / latency ekstra pada operasi Insert/Update.
- Two-Phase Commit (TwoPC) pada arsitektur sharded menjamin atomisitas tetapi tidak menyediakan full ACID isolation tradisional (dapat terjadi fractured reads pada cross-shard read).
- In-memory lab implementation mensimulasikan mekanisme routing dan tidak menyertakan disk persistence atau 2PC cross-shard engine.
