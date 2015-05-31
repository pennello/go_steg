Helps to understand how to construct masks.

    func tobin(x uint64) string {
        s := strconv.FormatUint(x, 2)
        s = strings.Repeat("0", 8 - len(s)) + s
        return s[0:4] + " " + s[4:]
    }
    
    func printRanges() {
        for byte_index := uint64(0); byte_index < 32; byte_index += 1 {
            low_int := byte_index * 8
            high_int := low_int + 8 - 1
            low_str := tobin(low_int)
            high_str := tobin(high_int)
            fmt.Printf("%02d %s -> %s\n", byte_index, low_str, high_str)
        }
    }

Table used to construct original 1-byte atom/32-byte chunk masks.

    0 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
    1 cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc
    2 f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0
    3 00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff
    4 0000ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000ffff
    5 00000000ffffffff00000000ffffffff00000000ffffffff00000000ffffffff
    6 0000000000000000ffffffffffffffff0000000000000000ffffffffffffffff
    7 00000000000000000000000000000000ffffffffffffffffffffffffffffffff

Table of atom/chunk size trade-offs.

      atom size                 chunk size
    bytes    bits           bits          bytes
    1B         8b        256   b         32   B
    2B        16b         64 Kib          8 KiB
    3B        24b         16 Mib          2 MiB
    4B        32b          4 Gib        512 MiB
    5B        40b          1 Tib        128 GiB
    6B        48b        256 Tib         32 TiB
    7B        56b         64 Pib          8 PiB
    8B        64b         16 Eib          2 EiB

Naive implementation of read, doing one bit at a time.

    func (ctx *Ctx) read(c []byte) []byte {
        rlen := uint(ctx.atomSize)
        r := make([]byte, rlen)
        // cBi:  chunk  byte index
        // cbsi: chunk  bit  sub-index
        // cbi:  chunk  bit  index
        // rBi:  return byte index
        // rbsi: return bit  sub-index
        // rbi:  return bit  index
        for cBi := uint(0); cBi < ctx.chunkSize; cBi++ {
            B := c[cBi]
            for cbsi := uint(0); cbsi < 8; cbsi++ {
                cbi := cBi<<3 | cbsi
                for rBi := uint(0); rBi < rlen; rBi++ {
                    for rbsi := uint(0); rbsi < 8; rbsi++ {
                        rbi := rBi<<3 | rbsi
                        if (1<<rbi)&cbi != 0 {
                            b := B & (1 << cbsi)
                            b >>= cbsi
                            b <<= rbsi
                            r[rBi] ^= b
                        }
                    }
                }
            }
        }
        return r
    }
