# Coding Style

## Go
- ใช้ `gofmt`/`go vet` มาตรฐาน, ผ่าน `golangci-lint run` เสมอ
- Exported ใช้ PascalCase, unexported ใช้ camelCase
- Acronym: uppercase ตอน exported (`HTML`, `DB`, `URL`), lowercase ตอน unexported (`db`, `url`)
- Interface ตั้งชื่อตามหน้าที่ (`Database`) หรือใช้ suffix `-er` (`Handler`)
- Import order: standard library ก่อน, เว้นบรรทัด, third-party, เว้นบรรทัด, local packages
- ห้ามเขียน query ซับซ้อนใน handler — ห้าม raw SQL ถ้าไม่จำเป็น ให้ผ่าน sqlc-generated queries เท่านั้น
- Migration ต้องมีคู่ `.up.sql` / `.down.sql` ที่ rollback ได้เสมอ
- ทุก error ต้อง check และ wrap ด้วย `fmt.Errorf("...: %w", err)` เพื่อรักษา context
- ใช้ `errors.Join` เมื่อต้องรวมหลาย error
- Log ด้วย structured logging (`slog.Logger`, key-value pairs) ก่อน return error ที่สำคัญ
- Context: ฟังก์ชันที่ต้องใช้ ให้รับ `context.Context` เป็น parameter แรกเสมอ
- เพิ่ม compile-time interface check ทันทีหลัง struct definition: `var _ InterfaceName = (*StructName)(nil)`
- Comment เอกสาร export ทุกตัวขึ้นต้นด้วยชื่อของมัน: `// Handler handles requests.` ใช้ `//` เท่านั้น ไม่ใช้ `/* */` ยกเว้น package doc

## templ / HTMX
- Component ตั้งชื่อ PascalCase, ไฟล์อยู่ใน `internal/components/<feature>/`
- แก้เฉพาะ `.templ` source ห้ามแก้ generated `.go` — รัน `go tool templ generate -path ./internal/components` หลังแก้
- ใช้ `{ variable }` สำหรับ interpolation, `{ function() }` สำหรับเรียกฟังก์ชัน, `@ComponentName(args)` สำหรับ compose component อื่น
- Tag ทุกตัวต้องปิด (`<div></div>` หรือ self-closing `<br/>`)
- Interactivity ใช้ HTMX attribute (`hx-get`, `hx-post`, `hx-target`, `hx-swap`, `hx-trigger`) แทน JS framework — เขียน vanilla JS เฉพาะเมื่อ HTMX ทำไม่ได้จริง ๆ
- แสดง loading/error/success state ทุกครั้งที่มี HTMX request (เช่น `htmx-indicator`, error partial)

## SQLC / Database
- Query annotation บังคับ: `-- name: FunctionName :one|:many|:exec|:execresult`
- Source of truth คือ `internal/db/queries/query.sql`, ห้ามแก้ generated code ใน `internal/db/queries/` โดยตรง — รัน `go tool sqlc generate`
- ใช้ `?` placeholder สำหรับ SQLite
- รัน `go tool sqlc vet` ก่อน commit เมื่อแก้ query
- แก้ schema ต้องสร้าง migration ใหม่เสมอ (`migrate create -ext sql -dir internal/db/migrations <name>`) ห้ามแก้ migration ที่ apply ไปแล้ว

## Tailwind CSS
- Utility-first, class เล็ก single-purpose (`text-center`, `bg-blue-500`, `p-4`)
- Responsive prefix ตาม breakpoint (`sm:`, `md:`, `lg:`)
- แก้ที่ `styles/input.css` เท่านั้น ห้ามแก้ generated `internal/dist/assets/css` ตรง ๆ

## General
- ตั้งชื่อไฟล์ให้สื่อความหมาย
- ลบ code ที่ไม่ได้ใช้
- หลีกเลี่ยง duplicate logic
- เขียน comment เฉพาะส่วนที่ non-obvious (hidden constraint, workaround, subtle invariant)
- ห้ามสร้าง documentation file ถ้าไม่ได้รับคำสั่ง

## Testing
- ทุก change ต้องมี test ที่เกี่ยวข้อง
- Unit test: `go test -v ./...`, race detector: `go test -race ./...`
- Single package/test: `go test -v ./internal/server/handler -run TestHome`
- E2E (Playwright, tag `e2e`): `go test -v ./e2e -tags=e2e -run TestName`
- ห้ามลบ test โดยไม่ได้รับอนุมัติ

## Review Code
ใช้ skill `code-review` หรือ `review` เป็น final review gate สำหรับ ticket GOX-XXXX

อ่าน:
- spec/context ที่อ้างใน ticket
- related files
- tests ที่เกี่ยวข้อง
- acceptance criteria

ให้ตอบเป็น:
- Verdict: PASS / NEEDS FIX / BLOCKED
- Findings เรียงตาม severity
- Missing tests ถ้ามี
- Commands ที่ควรรัน (`go test`, `golangci-lint run`, `go tool sqlc vet`, `go build`)
- ควรย้ายไป done ได้หรือยัง
