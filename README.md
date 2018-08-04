# ip matcher
ip matcher/ranger is a btree collection of multiple ip range in CIDR format (10.1.2.100/24) or range format (10.1.2.1-10.1.2.100), alse support port matcher/ranger.

## examples
```go
matcher := NewIPMatcher()

itemA, _ := StringToIPRange("86.100.32.0/24", "A")
matcher.Add(itemA)
matcher.AddRange(net.ParseIP("86.100.32.50"), net.ParseIP("86.100.32.100"), "B")

result := matcher.Match(net.PraseIP("86.100.32.60"))
fmt.Println(result) // sholud match itemA