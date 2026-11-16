# Failure Handling and Scenarios

## 1. Purpose

A request in a microservices system may depend on multiple services,
databases, and network calls. A failure in one component can prevent the
overall request from completing, even if another service has already
committed its work.

This document explains failure behavior across:

-   Signup
-   Signin
-   Refresh token
-   Signout
-   Get, update, and delete user profile
-   Service communication and shutdown

It records the current architecture, expected behavior, and gaps to
verify or implement. Recommendations in this document are not claims
that the mechanisms already exist.

## 2. System context

-   **API Gateway** receives HTTP requests and calls backend services
    over gRPC.
-   **Auth Service** owns credentials, password hashes, roles, and
    refresh-token records.
-   **User Service** owns user profile data.
-   Auth Service and User Service use separate databases.
-   During signup, Auth Service calls User Service to create a profile,
    then stores credentials in Auth DB using the profile ID returned by
    User Service.
-   Signin is handled by Auth Service and does not need to call User
    Service.

Each service owns its database. A local transaction in one service
cannot roll back a transaction already committed by another service.

## 3. Failure-handling principles

1.  Use context deadlines for inter-service calls and database
    operations.
2.  Fail safely when a required dependency is unavailable.
3.  Retry only operations that are safe to retry.
4.  Do not retry permanent failures such as invalid input or invalid
    credentials.
5.  Keep database transactions short.
6.  Never assume a failed RPC means the remote operation did not
    execute.
7.  Log failures with request IDs and operation context.
8.  Do not expose internal database errors, connection details, or
    secrets to clients.
9.  Preserve service database ownership; do not work around a failed
    service by writing directly to its database.
10. Make cross-service consistency workflows explicit.

## 4. Signup failures

Signup crosses Auth Service and User Service.

### 4.1 User Service is unavailable

**Flow:** 1. Gateway sends signup to Auth Service. 2. Auth Service calls
`UserService.CreateUser`. 3. User Service cannot complete the RPC.

**Expected behavior:** - Auth Service must not proceed to create
credentials if it did not receive a successful profile-creation
result. - Return a controlled dependency/service error. - Do not expose
internal connection details.

**Suggested gRPC status:** `Unavailable` when the dependency is
unavailable.

**Current-code note:** The Auth Service returns the downstream error.
Verify that the outbound gRPC call has a suitable deadline and that the
handler maps the error intentionally.

### 4.2 Profile is created, but Auth Service later fails

**Flow:** 1. User Service creates and commits the profile. 2. Auth
Service hashes the password and begins its Auth DB transaction. 3. A
later Auth DB operation fails, or token generation fails. 4. Auth DB
transaction rolls back if it has started and has not committed.

**Problem:** The profile remains in User DB. Auth DB rollback cannot
undo the User DB commit.

**Possible future solutions:** - A compensating operation to remove or
mark the newly created profile. - A pending-registration workflow. - A
durable workflow with retries and reconciliation.

A compensation can also fail, so it must be observable and recoverable.
Do not describe a Saga or compensation as implemented unless the code
actually implements it.

### 4.3 Signup is retried after a timeout or lost response

A client may not know whether signup completed.

For example: 1. User Service creates the profile. 2. Auth Service
commits credentials and refresh-token data. 3. The response is lost. 4.
The client retries signup.

A unique email constraint helps prevent duplicate rows within a
database, but it does not make the whole cross-service workflow
idempotent.

**Possible future solutions:** - An idempotency key with a defined
retention and replay contract. - Idempotent profile creation or a way to
retrieve the result of the original operation. - Reconciliation for
profile records that have no matching credentials. - Stable handling of
duplicate-registration outcomes.

Do not blindly retry profile creation after an ambiguous timeout.

## 5. Signin failures

Signin primarily depends on Auth Service and Auth DB.

### 5.1 Auth DB is unavailable

**Expected behavior:** - Do not issue tokens without successfully
verifying credentials. - Return a controlled service error. - Log the
database failure without exposing SQL details.

**Suggested gRPC status:** `Unavailable` for a dependency outage.

### 5.2 Invalid email or password

**Expected behavior:** - Return a generic authentication failure. - Do
not reveal whether the email exists. - Do not issue access or refresh
tokens.

**Suggested gRPC status:** `Unauthenticated`.

### 5.3 Token generation or refresh-token persistence fails

**Expected behavior:** - Do not return a successful signin response. -
Do not return tokens if the required persistence step failed. - Log the
internal failure safely.

