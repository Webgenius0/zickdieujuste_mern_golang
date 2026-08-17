# Software Requirements Specification — ZICK

**Stack:** Go (Echo framework) · PostgreSQL · you will follow my current code structure
**Version:** 2.0 (corrected against Figma + original SRS/ERD review)

---

## 1. Project Overview

**Purpose:** A faith-based mobile application delivering daily spiritual guidance — prayers, worship audio, motivational video, proverbs, and age-segmented (Kids/Teens) content — with scheduled push reminders and a premium subscription tier.

**Objectives:**
- Serve categorized, mixed-media content (audio, video, text) via a paginated, searchable, filterable API.
- Let users set personalized daily prayer reminder times with reliable, timezone-aware push delivery.
- Gate premium content behind subscription status validated against Apple/Google receipts.

**Target Users:** Individuals following a daily Christian/spiritual routine, including parents seeking age-appropriate content for children and teens.

---

## 2. User Roles

| Role | Access |
|---|---|
| **Free User** | Non-premium content, schedule management, profile settings |
| **Premium User** | All content including `is_premium = true` items, offline sync |
| **Admin** (inferred) | Content CMS, user management, support tooling |

---

## 3. Functional Requirements

### 3.1 Authentication & Onboarding
- Email/password registration with bcrypt hashing.
- Google/Apple OAuth (`auth_provider` on `users`; `password_hash` nullable for OAuth accounts).
- Forgot Password → 5-digit OTP sent to email → verify → set new password. **(Corrected: UI shows 5 OTP digit boxes, not 4.)**
- T&C acceptance recorded via `terms_accepted_at`.
- JWT access token + refresh token pair issued on login; refresh token persisted (hashed) server-side for revocation on logout.

### 3.2 Home Dashboard
- Personalized greeting, today's prayer schedule summary (morning/night), today's `DAILY_QUOTE` content item, quick-access grid to all content categories.

### 3.3 Content Library & Consumption
- Seven content types: `PRAYER`, `MOTIVATION`, `WORSHIP`, `PROVERB`, `DAILY_QUOTE`, `ILLUSTRATION`, `ENCOURAGEMENT`.
- Sub-tabs per type: Prayer (Morning/Night), Worship (Day/Night), Kids/Teens (Prayers/Faith).
- Kids and Teens sections can share the *same underlying content row* tagged for both audiences (`content_audiences` join) — not duplicated data.
- Illustrations and Encouragement items carry **both** a video (`media_url`) and full body text (`body_text`) — detail screens show a player above scrollable text.
- Proverbs and Daily Quotes are text-only, date-published, with a reverse-chronological "Previous" list (queried by `published_at DESC`, no separate join table needed).
- Motivation detail screens show "Related Motivation" via `related_content`.
- Free users requesting `is_premium = true` content receive `403 Forbidden` with an upgrade-prompt payload instead of the media URL.
- Library search bar hits full-text search (`search_vector` GIN index) across title + body.

### 3.4 Custom Scheduling
- User sets `morning_prayer_time`, `night_prayer_time`, and a global `push_enabled` toggle, stored with an IANA `timezone`.
- **Open item flagged back to design:** the Schedule screen's top "Your Schedule" summary (05:30 / 21:00) doesn't match the individual reminder rows below it (09:30 AM / 08:30 PM) in the mockup — confirm which is the source of truth before backend defaults are set.
- A worker evaluates due schedules against `device_tokens` and dispatches FCM/APNS pushes, respecting each user's stored `timezone`.

### 3.5 Subscriptions
- Three plans (Biannual, Annual, Friends & Family) sourced from an admin-editable `subscription_plans` table rather than hardcoded — prices are UI-visible and expected to change.
- Purchase flow: client buys via Apple/Google IAP → sends receipt to backend → backend validates against the originating `store` → updates `subscriptions` + `users.is_premium`.
- Store webhooks handle renewals, cancellations, and billing failures asynchronously.

### 3.6 Account Management
- Profile: name, email, location, avatar, theme (`IVORY`/`NAVY`), language.
- Change password (old/new/confirm).
- Delete account: soft-delete (`deleted_at`) with a grace period, then a scheduled purge job hard-deletes PII — satisfies both the "permanent, cannot be undone" UI copy and GDPR right-to-erasure without an irreversible synchronous delete.
- Log out: revokes the active refresh token.

---

## 4. Authentication & Authorization

- **JWT access tokens** (short-lived, ~15 min) + **refresh tokens** (long-lived, persisted hashed in `refresh_tokens`, revocable).
- Middleware: `RequireAuth` (valid access token) and `RequirePremium` (checks `users.is_premium` or active `subscriptions` row) gate content endpoints.
- OTP: 5-digit numeric code, expires in 5–10 minutes, single-use (`is_used`), rate-limited per email/IP.
- Passwords hashed with bcrypt (cost ≥ 12).

---
 

**Layering rule:** handlers depend on services; services depend on domain interfaces; postgres repositories implement those interfaces. This keeps the service layer testable via mock repositories and keeps Echo/pgx out of business logic.

---

## 6. API Requirements

