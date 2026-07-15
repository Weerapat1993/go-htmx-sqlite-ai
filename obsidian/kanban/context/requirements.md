# Website Requirements: go-htmx-sqlite-ai

## Product-Level Source Of Truth

โปรเจกต์นี้ไม่มี PRD แยกไฟล์ — ใช้ `CLAUDE.md` และ `AGENTS.md` ที่ root เป็นแหล่งอ้างอิงภาพรวมของ architecture, wiring flow, security defaults และ command reference

Agent AI ต้องอ่าน `CLAUDE.md` และ `AGENTS.md` ก่อนทำงานเหล่านี้:

- สร้าง ticket ใหม่หรือแตก requirement เป็นงานย่อย
- วางแผน feature หรือแก้ routing/handler/middleware
- แก้ database schema หรือ query
- review งานที่เกี่ยวกับ full-stack flow (handler → templ → HTMX)

ถ้า requirement ใหม่ขัดกับ `CLAUDE.md`/`AGENTS.md` ให้แจ้ง conflict และถามผู้ใช้ว่าจะอัปเดตเอกสารหรือให้ถือเป็น exception เฉพาะงานนั้น

## Goal
เว็บแอป Go server-rendered ด้วย HTMX + SQLite เป็น boilerplate/demo ที่มี Todo list, health check, และ static pages เป็นตัวอย่าง feature

## Tech Stack
- Go 1.26+ (stdlib `net/http` ServeMux, ไม่มี web framework)
- templ (type-safe HTML templating, codegen จาก `.templ`)
- HTMX (interactivity ฝั่ง client แบบ server-driven, ไม่มี SPA framework)
- Tailwind CSS v4 (CLI-based, generate จาก `styles/input.css`)
- SQLite (ผ่าน `db.sqlite3`, migration ด้วย `migrate.sh`)
- sqlc (generate type-safe Go จาก `internal/db/queries/query.sql`)
- Air (live reload สำหรับ dev)
- golangci-lint (lint), govulncheck (vulnerability scan)
- Playwright (e2e tests, Go build tag `e2e`)

## Functional Requirements
- Server-rendered pages ผ่าน templ components ใน `internal/components/`
- Interactivity ฝั่ง client ใช้ HTMX attributes (`hx-get`, `hx-post`, `hx-swap` ฯลฯ) แทน client-side JS framework
- รองรับ CRUD หลักของระบบ (เช่น Todo) ผ่าน sqlc-generated queries
- มี validation ทั้งฝั่ง handler (ก่อนเขียน DB) และแสดง error state กลับผ่าน templ partial
- แสดง loading, error และ success state ทุกครั้งที่มี HTMX interaction (เช่น `hx-indicator`, swap partial ที่มี error message)
- `GET /api/health` endpoint สำหรับ health check

## Backend Requirements
- Routes ลงทะเบียนใน `internal/server/router/router.go` ด้วย stdlib `http.ServeMux` pattern แบบ method-prefixed (เช่น `"GET /{$}"`)
- Handler เป็น struct-based (`handler.New(logger, database)`) DI ด้วย logger + `db.Database`, อยู่ใน `internal/server/handler/`
- Handler ต้องบาง ไม่ใส่ business logic เยอะ — ถ้า logic ซับซ้อนให้แยกเป็น package/service แยก
- ใช้ sqlc-generated `*queries.Queries` เท่านั้นสำหรับ DB access ห้าม raw SQL ปนใน handler
- Schema/migration อยู่ใน `internal/db/migrations`, query source of truth คือ `internal/db/queries/query.sql` (annotation `-- name: X :one/:many/:exec`)
- ใช้ `migrate.sh` สำหรับ apply migration ทุกครั้งที่แก้ schema — server ไม่ migrate เอง on startup
- Middleware chain (`Recovery` → `Logging` → `Security` → `RateLimit` → `CSRF`) ต้องคงลำดับเดิมเสมอเมื่อแก้ routing
- CSRF ใช้ Go 1.25+ native `http.CrossOriginProtection` (header-based, ไม่มี token)

## Frontend Requirements
- Templates อยู่ใน `internal/components/` แยกตาม feature (เช่น `home/`, `todo/`, `about/`, `layout/`, `core/`)
- ใช้ templ syntax เขียน component แบบ type-safe, ไม่เขียน raw HTML string ใน Go
- ใช้ HTMX attributes สำหรับ partial update แทน full page reload เมื่อทำได้
- ห้าม hardcode path —ใช้ constant/route helper ที่มีอยู่ในโปรเจกต์แทน string literal ซ้ำ ๆ
- Styling ใช้ Tailwind CSS utility classes โดยตรงใน `.templ`, ห้ามแก้ generated CSS (`internal/dist/assets/css`) ให้แก้ `styles/input.css`
- หลีกเลี่ยง component/templ file ที่ยาวเกินไป — แยก partial ย่อยเมื่อ logic/markup ซับซ้อน
- ออกแบบ UI ต้องคำนึงถึง Responsive Design
- เขียนโค้ดให้สั้นและเข้าใจง่าย (Clean Code)
- ห้ามแก้ generated `.go` ใน `internal/components/**` โดยตรง — แก้ที่ `.templ` source แล้วรัน `go tool templ generate`

## API Response Format
Endpoint ที่ตอบ JSON (เช่น `/api/health`) ใช้รูปแบบ:
```json
{
  "status": 200,
  "message": "Success",
  "data": {}
}
```
Endpoint ที่ตอบ HTML/HTMX partial ให้ตอบ templ-rendered fragment ตรง ๆ ไม่ wrap เป็น JSON
