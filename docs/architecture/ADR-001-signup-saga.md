# ADR-001: Signup Saga Orchestration

- **Status:** Accepted
- **Pattern:** Saga Orchestration
- **Orchestrator:** Auth Service
- **Participants:** Auth Service and User Service

## Context

Signup involves two independent services and databases.

Auth Service requests profile creation from User Service. User Service commits the profile in its own database.

Auth Service subsequently creates credentials and refresh-token records in its database.

If Auth DB operations fail, the Auth transaction can roll back, but the User DB profile remains committed.

## Decision

Implement a synchronous Saga with Auth Service acting as the orchestrator.

## Workflow

1. Auth Service requests profile creation from User Service.
2. User Service creates and commits the profile.
3. Auth Service starts a database transaction.
4. Auth Service creates credentials and refresh-token records.
5. Auth Service commits the transaction.
6. Signup succeeds.

## Compensation

If a known pre-commit Auth DB operation fails:

1. Roll back the Auth DB transaction if one exists.
2. Call the User Service compensation RPC.
3. User Service removes the profile created by this signup attempt.
4. Return the original signup failure.

Compensation must not delete a pre-existing profile.

If compensation fails, preserve the original failure and record the compensation failure for investigation.

If Auth DB commit outcome is ambiguous, do not blindly compensate. The result must be treated as uncertain.

## Scope

This initial implementation includes:
- Saga orchestration in Auth Service.
- Profile creation in User Service.
- A dedicated compensating RPC.
- Tests for successful signup and compensation on failure.

This initial implementation excludes:
- Transactional Outbox.
- Kafka-based Saga messaging.
- Persistent workflow state.
- Automatic retry and reconciliation.
- Distributed workflow recovery after process crashes.

These may be considered as future improvements.