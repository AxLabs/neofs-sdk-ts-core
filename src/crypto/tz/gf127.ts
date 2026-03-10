/**
 * GF(2^127) - Galois Field with 2^127 elements
 * 
 * Elements are represented as two 64-bit values (as BigInts in JS).
 * The reduction polynomial is x^127 + x^63 + 1.
 */

const MSB64 = 1n << 63n;
const MAX_UINT64 = (1n << 64n) - 1n;

/**
 * GF127 represents an element of GF(2^127) as [lo, hi] where value = hi * x^64 + lo
 */
export class GF127 {
  lo: bigint;
  hi: bigint;

  constructor(lo: bigint = 0n, hi: bigint = 0n) {
    this.lo = lo & MAX_UINT64;
    this.hi = hi & MAX_UINT64;
  }

  clone(): GF127 {
    return new GF127(this.lo, this.hi);
  }

  /**
   * Set this to a + b (XOR in GF(2^127))
   */
  static add(a: GF127, b: GF127, c: GF127): void {
    c.lo = a.lo ^ b.lo;
    c.hi = a.hi ^ b.hi;
  }

  /**
   * Multiply by x (shift left by 1 with reduction)
   * Reduction polynomial: x^127 + x^63 + 1
   */
  static mul10(a: GF127, b: GF127): void {
    const carry = a.lo >> 63n;
    b.lo = (a.lo << 1n) & MAX_UINT64;
    b.hi = ((a.hi << 1n) ^ carry) & MAX_UINT64;

    // Reduction: if bit 127 is set, reduce by x^127 + x^63 + 1
    const mask = b.hi & MSB64;
    if (mask !== 0n) {
      b.lo ^= mask | (mask >> 63n);
      b.hi ^= mask;
    }
  }

  /**
   * Multiply by x + 1 (shift left by 1 and add original)
   */
  static mul11(a: GF127, b: GF127): void {
    const carry = a.lo >> 63n;
    b.lo = (a.lo ^ (a.lo << 1n)) & MAX_UINT64;
    b.hi = (a.hi ^ (a.hi << 1n) ^ carry) & MAX_UINT64;

    // Reduction
    const mask = b.hi & MSB64;
    if (mask !== 0n) {
      b.lo ^= mask | (mask >> 63n);
      b.hi ^= mask;
    }
  }

  /**
   * Copy a to b
   */
  static mul1(a: GF127, b: GF127): void {
    b.lo = a.lo;
    b.hi = a.hi;
  }

  /**
   * Full multiplication in GF(2^127)
   */
  static mul(a: GF127, b: GF127, c: GF127): void {
    const r = new GF127();
    const d = a.clone();
    
    // Process low 64 bits
    for (let i = 0n; i < 64n; i++) {
      if ((b.lo & (1n << i)) !== 0n) {
        GF127.add(r, d, r);
      }
      GF127.mul10(d, d);
    }
    
    // Process high 63 bits (bit 127 is always 0)
    for (let i = 0n; i < 63n; i++) {
      if ((b.hi & (1n << i)) !== 0n) {
        GF127.add(r, d, r);
      }
      GF127.mul10(d, d);
    }
    
    c.lo = r.lo;
    c.hi = r.hi;
  }

  /**
   * Convert to 16-byte array (big-endian)
   */
  toBytes(): Uint8Array {
    const buf = new Uint8Array(16);
    // High part first (big-endian)
    for (let i = 0; i < 8; i++) {
      buf[i] = Number((this.hi >> BigInt(56 - i * 8)) & 0xffn);
    }
    // Low part
    for (let i = 0; i < 8; i++) {
      buf[8 + i] = Number((this.lo >> BigInt(56 - i * 8)) & 0xffn);
    }
    return buf;
  }

  /**
   * Create from 16-byte array (big-endian)
   */
  static fromBytes(buf: Uint8Array): GF127 {
    if (buf.length !== 16) {
      throw new Error('GF127: data must be 16 bytes');
    }
    let hi = 0n;
    let lo = 0n;
    for (let i = 0; i < 8; i++) {
      hi = (hi << 8n) | BigInt(buf[i]);
    }
    for (let i = 0; i < 8; i++) {
      lo = (lo << 8n) | BigInt(buf[8 + i]);
    }
    return new GF127(lo, hi);
  }
}
