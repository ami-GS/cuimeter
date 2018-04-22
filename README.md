# cuimeter
meter visualization framework on cui

### Introduction
cuimeter enables visualization like [speedometer](http://excess.org/speedometer/) for any telemetry.

The motivation is to compare telemetry with several settings on one cui screen.

What user have to do is just to implement interfaces defined bellow, then this framework should run with your definition.

Please refere [examples](https://github.com/ami-GS/cuimeter/tree/master/examples) for easy understandings.

### Example

```sh
>> git clone https://github.com/ami-GS/cuimeter
>> cd cuimeter/example
>> go run oneLine.go --target ./target1.txt --target ./target2.txt
```

### Interfaces
Currently you have to implement interfaces bellows
- `Parse(data string) (map[string]int64, error)`
  - parse data following your telemetry format
- `Get(retData \*int64, wg \*sync.WaitGroup)`
  - How to get Data. (e.g. file open, http get and so on)
- `GetUnit() string`
  - already defined, but you can override
- `GetInterval() time.Duration`
  - already defined, but you can override
