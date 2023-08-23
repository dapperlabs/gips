## gnatips-go

A tiny convenience package to get NAT ips from a region of a particular GCP project.

I wanted a nice very simple package to do this and didn't see one.

```go
ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()
ips, err := gnatips.Get(ctx, "project-name-goes-here", "us-west2")
if err != nil {
    // Do something.
}
// ips: ["127.0.0.1", "8.8.8.8", "1.1.1.1"]
```