# JustDoc — Product Summary

## Vision

To provide developers with a simple, reliable way to store and retrieve JSON documents via a clean API, without the overhead of setting up and managing a database.

## Target Audience

Frontend developers building apps that need persistent data storage for prototypes and small applications.

## Core Features

- `POST /<channel>/<document>` — Store or update a JSON document
- `GET /<channel>/<document>` — Retrieve a JSON document
- `GET /<channel>/` — List all documents in a channel
- `GET /` — List all channels

## Key Constraint

No authentication for v1 — designed for simplicity and rapid integration.