**Current-code note:** The shared signin implementation stores the
refresh token through the DB without showing an explicit transaction
around the whole signin operation. Do not claim that all signin writes
are rolled back atomically unless the implementation is changed and
verified.

## 6. Refresh-token failures

Refresh-token rotation needs special care because the same token may be
presented concurrently.

### 6.1 Invalid, expired, or revoked token

**Expected behavior:** - Reject the request. - Do not issue replacement
tokens. - Return an authentication error.

**Suggested gRPC status:** `Unauthenticated`.

### 6.2 Auth DB is unavailable

**Expected behavior:** - Do not issue a replacement token if validation
or persistence cannot be completed. - Return a controlled service
error. - A retry must account for whether the original rotation may have
committed.

**Suggested gRPC status:** `Unavailable`.

### 6.3 Two requests use the same refresh token concurrently

If both requests validate the token before either revokes it, both may
attempt to rotate it.

A conditional update can ensure that only one request consumes the old
token:

``` sql
UPDATE refresh_tokens
SET revoked = TRUE
WHERE token_hash = $1
  AND revoked = FALSE
  AND expires_at > NOW()
RETURNING user_id;
```

The update and insertion of the replacement token must happen in the
same Auth DB transaction. If no row is returned, the token was invalid,
expired, or already consumed.

**Current-code note:** The previously shared service code starts a
transaction for rotation, but passes the raw refresh token to store
methods that were designed to work with hashes. Hashing must be
consistent during creation, lookup, and revocation. Verify the store
implementation and tests before treating rotation as production-ready.

## 7. Signout failures

Signout revokes a refresh token.

### 7.1 Invalid or already-revoked token

Repeated signout requests should have a deliberate, consistent contract.
For example, an API may treat an already-revoked token as successfully
signed out, or return an authentication error. Choose and document one
behavior.

### 7.2 Auth DB is unavailable

Do not report successful revocation if the database operation failed.
Return a controlled error. A retry may be safe if revocation is
implemented idempotently.

### 7.3 Access-token behavior

Revoking a refresh token does not automatically invalidate an
already-issued stateless access token. It may remain usable until expiry
unless the system adds access-token revocation or online validation.

## 8. User-profile failures

Profile operations are handled by User Service.

### 8.1 User Service is unavailable

Gateway cannot complete a profile operation if User Service is required
and unavailable.

**Expected behavior:** - Return a controlled dependency error. - Do not
access User DB directly from Gateway or Auth Service as a workaround.

**Suggested gRPC status:** `Unavailable`.

### 8.2 User does not exist

Return a not-found result rather than exposing a database-specific
`sql.ErrNoRows` error.

**Suggested gRPC status:** `NotFound`.

### 8.3 User accesses another user's profile

The service must authorize the requested resource, not merely check that
a JWT exists.

**Expected behavior:** - Compare the authenticated identity with the
requested user ID. - Permit access to another user's profile only when
an explicit authorization rule allows it, such as an administrative
role. - Reject unauthorized access.

**Suggested gRPC status:** `PermissionDenied`.

### 8.4 Profile update fails

If the update uses a local database transaction, report success only
after commit. If no transaction is used, describe the actual
single-statement behavior rather than claiming a rollback.

**Email consistency concern:** Email is stored in both Auth Service and
User Service in the current design. If profile update changes email in
only User Service, the records can diverge. An email change therefore
needs an explicit cross-service workflow, or a decision about which
service is authoritative for that attribute.

### 8.5 Profile deletion partially succeeds

Deleting a profile in User Service does not automatically delete
credentials or revoke refresh tokens in Auth Service.

**Possible future solutions:** - A coordinated deletion workflow, such
as a Saga. - A pending-deletion state followed by cleanup and
reconciliation. - An Auth Service operation that disables credentials
and revokes refresh tokens, coordinated with profile deletion.

Do not assume that deleting one service's record deletes another
service's data.

## 9. Retry strategy

Retries can help with transient failures but can also duplicate writes.

  -----------------------------------------------------------------------
  Failure                             General policy
  ----------------------------------- -----------------------------------
  Invalid request                     Do not retry

  Invalid credentials                 Do not retry

  User not found                      Do not retry

  Duplicate email                     Do not retry blindly

  Temporary network failure           Bounded retry only when safe

  Service unavailable                 Bounded retry only when safe

  Deadline exceeded                   Treat outcome as potentially
                                      unknown

  Database deadlock                   Retry the transaction when safe and
                                      within the request deadline
  -----------------------------------------------------------------------

