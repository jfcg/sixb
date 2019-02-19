## sixb
Some string utility functions

### Txt2int collisions
We have searched for collisions on Txt2int 64-bit hash function for more than 5x2^37 **short utf8 text** inputs but could not identify any. If you can, please do let me know ;)

In general I want to find out the "minimum sum of colliding input lengths" len\(s1\)+len\(s2\) where:
- s1, s2 are distinct utf8 texts
- Txt2int\(s1\) = Txt2int\(s2\)

### Contributors
- Chris Burkert \(burkert.chrisATgmail\) for searching collisions on his 6 TiB ram 240 cores server
- Michael T. Jones \(michael.jonesATgmail\) for his valuable suggestions on hash design
