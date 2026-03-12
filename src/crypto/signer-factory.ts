/**
 * Unified signer factory that auto-detects key format.
 */

import { hexToBytes } from '../utils/buffer';
import { Signer } from './signer';
import { ECDSASigner } from './ecdsa';

/**
 * Create an {@link ECDSASigner} from a private key string.
 *
 * Accepts two formats:
 *
 * - **WIF** (Wallet Import Format): a Base58Check-encoded string
 *   (e.g. `L1aW4aubDFB7yfras2S1mN3bqg9nwySY8nkoLmJebSLD5BWv3ENZ`)
 * - **Hex**: a 64-character hex string with optional `0x` prefix
 *   (e.g. `0xabcdef0123456789...` or `abcdef0123456789...`)
 *
 * @throws Error if the key cannot be parsed in either format
 *
 * @example
 * ```ts
 * import { createSigner } from '@axlabs/neofs-sdk-ts-core';
 *
 * const signer = createSigner(process.env.NEOFS_PRIVATE_KEY!);
 * const client = new NeoFSClient({ endpoint, signer });
 * ```
 */
export function createSigner(privateKey: string): Signer {
  const trimmed = privateKey.trim();
  const hex = trimmed.toLowerCase().replace(/^0x/, '');

  // Try hex first (64 hex chars = 32 bytes)
  if (/^[0-9a-f]+$/i.test(hex) && hex.length === 64) {
    return ECDSASigner.fromPrivateKeyBytes(hexToBytes(hex));
  }

  // Try WIF
  try {
    return ECDSASigner.fromWIF(trimmed);
  } catch {
    // fall through
  }

  throw new Error(
    'Invalid private key: expected 64-char hex (with optional 0x prefix) or WIF (Base58Check).',
  );
}