Recommended practices: - Use a small maximum attempt count. - Use
exponential backoff with jitter. - Respect context cancellation and
deadlines. - Avoid retrying independently at every service layer. - Use
idempotency or reconciliation for operations that may be repeated. - Do
not retry indefinitely.

A deadline-exceeded error does not prove that the remote operation
failed. The remote service may have committed its work while the
response was delayed or lost.

## 10. gRPC error mapping

Use meaningful gRPC status codes instead of returning every error as
`Internal`.

  Situation                     Suggested status
  ----------------------------- --------------------
  Invalid request               `InvalidArgument`
  Invalid credentials           `Unauthenticated`
  Missing permission            `PermissionDenied`
  Resource not found            `NotFound`
  Duplicate resource            `AlreadyExists`
  Dependency unavailable        `Unavailable`
  Deadline exceeded             `DeadlineExceeded`
  Unexpected internal failure   `Internal`

These are recommended mappings. Verify the actual handler and service
error mapping before claiming they are implemented.

Log detailed internal errors on the server side while returning safe
messages to clients.

## 11. Observability

Failures should be diagnosable across Gateway, Auth Service, and User
Service.

Recommended practices: - Propagate a request or trace ID across HTTP and
gRPC boundaries. - Log downstream RPC failures and durations. - Record
request counts, latency, and error counts. - Monitor database
connection-pool health and transaction failures. - Track partial signup
and deletion workflows. - Track refresh-token failures and suspicious
reuse. - Add distributed tracing when the tracing infrastructure is
implemented. - Never log passwords, access tokens, raw refresh tokens,
or other secrets.

Useful metrics include: - HTTP request count, duration, and in-flight
requests. - gRPC request count, duration, and status-code counts. -
Database query errors and connection-pool wait metrics. - Signup
failures and partial-registration counts. - Refresh-token rotation
failures.

## 12. Graceful shutdown

Graceful shutdown is service-lifecycle handling rather than a single
request failure, but it matters because deployments and restarts can
interrupt in-flight work.

A service should: 1. Stop accepting new work. 2. Allow in-flight
requests to finish within a configured grace period. 3. Cancel or
terminate remaining work after the deadline. 4. Close gRPC
listeners/connections and database resources. 5. Exit cleanly.

Verify the actual shutdown code for each application. Do not claim
graceful shutdown is implemented merely because the process handles a
signal.

## 13. Additional security boundary to verify

User Service's `CreateUser` RPC was previously exempted from the User
Service JWT interceptor. That means any client able to reach that RPC
may be able to invoke it unless another protection exists.

Public signup can still enter through Gateway, but the internal
Auth-Service-to-User-Service call should have an explicit
service-to-service trust mechanism, such as mTLS or authenticated
service credentials, plus authorization at the RPC boundary.

End-user JWT authentication and service-to-service authentication solve
different problems.

## 14. Production-readiness checklist

-   [ ] Configure deadlines for outbound gRPC calls.
-   [ ] Propagate request/trace IDs across service boundaries.
-   [ ] Map application errors to intentional gRPC and HTTP responses.
-   [ ] Test behavior when a dependency is unavailable.
-   [ ] Test signup when profile creation succeeds but Auth DB work
    fails.
-   [ ] Define safe behavior for signup retries after ambiguous
    outcomes.
-   [ ] Define and test idempotency for retryable write operations.
-   [ ] Make refresh-token rotation atomic.
-   [ ] Hash refresh tokens consistently at create, lookup, and revoke.
-   [ ] Prevent concurrent use of the same refresh token from succeeding
    twice.
-   [ ] Define signout behavior for already-revoked tokens.
-   [ ] Coordinate email changes across service ownership boundaries.
-   [ ] Coordinate profile deletion, credential disabling, and token
    revocation.
-   [ ] Enforce profile authorization, not just token presence.
-   [ ] Secure service-to-service RPCs, including `CreateUser`.
-   [ ] Add structured logs and useful metrics without logging secrets.
-   [ ] Verify graceful shutdown in Gateway, Auth Service, and User
    Service.

## 15. Scope and implementation status

This document is a failure analysis and design guide based on the
architecture and code reviewed so far.

It distinguishes: - **Current behavior:** behavior visible in the shared
implementation. - **Expected behavior:** behavior the system should
provide. - **Gap / future work:** behavior that still needs
implementation or verification.

Do not describe Saga, compensation, idempotency, retries, circuit
breakers, distributed tracing, or graceful shutdown as implemented
unless the corresponding code and tests exist.

Introduce resilience mechanisms based on actual requirements and
measured failure modes. Not every operation needs a Saga, retries, or a
distributed transaction.
