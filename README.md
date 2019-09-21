## sixb [![Go Report Card](https://goreportcard.com/badge/github.com/jfcg/sixb)](https://goreportcard.com/report/github.com/jfcg/sixb) [![GoDoc](https://godoc.org/github.com/jfcg/sixb?status.svg)](https://godoc.org/github.com/jfcg/sixb)
Some string utility functions including slice conversions and a fast hash.

### Txt2int collisions
We have searched for collisions on Txt2int 64-bit hash function for more than 5x2^37 **short utf8 text** inputs but could not identify any. Collisions were found for longer inputs and it is not a cryptographic hash.

If you can find a collision for short\(<9 bytes\) inputs, please do let me know. Also I want to find out the "minimum sum of colliding input lengths" len\(s1\)+len\(s2\) where:
- s1, s2 are distinct utf8 texts
- Txt2int\(s1\) = Txt2int\(s2\)

### Contributors
- Chris Burkert \(burkert.chrisATgmail\) for searching collisions on his 6 TiB ram 240 cores server
- Michael T. Jones \(michael.jonesATgmail\) for his valuable suggestions on hash design
