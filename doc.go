// chris 061515

// Package steg implements a simple steganographic algorithm.  The
// algorithm is based on a puzzle given to the author by a friend.
//
// The Puzzle
//
// Consider an eight by eight bit matrix of all zeroes.
//
//	0 0 0 0 0 0 0 0
//	0 0 0 0 0 0 0 0
//	0 0 0 0 0 0 0 0
//	0 0 0 0 0 0 0 0
//	0 0 0 0 0 0 0 0
//	0 0 0 0 0 0 0 0
//	0 0 0 0 0 0 0 0
//	0 0 0 0 0 0 0 0
//
// Alice receives this bit matrix and must flip one bit before sending
// it to Bob.  How many bits can she transmit to him if they agree on a
// pre-arranged scheme?  Trivially, she can send one.  Select a
// particular bit.  If it's already the value Alice wants to send to
// Bob, she can flip any other bit.  If it's not, she can flip it to
// make sure he receives what she wants to send.  But can she do better?
//
// Sure!  She can flip any one of the sixty-four values, effectively
// transmitting six bits.
//
// Now suppose the bit matrix is randomized before receipt by Alice.
// For example, consider the following.
//
//	1 1 1 1 1 0 1 0
//	0 0 0 1 0 1 0 0
//	1 1 1 0 1 0 0 0
//	1 1 1 0 0 1 1 0
//	1 1 1 0 1 1 1 0
//	1 0 0 1 1 1 0 0
//	1 1 1 1 1 1 0 1
//	0 0 1 1 1 0 1 0
//
// Now how many bits can she send?  Again, trivially, she can send one,
// by the first strategy outlined.  But can she do better?
//
// To simplify the problem, let's consider a two by two bit matrix.
//
//	# #
//	# #
//
// The question becomes: is there a function of this matrix that yields
// two bits such that Alice can flip exactly one bit in it to ensure
// that Bob receives exactly the bits she intends?  Here is an approach.
//
//	    A
//	    ^
//	    |
//	B <-# #-
//	    # #
//	    |
//
// Let A be the XOR of the first and third bits, and let B be the XOR of
// the first and second bits.  If both A and B are the bits Alice wishes
// to send, then she can simply flip the fourth bit.  If one differs,
// she can simply flip the second or third bit, depending on which
// differs.  And if both differ, she can flip the first bit, flipping
// both A and B.
//
// Let's increase the size of the matrix to the next power of two to see
// if we can work towards addressing the original question about the
// eight by eight bit matrix.  Instead of thinking of a two by two
// matrix, let's think of a two by two by two cube.
//
//	  #---#
//	 /|  /|
//	#---# |
//	| # | #
//	|/  |/
//	#---#
//
// We might think to naively extend the previous solution, selecting
// three orthogonal lines of XORs.
//
//	     #---#
//	   |/|  /|
//	   #---# |
//	   | # | #
//	   |/  |/
//	A<-#---#-
//	  /|
//	 L v
//	B  C
//
// But what if two of the three differ?  What one bit can be flipped to
// alter the result appropriately?  A less naive approach is required.
// Instead of considering lines with our two by two by two solution,
// consider planes.
//
//	       A
//
//	       ^
//	       |
//	     #---#
//	    /|  /|
//	   #---# | -> B
//	   | # | #
//	   |/  |/
//	   #---#
//	  /
//	 L
//	C
//
// Let A, B, and C be the XORs of the vertices of three orthogonal sides
// of the cube.  If all three differ, then the common vertex can be
// flipped.  If only two out of the three differ, then the vertex common
// between the differing two, but uncommon to the third can be flipped.
// If only one differs, then the vertex uncommon to the other two can be
// flipped.  If the values of A, B, and C, are already what Alice wishes
// to send, then she can flip the vertex uncommon to the three planes as
// a "garbage" bit, leaving A, B, and C unchanged.
//
// So, this approach suggests that given 2^n random bits, we can in fact
// encode n specific bits within them by flipping exactly one bit.
// Therefore, the answer to the puzzle is that Alice can do better!  She
// can encode all 6 bits, even if the eight by eight bit matrix she
// receives is entirely random.
//
// An Algorithm
//
// The general idea of the above solution to the puzzle is that given
// 2^n carrier bits, select n sets of bits, each of size 2^(n-1), or
// half.  But how to select the sets, in general?  Bob must be able to
// read n message bits.  Consider the message bit at index 0.  XOR
// together all of the bits in the carrier whose index value itself has
// a 1 at its index 0.  And so on.  For example, to read three message
// bits A, B, and C, from eight carrier bits:
//
//	0 0 0
//	0 0 1     a
//	0 1 0   b
//	0 1 1   b a
//	1 0 0 c
//	1 0 1 c   a
//	1 1 0 c b
//	1 1 1 c b a
//	      | | |
//	      | | +-> A
//	      | +---> B
//	      +-----> C
//
// Compute A by XORing the 'a's; B by XOR the 'b's, etc.
//
// But how does Alice determine which bit to flip?  Simple: first, read
// the message bits out of the carrier bits she's given.  Then, XOR the
// message bits in the carrier with the message bits she wishes to send.
// There will be 1s where the bits differ.  This gives the exact address
// of the bit to flip that will alter the message bits that Bob will
// read out of the carrier such that it will match Alice's desired
// message.  Why is this so?  Suppose the message bits Alice reads
// differ from what she wants to send at a particular index.  If she
// flips a bit in the carrier that has a 1 in that index position, then
// the resulting message bits that Bob reads will have that bit, and
// only that bit flipped, by the reading scheme outlined earlier.
//
// This is the algorithm that package steg implements.
//
// Implementation
//
// Since thinking about sending message bits is unconventional at best,
// package steg takes a byte-based approach.  It defines an atom, the
// number of message bytes that are read or written at a time.
// Correspondingly, this implicitly defines a chunk, the number of
// carrier bytes used to embed an atom.  The following atom sizes and
// implicit chunk sizes are available.
//
//	atom size   chunk size
//	   1B           32B
//	   2B          8KiB
//	   3B          2MiB
//
// First, create a context with an atom size.  Then, create readers,
// writers, and muxes from the context.  By design, the implementation
// makes no effort to be aware of the character of the carrier data.
//
// References
//
// https://en.wikipedia.org/wiki/Steganography
package steg
