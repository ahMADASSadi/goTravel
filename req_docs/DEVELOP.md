# Multi-Transport Ticketing System Roadmap

## Overview

This document outlines how to extend the existing **bus ticketing system** into a platform that supports **buses**, **planes**, **ships**, and **trains**. We will build on the current architecture (using Gin, GORM) and add new features to support different transport types.

---

## Current Architecture

Right now, the project is organized like this:

```txt
cmd/
  main.go                <-- The main entry point of the app
internal/
  api/
    routes/              <-- Defines API routes (Gin)
    handler/             <-- Handles HTTP requests
  config/                <-- Settings and configuration
  db/
    sql_commands/        <-- Database migration files
  errors/                <-- Centralized error handling
  responses/             <-- Standard API responses
  models/                <-- Database models (Bus, Schedule, etc.)
  repository/            <-- Manages database access
  service/               <-- Business logic
```

The structure is clean and easy to extend.

---

## How We’ll Add Multi-Transport Support

### Step 1: Generalize Models

We need to update the current bus-specific models to support multiple transport types.

#### Create a `Transport` Model

Instead of having separate models for each transport type, we’ll use one model for all:

```go
type Transport struct {
    ID           uint
    Type         string // "bus", "airplane", "ship", "train"
    TransportNo  string // Transport number (e.g., license plate or flight number)
    Capacity     uint   // Number of seats
    Seats        []Seat // List of seats available
}
```

- `Type` will indicate what kind of transport it is (bus, airplane, etc.).
- `TransportNo` is a unique identifier (like a bus number or flight number).

#### Specialized Tables

We’ll still have specialized tables for things like bus details, airplane details, etc., but these will link back to the common `Transport` table.

#### Generalize Schedule Model

We’ll update the schedule model to support any transport type:

```go
type Schedule struct {
    ID              uint
    TransportID     uint
    OriginCode      string
    DestinationCode string
    DepartureTime   time.Time
    ArrivalTime     time.Time
}
```

---

### Step 2: Refactor Services and Repositories

We need to update the code that interacts with the database to handle the new `Transport.Type` field. Each transport type (bus, plane, ship, train) will have its own special rules and business logic:

- **Bus**: Direct seat selection.
- **Airplane**: Class-based seat selection (economy/business).
- **Ship**: Cabin vs free seating.
- **Train**: Carriage and seat mapping.

---

### Step 3: New API Endpoints

We’ll add new endpoints to support multi-transport searches and reservations:

| Endpoint                                | Description                     |
| --------------------------------------- | ------------------------------- |
| `GET /api/v1/transport/{type}/search`   | Search trips by transport type  |
| `POST /api/v1/transport/{type}/reserve` | Reserve a ticket for a transport |
| `GET /api/v1/tickets/{id}`              | Get ticket and transport info   |

Here, `{type}` can be `bus`, `airplane`, `ship`, or `train`.

---

### Step 4: Business Rules for Each Transport

Each transport type has its own set of rules:

| Transport | Business Rules                                |
| --------- | --------------------------------------------- |
| Bus       | Simple seat selection.                        |
| Airplane  | Seat classes, luggage rules, boarding passes. |
| Ship      | Cabin vs free seating, port schedules.        |
| Train     | Carriage mapping, multi-leg journeys.         |

---

## Folder Structure Changes

We will add new files and models to accommodate the multi-transport support:

```txt
internal/
  models/
    transport.go           <-- General transport model
    vehicle_bus.go         <-- Bus-specific details
    vehicle_plane.go       <-- Airplane-specific details
  api/
    handler/
      transport_handler.go <-- Handles transport-related requests
    routes/
      transport_routes.go  <-- Defines routes for transport operations
```

---

## Future Features (Advanced)

- **User accounts** to manage reservations.
- **Payment integration** (e.g., Stripe, PayPal).
- **QR codes** for boarding passes.
- **Notifications** (email/SMS).
- **Multi-language support**.
- Support for **multi-leg journeys** (especially for trains and ships).

---

## Next Steps

1. Refactor the current models to support multiple transport types.
2. Update the business logic to handle different transport types.
3. Extend API routes to support searching and reserving tickets for buses, planes, ships, and trains.
4. Test each transport type separately to ensure everything works smoothly.
5. Add any specific features for each transport type (e.g., seat classes, cabins, etc.).
