# API Versioning Strategy

This document outlines the versioning strategy for the Resume API, explaining how API versions are managed, how clients can request specific versions, and how developers should implement new versions.

## Overview

The Resume API uses a versioning strategy to ensure backward compatibility while allowing the API to evolve. The versioning system supports multiple methods for clients to specify which API version they want to use.

## Supported Versions

Currently supported API versions:

- `v1` - Initial API version

## How to Request a Specific Version

Clients can specify which API version they want to use in three ways:

### 1. URI Path (Recommended)

Include the version in the URI path:

```
GET /api/v1/profile
```

### 2. Accept Header

Specify the version in the Accept header:

```
Accept: application/json;version=1
```

### 3. Query Parameter

Add a version query parameter:

```
GET /api/profile?version=1
```

## Version Negotiation

The API uses the following process to determine which version to use:

1. Check the URI path for a version
2. If not found, check the Accept header
3. If not found, check the version query parameter
4. If no version is specified, use the latest stable version

## Implementing New API Versions

When implementing a new API version, follow these guidelines:

### When to Create a New Version

Create a new API version when making breaking changes, such as:

- Removing or renaming fields
- Changing field types or formats
- Changing response structure
- Removing endpoints
- Changing endpoint behavior in incompatible ways

### How to Implement a New Version

1. Add the new version to the `Version` constants in `internal/versioning/versioning.go`
2. Create version-specific handlers if needed
3. Register routes for the new version using the versioning router
4. Update tests to cover the new version
5. Update documentation to describe the changes

### Example: Adding a New Version

```go
// 1. Add the new version constant
const (
    V1 Version = "v1"
    V2 Version = "v2"
    LatestVersion = V2
)

// 2. Update the All() function
func All() []Version {
    return []Version{V1, V2}
}

// 3. Register routes for both versions
versionedRouter := versioning.NewRouter(router)
versionedRouter.Group(versioning.V1).GET("/profile", v1Handler.GetProfile)
versionedRouter.Group(versioning.V2).GET("/profile", v2Handler.GetProfile)

// 4. Or register the same handler for multiple versions if the functionality is the same
versionedRouter.RegisterVersionedEndpoint("/skills", "GET", []versioning.Version{versioning.V1, versioning.V2}, handler.GetSkills)
```

## Deprecation Policy

When a new API version is released:

1. The previous version will be supported for at least 12 months
2. Deprecation notices will be included in API responses for deprecated versions
3. Documentation will clearly mark deprecated versions
4. Clients will be encouraged to migrate to the newest version

## Best Practices for API Consumers

- Always specify the API version you're using
- Check for deprecation notices in API responses
- Test your integration with new API versions before migrating
- Subscribe to the API changelog for updates on new versions and deprecations