# Database Design and Ownership

## 1. Overview

The project follows the **Database per Service** pattern.

Each backend service owns its database and is responsible for managing its schema and data.

| Service | Database | Responsibility |
|---|---|---|
| Auth Service | Auth PostgreSQL DB | Authentication credentials and token state |
| User Service | User PostgreSQL DB | User profile information |

The API Gateway does not directly access either database.

## 2. User Service Database

User Service owns the `users` table.

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### Responsibilities

- Stores user profile information.
- Generates the user UUID.
- Enforces profile-level constraints.
- Handles profile creation, retrieval, updates, and deletion.

User Service is the source of ownership for profile information.

## 3. Auth Service Database

Auth Service owns authentication-related data.

### Credentials Table

The Auth Service uses the `credentials` table to store authentication information.

Conceptual schema:

```sql
CREATE TABLE credentials (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

The `id` is supplied by User Service during signup.

### Responsibilities

- Stores password hashes.
- Stores authentication-related email information.
- Maintains user roles, if configured.
- Supports signin and credential verification.

### Refresh Tokens Table

Auth Service also owns refresh-token records.

Conceptual fields:

| Column | Purpose |
|---|---|
| id | Refresh-token record identifier |
| user_id | Logical user UUID |
| token_hash | Hash of the refresh token |
| expires_at | Token expiration time |
| revoked | Whether the token has been revoked |
| created_at | Record creation time |

The actual schema must match the Auth Service migrations.

Refresh tokens should be stored as cryptographic hashes rather than plaintext tokens.

## 4. Shared User Identity

Both databases use the same logical user UUID.

Example:

**User DB**

```text
users
--------------------------------
id       = 8d3f...
name     = John
email    = john@example.com
```

**Auth DB**

```text
credentials
--------------------------------
id            = 8d3f...
email         = john@example.com
password_hash = <hashed-password>
```

The UUID connects the two records logically.

There is no foreign key between the databases.

Auth Service receives the UUID from User Service through the `CreateUser` gRPC response.

## 5. Why Separate Databases?

Database ownership provides several benefits:

### Independent Ownership

Each service controls its own schema, migrations, and persistence logic.

### Loose Coupling

Auth Service does not depend on the internal table structure of User Service.

It depends on the User Service gRPC contract instead.

### Independent Evolution

A service can evolve its database schema without requiring another service to query or understand its internal tables.

### Independent Scaling

Database resources can be scaled according to each service's workload.

## 6. Cross-Service Data Consistency

Separate databases introduce consistency challenges.

For example, during signup:

1. User Service creates a profile.
2. Auth Service creates credentials.
3. Both records use the same UUID.

If User Service succeeds but Auth Service fails, the databases become inconsistent.

There is no single SQL transaction covering both databases.

The application must handle partial failures through an explicit recovery strategy.

Possible approaches include:

- Saga pattern.
- Compensating operations.
- Pending registration.
- Reconciliation.

## 7. Schema Migration Strategy

Each service maintains its own database migrations.

Example:

```text
auth-service/
└── migrations/
    ├── 000001_create_credentials.up.sql
    ├── 000001_create_credentials.down.sql
    ├── 000002_create_refresh_tokens.up.sql
    └── 000002_create_refresh_tokens.down.sql

user-service/
└── migrations/
    ├── 000001_create_users.up.sql
    └── 000001_create_users.down.sql
```

Migration responsibilities:

- Each service applies migrations to its own database.
- Schema changes must be version-controlled.
- Up migrations apply changes.
- Down migrations provide rollback operations where safe.
- SQL queries and store implementations must be updated when schemas change.

A table rename does not automatically rename its columns. For example:

```sql
ALTER TABLE users RENAME TO credentials;
```

The existing `id` column remains named `id`.

## 8. Database Access Rules

1. Auth Service directly accesses only Auth DB.
2. User Service directly accesses only User DB.
3. API Gateway does not directly access backend databases.
4. Services communicate through their gRPC APIs.
5. Cross-service UUID relationships are maintained at the application level.
6. Each service owns its database migrations.

## 9. Future Improvements

Potential improvements include:

- Database connection pooling and monitoring.
- Transaction and query performance monitoring.
- Index review based on actual query patterns.
- Backup and restore procedures.
- Migration testing in CI.
- Handling incomplete registrations.
- Database availability and recovery testing.