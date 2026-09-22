# Prayer Schedule QA Test Plan

This document outlines the comprehensive test cases and frontend implementation checklist for the "Prayer Schedule" domain, ensuring flawless integration and robust security.

## 1. Default Schedule Creation Check

**Objective**: Verify that a newly registered user who has never configured a schedule receives the correct default values on their first GET request, rather than a `404 Not Found`.

**cURL Request**:
```bash
curl -X GET "http://localhost:8080/api/v1/schedules/me" \
  -H "Authorization: Bearer <NEW_USER_JWT>" \
  -H "Accept: application/json"
```

**Expected Result (Status 200 OK)**:
```json
{
  "id": "generated-uuid-here",
  "morning_prayer_time": "05:00",
  "night_prayer_time": "21:00",
  "push_enabled": true,
  "timezone": "UTC",
  "user_id": "<NEW_USER_ID>",
  "updated_at": "2026-09-22T00:00:00Z"
}
```
> [!IMPORTANT]
> The backend logic must intercept the `record not found` error from the DB, create this default record, save it, and return it.

---

## 2. UI Payload Mapping & Logic Verification

> [!TIP]
> **Checklist for the Frontend/Mobile Developer**

- [ ] **Start Time vs. End Time**: The API only tracks `morning_prayer_time` and `night_prayer_time` (the start times). The UI displays an 8-minute duration (e.g., `09:30 AM` to `09:38 AM`). The frontend must dynamically calculate the End Time by adding 8 minutes to the retrieved start time for display purposes. There is no need to send the End Time to the backend.
- [ ] **Time Format Conversion**: The API strictly expects and returns the 24-hour `HH:mm` string format (e.g., `"14:30"`). The UI displays a 12-hour AM/PM format. The frontend must parse the API's 24-hour string, format it for the user in 12-hour format, and convert it back to 24-hour format before sending the PUT request.
- [ ] **Timezone Handling**: The frontend should detect the device's native timezone using standard libraries (e.g., `Intl.DateTimeFormat().resolvedOptions().timeZone` in JS/React Native) and send that string (e.g., `"America/New_York"`) in the `timezone` field to ensure push notifications trigger at the correct local time.

---

## 3. Input Validation Tests

### A. Invalid Timezone String
**Objective**: Ensure the backend rejects timezone strings that are not valid IANA identifiers.

**cURL Request**:
```bash
curl -X PUT "http://localhost:8080/api/v1/schedules/me" \
  -H "Authorization: Bearer <VALID_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "morning_prayer_time": "08:00",
    "night_prayer_time": "22:00",
    "push_enabled": true,
    "timezone": "America/Invalid"
  }'
```

**Expected Result (Status 400 Bad Request)**:
```json
{
  "code": 400,
  "message": "Validation failed",
  "details": "timezone must be a valid IANA timezone string"
}
```

### B. Invalid Time Formats
**Objective**: Ensure the backend strictly enforces the `HH:mm` format and rejects out-of-bounds or malformed times.

**cURL Request**:
```bash
curl -X PUT "http://localhost:8080/api/v1/schedules/me" \
  -H "Authorization: Bearer <VALID_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "morning_prayer_time": "25:00",
    "night_prayer_time": "8 PM",
    "push_enabled": true,
    "timezone": "America/New_York"
  }'
```

**Expected Result (Status 400 Bad Request)**:
```json
{
  "code": 400,
  "message": "Validation failed",
  "details": "morning_prayer_time must be in valid HH:mm format, night_prayer_time must be in valid HH:mm format"
}
```

---

## 4. Authentication & Security Tests

### A. Unauthenticated Access
**Objective**: Ensure the endpoint cannot be accessed without a valid JWT.

**cURL Request**:
```bash
curl -X GET "http://localhost:8080/api/v1/schedules/me" \
  -H "Accept: application/json"
```

**Expected Result (Status 401 Unauthorized)**:
```json
{
  "message": "missing or malformed jwt"
}
```

### B. Insecure Direct Object Reference (IDOR) Attempt
**Objective**: Because the route is `/me`, it strictly relies on the JWT payload. We will attempt to maliciously inject a `user_id` into the PUT body to ensure the backend ignores it and doesn't modify another user's schedule.

**cURL Request**:
```bash
curl -X PUT "http://localhost:8080/api/v1/schedules/me" \
  -H "Authorization: Bearer <USER_A_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "morning_prayer_time": "08:00",
    "night_prayer_time": "22:00",
    "push_enabled": true,
    "timezone": "America/New_York",
    "user_id": "<USER_B_ID>"
  }'
```

**Expected Result (Status 200 OK)**:
The backend should process the update, but a subsequent `GET` using `<USER_A_JWT>` must verify that the schedule was updated for **User A**, and User B's schedule remains untouched. The injected `user_id` must be completely ignored by the Go server logic during the struct binding and update phases.