| Module | Method | Endpoint | Purpose | Auth | Role |
|---|---|---|---|---|---|
| Auth | POST | `/api/v1/auth/register` | Register email user | No | Any |
| Auth | POST | `/api/v1/auth/login` | Email/social login | No | Any |
| Auth | POST | `/api/v1/auth/refresh` | Exchange refresh token for new access token | No (refresh token in body) | Any |
| Auth | POST | `/api/v1/auth/logout` | Revoke refresh token | Yes | User |
| Auth | POST | `/api/v1/auth/forgot-password` | Request OTP | No | Any |
| Auth | POST | `/api/v1/auth/reset-password` | Verify OTP & set new password | No | Any |
| User | GET | `/api/v1/users/me` | Get profile & settings | Yes | User |
| User | PUT | `/api/v1/users/me` | Update profile info | Yes | User |
| User | PUT | `/api/v1/users/me/password` | Change password | Yes | User |
| User | DELETE | `/api/v1/users/me` | Soft-delete account | Yes | User |
| Device | POST | `/api/v1/devices` | Register/refresh FCM/APNS token | Yes | User |
| Schedule | GET | `/api/v1/schedules/me` | Get current schedule | Yes | User |
| Schedule | PUT | `/api/v1/schedules/me` | Update prayer times/timezone/push toggle | Yes | User |
| Content | GET | `/api/v1/content` | List/search/filter content (`type`, `sub_type`, `audience`, `q`) | Yes | User |
| Content | GET | `/api/v1/content/daily-quote` | Today's quote | Yes | User |
| Content | GET | `/api/v1/content/{id}` | Detail + media URL (premium-gated) | Yes | User |
| Content | GET | `/api/v1/content/{id}/related` | Related content list | Yes | User |
| Subs | GET | `/api/v1/subscriptions/plans` | List active plans (pricing) | Yes | User |
| Subs | POST | `/api/v1/subscriptions/verify` | Verify store receipt | Yes | User |
| Subs | POST | `/api/v1/subscriptions/webhook` | Handle store webhooks | No (store-signed) | Server |

---

## 7. Database Requirements

- **Engine:** PostgreSQL  

---

## 8. Notifications & Real-Time

- FCM (Android) and APNs (iOS), tokens stored per-device in `device_tokens` (a user may have several).
- Dispatch worker: rather than a naive per-minute polling loop, use a scheduled job (cron trigger every minute) that queries `user_schedules` where `morning_prayer_time`/`night_prayer_time` matches the current minute **in that user's stored timezone**, then enqueues sends to a queue (e.g., Redis-backed via `asynq`, Go's equivalent of BullMQ/Celery) for scalable, retryable delivery.

---

## 9. File & Media Management

Claudinary will be used for storing media files

---

## 10. Payment / Subscription

**Flow:** Client purchases via Apple/Google IAP → sends receipt to `/subscriptions/verify` → backend validates against the correct `store` client → upserts `subscriptions` row (linked to `subscription_plans`) → updates `users.is_premium`.

**Webhooks:** `/subscriptions/webhook` handles server-to-server notifications (App Store Server Notifications v2 / Google RTDN) for renewals, cancellations, refunds, billing retries — keeps `is_premium` accurate without relying solely on client-side polling.

---

## 11. Admin / Dashboard Requirements (Inferred)

- Content CMS: create/edit content rows (all types), assign `sub_type`/`category_tag`/`content_audiences`, upload media via presigned URL, schedule `published_at`.
- Subscription plan management: edit `subscription_plans` pricing/features without a deploy.
- User support: view users, manually trigger password reset, handle deletion requests.

---

## 12. Business Rules

- Media URLs for premium content are always short-lived presigned URLs — never stored/returned as permanent public links.
- Email uniqueness enforced at the DB level (`UNIQUE` constraint) in addition to application validation.
- Account deletion: soft-delete immediately (hides from all reads), scheduled job hard-deletes/anonymizes PII after the grace period, satisfying both UI copy and compliance requirements.
- A `content` row may belong to multiple audiences (`content_audiences`) — Kids and Teens sections can render identical underlying content without duplication.

---

## 13. Edge Cases & Error Handling

- **Offline access:** for premium users, provide a `/content/sync?since=<timestamp>` endpoint returning changed/new rows for local caching, plus `ETag`/`Last-Modified` support on media so the client can cache without re-downloading unchanged files.
- **Timezones:** schedule matching must run against each user's stored `timezone`, not server local time — critical since a fixed UTC cron loop would fire prayer reminders at the wrong local hour.
- **Duplicate OTP requests:** invalidate prior unused OTPs for the same user when a new one is requested.
- **Content with no media** (Proverb/Daily Quote): `media_type = NONE`, `media_url = NULL` — client renders text-only.

---

## 14. Non-Functional Requirements

- **Scalability:** notification dispatch decoupled from the cron trigger via a queue (`asynq` + Redis) so send volume doesn't block the scheduler.
- **Security:** rate limiting (Echo middleware, e.g. `golang.org/x/time/rate` or Redis-backed) on `/auth/login`, `/auth/forgot-password`, `/auth/reset-password`.
 
- **Observability:** structured logging (e.g., `zerolog`), request tracing via Echo middleware.

---

 
## 16. Open Questions (unchanged from original analysis, still unresolved by the designs)

- How is content localized when a user switches app language — separate media/text rows per locale, or is `language_preference` purely a UI-string setting for now?
- Does audio/video playback need last-position resume, or is V1 stateless? (Assumed stateless for V1.)
- Any analytics/tracking requirements beyond what's shown?
- Schedule screen time mismatch (see §3.4) — needs design clarification, not a backend assumption.



- You will keep updated the .env and .env.example file 

- you will make a .agent file where you will write the code structures so that every time you can understand codebase quickly

- Write Proper swagger docs for all the API endpoints 
- You will follow the current code structure and code should be neat and clean and professional 


