# DESIGN.md

## Architecture

The architecture of this project is an `idiomatic/layerd structure` which seperates the project into different layers including `cmd` layer, `internal` layer which itself contains the `api`(`routes`, `handlers`), `config`, `db`, `errors`, `reponses`, `models`, `repository` and the `service` layers, along with the auto-generated `swagger` api documentation layer in `docs`.

### Structure layout

``` txt
cmd/              <-- The entry point of your app (main.go lives here)
internal/         <-- All your application-specific code (private to your project)
    api/          <-- HTTP-related logic (Gin routes, handlers)
        routes/         <-- Route definitions (URL paths, groupings)
        handler/        <-- HTTP handlers (business logic triggered by API calls)

    config/       <-- Configuration (env vars, settings, constants)

    db/           <-- Database layer
        sql_commands/   <-- Raw SQL queries (if any migrations/scripts are here)

    errors/       <-- Custom error types & error handling logic

    responses/    <-- API response formats (e.g. success/error JSON templates)

    models/       <-- Data models (GORM structs for DB tables)

    repository/   <-- Data access layer (queries, CRUD, DB transactions)
    
    service/      <-- Business logic layer (pure Go logic, independent of HTTP or DB)
docs/             <-- Swagger Documentation 
```
