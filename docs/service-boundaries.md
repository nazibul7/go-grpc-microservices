# Service Boundaries and Data Ownership

## 1. Purpose

A microservice should own a specific business capability, its business logic, and the data required to support that capability.

Other services should interact with it through well-defined APIs rather than accessing its internal implementation or database directly.

This document defines the boundaries between the API Gateway, Auth Service, and User Service.

## 2. API Gateway

### Responsibilities

- Exposes HTTP endpoints to external clients.
- Routes authentication requests to Auth Service.
- Routes profile-related requests to User Service.
- Maintains gRPC client connections to backend services.
- Converts HTTP requests into gRPC requests.
- Converts gRPC responses into HTTP responses.
- Handles API-level concerns such as request routing and HTTP error mapping.

### Data Ownership

The API Gateway does not own user or authentication data.

It should not directly access either backend database.

### Dependencies

- Auth Service
- User Service

## 3. Auth Service

### Responsibilities

- Handles user registration.
- Authenticates users during signin.
- Hashes and verifies passwords.
- Manages authentication credentials.
- Generates access tokens.
- Manages refresh-token creation, rotation, and revocation.
- Coordinates user profile creation by calling User Service during signup.

### Data Ownership

Auth Service owns authentication-related data.

| Data | Description |
|---|---|
| User ID | Logical identifier shared with User Service |
| Email | Authentication lookup information |
| Password hash | Hashed user password |
| Role | Authorization role, if configured |
| Refresh tokens | Token hash, expiry, and revocation state |

### Database Ownership

Auth Service owns Auth DB and its authentication-related tables, including:

- `credentials`
- `refresh_tokens`

The exact tables depend on the current migrations.

### Dependencies

- Auth PostgreSQL DB
- User Service through gRPC

### Boundary Rules

- Auth Service must not directly query User Service's database.
- Auth Service owns password verification.
- Auth Service owns refresh-token lifecycle operations.
- Auth Service may call User Service APIs when a workflow requires profile operations.

## 4. User Service

### Responsibilities

- Creates user profiles.
- Retrieves user profile information.
- Updates profile information.
- Deletes profiles, if supported by the API.
- Exposes profile operations through gRPC.

### Data Ownership

User Service owns profile-related data.

| Data | Description |
|---|---|
| User ID | Unique profile identifier |
| Name | User's profile name |
| Email | User's profile email |
| Created At | Profile creation timestamp |
| Updated At | Last profile update timestamp |

### Database Ownership

User Service owns User DB and its profile-related tables, including:

- `users`

### Dependencies

- User PostgreSQL DB

### Boundary Rules

- User Service must not directly access Auth Service's database.
- User Service does not store or verify passwords.
- User Service does not manage refresh tokens.
- Profile operations must be handled by User Service.

## 5. Data Ownership Matrix

| Resource | API Gateway | Auth Service | User Service |
|---|---|---|---|
| HTTP routing | Owns | No | No |
| User profile | No | No | Owns |
| Password hash | No | Owns | No |
| Credential records | No | Owns | No |
| Access-token generation | No | Owns | No |
| Refresh-token lifecycle | No | Owns | No |
| Profile database | No | No | Owns |
| Auth database | No | Owns | No |

## 6. Shared User Identity

Both Auth Service and User Service maintain records associated with the same logical user UUID.

Example:

```text
User DB
-------
users
  id = 8d3f...
  name = John
  email = john@example.com

Auth DB
-------
credentials
  id = 8d3f...
  email = john@example.com
  password_hash = ...
```

The UUID connects the logical identity across services.

There is no direct foreign key between the two databases.

Each service remains responsible for its own record and lifecycle.

## 7. Why Auth Service Calls User Service During Signup

Registration requires creating both:

1. A user profile.
2. Authentication credentials.

Auth Service coordinates the registration workflow because it owns the authentication process.

It calls User Service through gRPC to create the profile, then stores credentials in its own database using the returned user ID.

This is service-to-service communication, not direct database sharing.

Signin does not require the same call because Auth Service already owns the credentials needed to verify the user's password.

## 8. Rules for Maintaining Service Boundaries

When adding new functionality, follow these rules:

1. A service should directly access only its own database.
2. Other services should use the owning service's API.
3. Do not duplicate business logic across services.
4. Keep protobuf contracts focused on service operations.
5. Avoid making one service dependent on another service's internal database schema.
6. Keep shared identifiers consistent across services.
7. Consider service availability and partial failures whenever one service calls another.

## 9. Future Considerations

As the system evolves, review:

- Whether new functionality belongs to an existing service or requires a separate service.
- Whether synchronous communication is appropriate for each workflow.
- How data consistency is maintained across services.
- How service failures affect dependent operations.
- Whether API contracts can evolve without breaking existing clients.