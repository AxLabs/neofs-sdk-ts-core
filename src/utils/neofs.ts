/**
 * NeoFS-specific utility helpers.
 */

import { ContainerID, ObjectAttribute } from '../types';
import { hexToBytes } from './buffer';
import { base58Decode } from './base58';

/**
 * Parse a container ID from a hex string (with optional `0x` prefix) or
 * a Base58-encoded string into a {@link ContainerID}.
 *
 * @example
 * ```ts
 * const cid = toContainerId('CYzn4HPHXRmGW5cwvrqHtosd79kjBtVYAP3ykH6RCCAa');
 * const cid2 = toContainerId('0xabcdef...');
 * ```
 *
 * @throws Error if the string is not valid hex or Base58
 */
export function toContainerId(containerId: string): ContainerID {
  const raw = containerId.trim();
  const hex = raw.toLowerCase().replace(/^0x/, '');
  if (/^[0-9a-f]+$/i.test(hex) && hex.length > 0 && hex.length % 2 === 0) {
    return { value: hexToBytes(hex) };
  }
  try {
    const bytes = base58Decode(raw);
    return { value: bytes };
  } catch {
    throw new Error(
      `Invalid container ID: expected hex (even length, optional 0x prefix) or Base58. Got: "${raw}"`,
    );
  }
}

/**
 * Convert an array of `{ key, value }` object attributes (as returned by the
 * NeoFS API) into a plain `Record<string, string>` for easier lookup.
 *
 * @example
 * ```ts
 * const head = await objectClient.head({ address, raw: false });
 * const attrs = decodeAttributes(head.attributes);
 * console.log(attrs['FileName']);
 * ```
 */
export function decodeAttributes(
  attrs?: ObjectAttribute[] | Array<{ key: string; value: string }>,
): Record<string, string> {
  const out: Record<string, string> = {};
  if (attrs) {
    for (const a of attrs) out[a.key] = a.value;
  }
  return out;
}

/**
 * Classify a NeoFS gRPC / SDK error as retryable.
 *
 * Returns `true` for transient failures such as expired sessions,
 * authentication issues that may self-resolve after session renewal,
 * and timeout / deadline errors.
 */
export function isRetryableNeoFSError(err: unknown): boolean {
  const msg = String((err as any)?.message || (err as any)?.details || '');
  return (
    msg.includes('session') ||
    msg.includes('expired') ||
    msg.includes('UNAUTHENTICATED') ||
    msg.includes('permission') ||
    msg.includes('deadline') ||
    msg.includes('timeout')
  );
}
