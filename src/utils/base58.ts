/**
 * Base58 encoding/decoding (Bitcoin alphabet).
 * Extracted from ECDSASigner so it can be shared across the SDK.
 */

const ALPHABET = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz';

/**
 * Decodes a Base58-encoded string to bytes.
 *
 * @throws Error if the string contains invalid Base58 characters
 */
export function base58Decode(encoded: string): Uint8Array {
  let value = 0n;

  for (let i = 0; i < encoded.length; i++) {
    const digit = ALPHABET.indexOf(encoded[i]);
    if (digit < 0) {
      throw new Error(`Invalid Base58 character '${encoded[i]}' at position ${i}`);
    }
    value = value * 58n + BigInt(digit);
  }

  const leadingZeros = encoded.match(/^1+/)?.[0].length || 0;
  const bytes: number[] = [];

  while (value > 0n) {
    bytes.unshift(Number(value & 0xFFn));
    value = value >> 8n;
  }

  const result = new Uint8Array(leadingZeros + bytes.length);
  result.set(bytes, leadingZeros);

  return result;
}

/**
 * Encodes bytes to a Base58 string.
 */
export function base58Encode(bytes: Uint8Array): string {
  let value = 0n;
  for (const b of bytes) {
    value = value * 256n + BigInt(b);
  }

  let encoded = '';
  while (value > 0n) {
    encoded = ALPHABET[Number(value % 58n)] + encoded;
    value = value / 58n;
  }

  for (const b of bytes) {
    if (b !== 0) break;
    encoded = '1' + encoded;
  }

  return encoded || '1';
}
