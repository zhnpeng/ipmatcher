# ip matcher
ip matcher/ranger is a btree collection of multiple ip range in CIDR format (10.1.2.100/24) or range format (10.1.2.1-10.1.2.100), alse support port matcher/ranger.

## examples
```go
portMatcherA := NewPortMatcher()
portr, _ := StringToPortRange("0-500", "A")
portMatcherA.Add(portr)

ipMatcher := NewIPMatcher()
item, _ := StringToIPRange("86.100.32.0/24", "A")
ipMatcher.Add(item)
ipMatcher.AddRange(net.ParseIP("86.100.32.50"), net.ParseIP("86.100.32.100"), "B")

result := ipMatcher.Match(net.ParseIP("86.100.32.1"))
fmt.Println(result) // should match ip range with portmMatcherA