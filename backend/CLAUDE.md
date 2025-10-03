# Media CMS Web Service Architecture (Go + MVCS, Ports & Adapters)

This document specifies the architecture for a Go service that powers a low‑cost, media‑heavy community website. It follows MVCS (Entities, View, Controller, Service) and a **Ports & Adapters** layout:

- **Service = business logic + ports (interfaces) it needs**
- **Clients = adapters** that satisfy those ports using concrete tech (DB, Object Storage)
- **HTTP mappers** convert transport payloads to/from **internal entities**

---

## 1. Goals

- Serve read‑heavy public content with strong cacheability.
- Keep infra costs low while supporting many images (and future video).
- Provide a simple admin workflow for CRUD of texts, images, and timeline entries.
- Integrate with an external Events API without duplicating its data.

## 2. Storage decision

### Object Storage for images + NoSQL DB for texts and metadata (chosen)
- DB stores texts, metadata, and storage URLs/paths. Images live in object storage.
- Client fetches media directly from object storage using URLs provided by the API.

**Pros:** low DB footprint/IO, generous free tiers, flexible for large media.  

## 3. Domain entities

- **Text**: `id`, `slug`, `content`, `page_id?`, `page_slug?`, `created_at`, `updated_at`
- **TimelineEntry**: `id`, `name`, `text`, `location`, `date`, audit fields
- **Image**: `id`, `slug?`, `object_url`, `name`, `text`, `date`, `location`, audit fields

> All entities carry standard audit fields (`created_at`, `updated_at`, `last_updated_by`).

## 4. REST API (JSON, `/api/v1`)

### Texts `/texts`
- `GET /texts/:slug` • `GET /texts/id/:id` • `GET /texts`  
- `GET /texts/page/:pageId` • `GET /texts/page/slug/:pageSlug`  
- `POST /texts` • `PUT /texts/:id` • `DELETE /texts/:id`

### Images `/images`
- `GET /images/:id` • `GET /images/gallery/:slug`  
- `POST /images` (metadata + optional base64 payload)  
- `PUT /images/:id` • `DELETE /images/:id`

### Timeline `/timelineentries`
- Standard CRUD.

### Events proxy `/events`
- `GET /events?limit=N&orderBy=startDate&desc=true` • proxy to community Events API only.

## 5. Request flows

### Read flow (public site)

![Read flow](sandbox:/mnt/data/ae7d6466-23ee-4392-bd75-f9e3c69bae4c.png)

1) Client requests texts or image metadata from the REST API.  
2) API reads DB and returns JSON including object storage URLs.  
3) Client fetches media directly from the object store.  
4) CDN/browser caching applies via `Cache-Control` and ETags.

### Write flow (admins)

![Write flow](sandbox:/mnt/data/b61b1575-df59-48c0-bc07-77dd0ad22e0a.png)

1) Admin posts text/image updates.  
2) API writes image bytes to object storage, persists metadata/paths to DB.  
3) Client refetches after updates — backend is the source of truth.

## 6. MVCS with Ports & Adapters

### 6.1 Entities
Go structs in `internal/entities` (no external DTOs).

### 6.2 Service (business logic + ports)
- Package `internal/service` defines the **ports** (interfaces) needed to execute use cases and implements orchestration logic.
- Examples of rules: slug normalization, page scoping, versioned object keys, signed URL strategy, deduplication, and transactional updates across DB + object store.

**Ports (examples):**
```go
// internal/service/ports.go
type DBPort interface {
    GetTextBySlug(ctx context.Context, slug string) (entities.Text, error)
    GetTextByID(ctx context.Context, id string) (entities.Text, error)
    ListTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error)
    CreateText(ctx context.Context, t entities.Text) (entities.Text, error)
    UpdateText(ctx context.Context, id string, patch entities.Text) (entities.Text, error)
    CreateImageMeta(ctx context.Context, img entities.Image) (entities.Image, error)
    UpdateImageMeta(ctx context.Context, id string, patch entities.Image) (entities.Image, error)
}

type ObjectStorePort interface {
    PutObject(ctx context.Context, key string, data []byte) (publicURL string, err error)
    DeleteObject(ctx context.Context, key string) error
    SignedURL(ctx context.Context, key string) (string, error)
}
```

**Service orchestration (example):**
```go
// internal/service/content_service.go
type ContentService interface {
    GetTextBySlug(ctx context.Context, slug string) (entities.Text, error)
    UploadImage(ctx context.Context, meta entities.Image, data []byte) (entities.Image, error)
}

type contentService struct { db DBPort; obj ObjectStorePort }

func NewContentService(db DBPort, obj ObjectStorePort) ContentService {
    return &contentService{db: db, obj: obj}
}

func (s *contentService) UploadImage(ctx context.Context, meta entities.Image, data []byte) (entities.Image, error) {
    key := "images/" + meta.Slug // apply versioning scheme here if needed
    url, err := s.obj.PutObject(ctx, key, data)
    if err != nil { return entities.Image{}, err }
    meta.ObjectURL = url
    return s.db.CreateImageMeta(ctx, meta)
}
```

