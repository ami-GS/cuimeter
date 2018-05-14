# cuimeter
meter visualization framework on cui

## Introduction
cuimeter enables visualization like [speedometer](http://excess.org/speedometer/) for any telemetry. The motivation is to compare telemetry with several settings on one cui screen. What user have to do is just to implement interfaces defined bellow, then this framework should run with your definition.

Please refere [examples](https://github.com/ami-GS/cuimeter/tree/master/examples) for easy understandings.

## Requirements
- go 1.10+

## Example

```sh
>> git clone https://github.com/ami-GS/cuimeter
>> cd cuimeter/example
>> go run oneLine.go --target ./target1.txt --target ./target2.txt
```

Output should be bellow
![alt text](https://user-images.githubusercontent.com/5763034/39187812-33d8839e-4809-11e8-8f6d-bc68bb162872.png)


## Interfaces
Developer have to implement interfaces bellows
- `func (self *YourHint) read() (string, error)`
  - How to get Data. (e.g. file open, http get and so on)
- `func (self *YourHint) parse(string) (int64, error)`
  - parse data following your telemetry format
- `func (self *YourHint) postProcess(int64) int64`
  - called after `parse` to process data, then data is ingested to graph
- `func (self *YourHint) getUnit() string`
  - already defined, but you can override
- `func (self *YourHint) getInterval() time.Duration`
  - already defined, but you can override

## Flags
Flags user can set are as bellows
- --target
  - Indicates which file/url data will be read, at least one target need to be set

## TODO
- Flexible number of targets (up to 2 targets as of now)
- Add command line argument template
  - interval, label and so on.
- Track Min as well (currently only max)
