# Version 2 Protocol Thoughts
## One location Three Items
Version 2 of the protocol can consider having location as the top level encode
and then dangle off \*Foreground \*Background and Ascii value

Then the window could just track the diffs and send down minimum amounts of
information for the change, thus we only have to specify the 2 byte location
once for ascii + color change.

## Color
Encode color as 2 bytes

 1  2-6        7-11       12-16
+--+----------+----------+----------+
|F | Red      | Green    | Blue     |
+--+----------+----------+----------+

33% reduction in color space requirements

## Color Again
We could do color compression but only have a specified color table ahead of
time that encodes up to 128 colors as 1 byte values, with one bit reserved for
F/B

66% savings on the most frequent

