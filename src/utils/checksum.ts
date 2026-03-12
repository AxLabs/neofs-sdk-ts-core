/**
 * Payload checksum helpers for NeoFS object headers.
 */

import { ChecksumType } from '../types/enums';
import { sha256 } from './crypto';
import { tzHash } from '../crypto/tz';

/**
 * Compute both SHA-256 and Tillich-Zémor (TZ) checksums for an object payload.
 *
 * NeoFS object headers require `payloadHash` (SHA-256) and `homomorphicHash` (TZ).
 * This helper produces both in the format expected by the `ObjectClient.put()` header.
 *
 * @example
 * ```ts
 * const checksums = payloadChecksums(payload);
 * await objectClient.put({
 *   header: { containerId, ownerId, ...checksums },
 *   payload,
 * });
 * ```
 */
export function payloadChecksums(payload: Uint8Array): {
  payloadHash: { type: number; sum: Uint8Array };
  homomorphicHash: { type: number; sum: Uint8Array };
} {
  return {
    payloadHash: {
      type: ChecksumType.SHA256,
      sum: sha256(payload),
    },
    homomorphicHash: {
      type: ChecksumType.TZ,
      sum: tzHash(payload),
    },
  };
}