### 6.3 Clients (adapters that satisfy ports)
- Package `internal/clients` provides **constructors that return the service ports**; implementations are unexported structs that delegate to repositories/gateways and use small local helpers for DTO↔entity mapping.

```go
// internal/clients/db_client.go
var _ service.DBPort = (*dbClient)(nil)

type dbClient struct {
    text repository.TextRepo   // raw DB adapter
    img  repository.ImageRepo
}

func NewDBClient(text repository.TextRepo, img repository.ImageRepo) service.DBPort {
    return &dbClient{text: text, img: img}
}
```

```go
// internal/clients/object_client.go
var _ service.ObjectStorePort = (*objectClient)(nil)

type objectClient struct { gw gateway.ObjectGateway }

func NewObjectClient(gw gateway.ObjectGateway) service.ObjectStorePort {
    return &objectClient{gw: gw}
}
```

### 6.4 Repository (raw DB access)
- Package `internal/repository` encapsulates the database SDK and returns DB‑native DTOs. No business logic.

### 6.5 Gateway (raw object storage access)
- Package `internal/gateway` encapsulates object storage SDK (upload/delete/signed URL). No business logic.

### 6.6 HTTP: handlers + mappers (transport layer)
- Package `internal/http/handlers` holds controllers (routing, validation, error mapping).  
- Package `internal/http/mapper` holds **HTTP mappers** that convert request/response payloads ↔ internal entities. These mappers are pure and transport‑specific.

```go
// internal/http/mapper/text.go
type CreateTextRequest struct {
    Slug string `json:"slug"`; Content string `json:"content"`; PageSlug string `json:"page_slug"`
}
func ToEntity(r CreateTextRequest) entities.Text { return entities.Text{Slug: r.Slug, Content: r.Content, PageSlug: r.PageSlug} }
```

Handlers **depend only on services**:

```go
// internal/http/handlers/texts.go
func (h *TextsHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req mapper.CreateTextRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.BadRequest(w, err); return }
    ent := mapper.ToEntity(req)
    out, err := h.svc.CreateText(r.Context(), ent)
    if err != nil { httpx.FromError(w, err); return }
    httpx.JSON(w, out, http.StatusCreated)
}
```

## 7. Caching, security, ops (highlights)

- Long `Cache-Control` for immutable media objects, short TTL for listing JSON; use CDN in front of object storage.
- Admin writes behind JWT/session; validate types/size; image content‑type allow‑list.
- Structured logs with request IDs; health checks for DB and storage; simple metrics.

## 8. Folder layout

```
/cmd/server/main.go
/internal/
  /entities/                 # pure domain types
  /service/                  # business logic + ports (DBPort, ObjectStorePort)
    ports.go
    content_service.go
  /clients/                  # adapters: constructors return ports; structs are unexported
    db_client.go
    object_client.go
  /repository/               # DB adapters (DTOs + minimal logic)
    firestore/...
    postgres/...
  /gateway/                  # object store adapter
    gcs/...  s3/...
  /http/
    /handlers/               # controllers
      texts_handler.go
      images_handler.go
    /mapper/                 # HTTP <-> entities mappers
      text_mapper.go
      image_mapper.go
  /platform/                 # config, logging, error helpers, middleware
  /test/fakes/               # in-memory implementations of service ports
/docs/architecture.md
/configs/{dev,prod}.yaml
```

## 9. Compatibility notes

- Public identifiers are slugs; numeric IDs remain internal.  
- After admin writes, clients should refetch — backend is the source of truth for media metadata and URLs.  
- Deploy on Cloud Run or any container host; object storage and DB can be Firebase Storage + Firestore (or equivalents).

---

## Appendix A. Example routes
```
GET /api/v1/texts
GET /api/v1/texts/{slug}
GET /api/v1/texts/id/{id}
GET /api/v1/texts/page/{pageId}
GET /api/v1/texts/page/slug/{pageSlug}
POST /api/v1/texts
PUT  /api/v1/texts/{id}
DEL  /api/v1/texts/{id}

GET /api/v1/images/{id}
POST /api/v1/images
PUT  /api/v1/images/{id}
DEL  /api/v1/images/{id}

GET /api/v1/timelineentries
POST /api/v1/timelineentries
...

GET /api/v1/events?limit=N&orderBy=startDate&desc=true
```

# Grupy Sanca API format

one of the client is the GrupyEventsClient, which calls the grupysanca events API, below is documentation about that API

# Grupy Sanca — Events API Overview

