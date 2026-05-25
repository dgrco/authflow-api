package ctxkeys

// This file exists to avoid awkward dependency management where files would otherwise
// depend on middleware just to extract this type.

type StringContextKey string
