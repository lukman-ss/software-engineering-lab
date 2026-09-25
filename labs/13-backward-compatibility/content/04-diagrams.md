# Diagrams

## 1. Arsitektur Expand-Migrate-Contract

Diagram di bawah ini mengilustrasikan bagaimana sistem menangani request secara backward-compatible dari Consumer Legacy maupun Modern menggunakan pipeline routing dan fallback backfill.

```text
[ Client V1 (Legacy) ]     [ Client V2 (Modern) ]
          │                         │
          ▼                         ▼
┌──────────────────────────────────────────────────┐
│              HTTP / Service Layer                │
│  - Transformation / Deprecation Pipeline         │
│  - Feature Flags (WriteMode, ReadMode)           │
│  - Observability Metrics (Legacy/New Traffic)    │
└─────────┬────────────────────────────────┬───────┘
          │                                │
          ▼                                ▼
┌──────────────────┐             ┌─────────────────┐
│ Legacy Storage   │◄──Backfill──┤ Modern Storage  │
│ (users.phone)    │   Worker    │ (user_phones)   │
└──────────────────┘             └─────────────────┘
```
