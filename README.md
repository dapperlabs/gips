## gips

`gips` = "Google IPs"

Finding an external IP inside of Google Cloud is frustrating - they're listed in a bunch of different places and you end up having to do a lot of work to find a particular one.

This little application finds all of the external IPs and throws them into an in-memory map.

NOTE: If you had a ton of projects and memory is becoming an issue - adding an additional adaptor would be pretty simple - `core.ProjectService` is the interface you'll need to conform to.

It's queryable with `curl` and outputs JSON:

```bash
/bin $ curl http://gips:8080/api/v1/project/project-name-goes-here
{
  "name": "project-name-goes-here",
  "regions": [
    {
      "region": "us-west1",
      "ips": [
        "34.2.123.129",
        "35.3.124.141",
        "34.4.125.109",
        "34.5.126.81",
        "34.6.127.84",
        "34.7.128.34"
      ]
    }
  ]
}
/bin $ curl http://gips:8080/api/v1/search/34.2.123.129
{"name":"project-name-goes-here","regions":[{"region":"us-west1","ips":["34.2.123.129","35.3.124.141","34.4.125.109","34.5.126.81","34.6.127.84","34.7.128.34"]}]}
```

## TODO

- [x] NAT IPs
- [x] golang package to query this - `search`
- [x] cli interface for search endpoint and to test `search` package
- [ ] Figure out proper amount of permissions this needs to run - add Terraform to allow those permissions
- [ ] CloudSQL IPs?
- [ ] Do we tag the IPs so we can find them quickly? OR is just knowing the project enough for this?
- [ ] Helm chart?