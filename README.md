# Go Templates

Template caching and paging utilities. Require go 1.22+

## Usage

```shell
go get github.com/mawngo/go-tmpls
```

## Template Caching

Cache the template for re-execution without having to parse it again, support template reload for development.

See [examples](/examples/main.go) for setup and integrating template cache.

### Built-in template funcs

By default, this library adds some [helpers](/internal/builtin.go) to the template.
To disable all built-in functions use`WithoutBuiltins()`, or `WithoutBuiltins('fn1', 'fn2', ...)` to disable specific
function.

You can add custom funcs using `WithFuncs`.

### Custom cache

By default, this library uses a map to store all parsed templates, thus make them never expire. If you want expiration,
use `WithCache(impl)` to provide your own `Cache[*template.Template]` implementation.

### No cache

When cache is enabled (default), change to the template that has been parsed will not be visible until you rerun the
project (or the cache expired if you use custom cache implementation).

Use `WithNocache(true)` to disable template cache, force template to parse again on each execution.

## Pagination

This library provides a simple pagination implementation for using in template.
See [page](/page) package and the [example](/examples).