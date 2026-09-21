#!/bin/bash

# Get Token
TOKEN=$(curl -s -X POST http://localhost:5525/api/v1/auth/admin/login -H "Content-Type: application/json" -d '{"email":"admin@altar.com","password":"Admin1234"}' | jq -r '.data.tokens.access_token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  TOKEN=$(curl -s -X POST http://localhost:5525/api/v1/auth/admin/login -H "Content-Type: application/json" -d '{"email":"admin@altar.com","password":"Admin1234"}' | jq -r '.access_token')
fi

# Create Page
curl -s -X POST http://localhost:5525/api/v1/admin/cms/pages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "terms-and-conditions",
    "title": "Terms & Condition",
    "intro_text": "By accessing or using Eventify, you agree to be bound by these Terms & Conditions. If you do not agree, please refrain from using the app.",
    "sections": [
      {
        "heading": "1. Acceptance of Terms",
        "content": "By accessing or using Eventify, you agree to be bound by these Terms & Conditions. If you do not agree, please refrain from using the app.",
        "sort_order": 0
      },
      {
        "heading": "2. Use of the App",
        "content": "Eventify provides a platform that connects users with event organizers, venues, and ticket sellers. You agree to use the app responsibly and only for lawful purposes.",
        "sort_order": 1
      },
      {
        "heading": "3. Event Services",
        "content": "• Event availability depends on venue capacity.\n• Tickets for events are subject to confirmation before purchase.\n• An event can only be booked after the payment is successfully processed.",
        "sort_order": 2
      },
      {
        "heading": "4. Payments",
        "content": "• Payments can be made online through the app or in cash at the venue.\n• Online payments may include a service fee.\n• Cash payments are recorded for transparency but collected by the venue staff.\n• All completed bookings are logged in the system.",
        "sort_order": 3
      },
      {
        "heading": "6. Organizer Responsibilities",
        "content": "• Event organizers must be verified before using the platform.\n• Subscription status determines account access and event management capabilities.",
        "sort_order": 4
      }
    ]
  }' | jq .
