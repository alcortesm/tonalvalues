# tonalvalues

tonalvalues is a command line tool that analyses the tonal values of a jpg file
and present them visually as another photo.
The goal is to help art students to understand tonal values.

## Usage

The command expects one argument, the jpg file to analyse.
It creates an `output.jpg` file showing the tonal values of the original photo,
using different amounts of gay tones.

For example...

```
; tonalvalues example.jpg
```

... generates an `output.jpg` file looking like this:

![example of the output](./output_explained.jpg)
