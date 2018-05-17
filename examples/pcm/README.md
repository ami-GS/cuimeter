# Processor Counter Monitor (PCM) visualization

This example visualize output from [pcm](https://github.com/opcm/pcm). This means the program reads pcm output from pipe, not as target.

## Requirements
- [pcm](https://github.com/opcm/pcm)


## Usage

You can ignore to use --target directive as follow (you can if you want and implement to use it by yourself).
```sh
$PATH_TO_PCM/pcm-memory.x | go run memory.go
```


## Warning

Because pcm doesn't seem bufferring output, this program has to read from beginning to the end of one piped data.
This means that user have to know buffer size which io.reader have to read at a time for one piped data.
