/**
 * Tillich-Zémor homomorphic hash implementation.
 * 
 * This is a cryptographic hash function based on matrix multiplication
 * over the Galois field GF(2^127). It has the homomorphic property:
 * H(m1 || m2) = H(m1) * H(m2)
 * 
 * Reference: https://link.springer.com/content/pdf/10.1007/3-540-48658-5_5.pdf
 */

import { GF127 } from './gf127';

/**
 * Size of the Tillich-Zémor hash output in bytes.
 */
export const TZ_SIZE = 64;

/**
 * Internal digest state - a 2x2 matrix over GF(2^127)
 * Stored as [c00, c10, c01, c11] for memory layout compatibility with Go
 */
class Digest {
  // Matrix stored as:
  // [ x[0] x[2] ]
  // [ x[1] x[3] ]
  x: [GF127, GF127, GF127, GF127];

  constructor() {
    this.x = [
      new GF127(1n, 0n),  // c00 = 1
      new GF127(0n, 0n),  // c10 = 0
      new GF127(0n, 0n),  // c01 = 0
      new GF127(1n, 0n),  // c11 = 1
    ];
  }

  /**
   * Reset to identity matrix
   */
  reset(): void {
    this.x[0] = new GF127(1n, 0n);
    this.x[1] = new GF127(0n, 0n);
    this.x[2] = new GF127(0n, 0n);
    this.x[3] = new GF127(1n, 0n);
  }

  /**
   * Process a single bit by multiplying with the appropriate generator matrix.
   * 
   * For bit 0: multiply by A = [[x, 1], [1, 0]]
   * For bit 1: multiply by B = [[x, x+1], [1, 1]]
   */
  private mulBitRight(bit: boolean, tmp: GF127): void {
    const c00 = this.x[0];
    const c10 = this.x[1];
    const c01 = this.x[2];
    const c11 = this.x[3];

    if (bit) {
      // Multiply by B = [[x, x+1], [1, 1]]
      // c00' = c00*x + c01
      // c01' = c00*(x+1) + c01 = c00*x + c00 + c01
      // c10' = c10*x + c11
      // c11' = c10*(x+1) + c11 = c10*x + c10 + c11
      
      GF127.mul1(c00, tmp);
      GF127.mul10(c00, c00);
      GF127.add(c00, c01, c00);
      GF127.mul11(tmp, tmp);
      GF127.add(c01, tmp, c01);

      GF127.mul1(c10, tmp);
      GF127.mul10(c10, c10);
      GF127.add(c10, c11, c10);
      GF127.mul11(tmp, tmp);
      GF127.add(c11, tmp, c11);
    } else {
      // Multiply by A = [[x, 1], [1, 0]]
      // c00' = c00*x + c01
      // c01' = c00
      // c10' = c10*x + c11
      // c11' = c10
      
      GF127.mul1(c00, tmp);
      GF127.mul10(c00, c00);
      GF127.add(c00, c01, c00);
      GF127.mul1(tmp, c01);

      GF127.mul1(c10, tmp);
      GF127.mul10(c10, c10);
      GF127.add(c10, c11, c10);
      GF127.mul1(tmp, c11);
    }
  }

  /**
   * Write data to the hash
   */
  write(data: Uint8Array): void {
    const tmp = new GF127();
    for (const byte of data) {
      // Process each bit of the byte, MSB first
      this.mulBitRight((byte & 0x80) !== 0, tmp);
      this.mulBitRight((byte & 0x40) !== 0, tmp);
      this.mulBitRight((byte & 0x20) !== 0, tmp);
      this.mulBitRight((byte & 0x10) !== 0, tmp);
      this.mulBitRight((byte & 0x08) !== 0, tmp);
      this.mulBitRight((byte & 0x04) !== 0, tmp);
      this.mulBitRight((byte & 0x02) !== 0, tmp);
      this.mulBitRight((byte & 0x01) !== 0, tmp);
    }
  }

  /**
   * Get the final hash value (64 bytes)
   */
  sum(): Uint8Array {
    const result = new Uint8Array(TZ_SIZE);
    
    // Output order: c00, c01, c10, c11 (row-major)
    const b0 = this.x[0].toBytes();
    const b1 = this.x[2].toBytes();
    const b2 = this.x[1].toBytes();
    const b3 = this.x[3].toBytes();
    
    result.set(b0, 0);
    result.set(b1, 16);
    result.set(b2, 32);
    result.set(b3, 48);
    
    return result;
  }
}

/**
 * Compute the Tillich-Zémor hash of data.
 * 
 * @param data - Input data to hash
 * @returns 64-byte hash
 */
export function tzHash(data: Uint8Array): Uint8Array {
  const d = new Digest();
  d.write(data);
  return d.sum();
}

/**
 * Create a new Tillich-Zémor hash instance for incremental hashing.
 */
export function createTzHash(): {
  write: (data: Uint8Array) => void;
  sum: () => Uint8Array;
  reset: () => void;
} {
  const d = new Digest();
  return {
    write: (data: Uint8Array) => d.write(data),
    sum: () => d.sum(),
    reset: () => d.reset(),
  };
}