> Base docs: `https://eventos.grupysanca.com.br/api`
> Base API URL: `https://eventos.grupysanca.com.br/api/v1`

---

## 1) API URL

* **Docs (root):** `https://eventos.grupysanca.com.br/api`
* **Base (v1):** `https://eventos.grupysanca.com.br/api/v1`

---

## 2) Common endpoints for listing events

* **List events (paged/sorted):**
  `GET /v1/events{?sort,filter,page[size],page[number]}`

* **List upcoming events (simple feed):**
  `GET /v1/events/upcoming`

* **Single event by ID or identifier:**
  `GET /v1/events/{event_identifier}`
  (The placeholder accepts a numeric ID or the event’s string identifier.)

* **Events for a specific group (optional):**
  `GET /v1/groups/{group_id}/events{?sort,filter,page[size],page[number]}`

**Examples**

```http
GET https://eventos.grupysanca.com.br/api/v1/events/upcoming
GET https://eventos.grupysanca.com.br/api/v1/events?page[size]=20&page[number]=1&sort=identifier
GET https://eventos.grupysanca.com.br/api/v1/events/12345
GET https://eventos.grupysanca.com.br/api/v1/events/b8324ae2
```

**Notes**

* Pagination uses `page[size]` (page length) and `page[number]` (1-based).
* Sorting uses `?sort=...` (prefix with `-` for descending if supported, e.g., `sort=-starts-at`).

---

## 3) Response structure (JSON:API)

The API follows **JSON:API 1.0**, so collection responses look like:

```json
{
  "meta": { "count": 1 },
  "data": [
    {
      "type": "event",
      "id": "1",
      "attributes": {
        "name": "Example Event",
        "description": "Optional blurb",
        "starts-at": "2016-12-13T23:59:59.123456+00:00",
        "ends-at": "2016-12-14T23:59:59.123456+00:00",
        "timezone": "UTC",
        "location-name": "Main Hall",
        "logo-url": "https://example.com/logo.png",
        "thumbnail-image-url": null,
        "large-image-url": null,
        "original-image-url": "https://example.com/logo.png",
        "icon-image-url": null,
        "identifier": "b8324ae2",
        "privacy": "public",
        "state": "draft",
        "created-at": "2017-06-26T15:22:37.205399+00:00"
      },
      "relationships": {
        "tickets": {
          "links": {
            "self": "/v1/events/1/relationships/tickets",
            "related": "/v1/events/1/tickets"
          }
        },
        "sessions": {
          "links": {
            "self": "/v1/events/1/relationships/sessions",
            "related": "/v1/events/1/sessions"
          }
        },
        "social-links": {
          "links": {
            "self": "/v1/events/1/relationships/social-links",
            "related": "/v1/events/1/social-links"
          }
        }
      },
      "links": { "self": "/v1/events/1" }
    }
  ],
  "jsonapi": { "version": "1.0" },
  "links": { "self": "/v1/events" }
}
```

**Headers**

* Send `Accept: application/vnd.api+json` (and `Content-Type` for non-GET).

---

## 4) Events (fields & usage tips)

**Common attributes you’ll render**

* Identity: `id` (numeric), `identifier` (string), `name`
* Timing: `starts-at`, `ends-at`, `timezone`
* Location: `location-name` (and related fields if present)
* Media: `logo-url`, `thumbnail-image-url`, `large-image-url`, `original-image-url`, `icon-image-url`
* Status/meta: `privacy`, `state`, `created-at`, `description` (when available)

**Relationships you might follow**

* `tickets`, `sessions`, `social-links`, `sponsors`, `tracks`, `microlocations`, `event-topic`, `event-sub-topic`, etc.
  Each relationship provides `links.self` and `links.related` endpoints.

**Practical listing recipe**

* If you only need the next few events: `GET /v1/events/upcoming` and render the first N.
* For archives or infinite scroll: `GET /v1/events?page[size]=20&page[number]=1&sort=-starts-at` and increment `page[number]`.

---

## Quick client snippet (fetch upcoming)

```js
const res = await fetch(
  "https://eventos.grupysanca.com.br/api/v1/events/upcoming",
  { headers: { Accept: "application/vnd.api+json" } }
);
const json = await res.json();
const events = json.data; // JSON:API list
```

---

## Error handling (typical JSON:API)

Errors generally arrive as:

```json
{
  "errors": [
    { "status": "404", "title": "Not Found", "detail": "Event not found" }
  ]
}
```

Check `response.ok` and handle `errors[]` accordingly.

---

### TL;DR

* **Base:** `https://eventos.grupysanca.com.br/api/v1`
* **List:** `/events` with `page[size]`/`page[number]` and `sort`
* **Upcoming:** `/events/upcoming`
* **Single:** `/events/{event_identifier}` (ID or identifier)
* **Format:** JSON:API with `data[]`, `attributes`, `relationships`, `links`
