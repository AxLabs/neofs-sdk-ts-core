/**
 * NeoFS User ID (Owner ID) utilities.
 * 
 * The owner ID is a 25-byte value derived from a public key:
 * - 1 byte: NEO address version prefix (0x35 for Neo N3)
 * - 20 bytes: script hash (RIPEMD160(SHA256(verification_script)))
 * - 4 bytes: checksum (first 4 bytes of double SHA256 of prefix + script hash)
 */

import { sha256, ripemd160, doubleSha256 } from '../utils/crypto';

// NEO N3 address version prefix
const NEO_ADDRESS_PREFIX = 0x35;

// Neo N3 CheckSig interop call hash (little-endian)
// "Neo.Crypto.CheckSig" => 0x27b3e756 (big-endian) => bytes: [0x56, 0xe7, 0xb3, 0x27]
const CHECKSIG_INTEROP = new Uint8Array([0x56, 0xe7, 0xb3, 0x27]);

/**
 * Creates a verification script from a compressed public key.
 * The script format is: PUSHDATA1 + length(33) + pubkey(33) + SYSCALL + interopHash(4)
 * 
 * Opcodes:
 * - 0x0C = PUSHDATA1 (push data with 1-byte length prefix)
 * - 0x21 = 33 (length of compressed public key)
 * - 0x41 = SYSCALL
 */
function createVerificationScript(compressedPublicKey: Uint8Array): Uint8Array {
  if (compressedPublicKey.length !== 33) {
    throw new Error(`Expected 33-byte compressed public key, got ${compressedPublicKey.length}`);
  }

  // PUSHDATA1 (0x0C) + length (0x21 = 33) + pubkey (33 bytes) + SYSCALL (0x41) + interop hash (4 bytes)
  const script = new Uint8Array(1 + 1 + 33 + 1 + 4);
  script[0] = 0x0c; // PUSHDATA1
  script[1] = 0x21; // 33 bytes length
  script.set(compressedPublicKey, 2);
  script[35] = 0x41; // SYSCALL
  script.set(CHECKSIG_INTEROP, 36);
  
  return script;
}

/**
 * Calculates the script hash (20 bytes) from a verification script.
 * Script hash = RIPEMD160(SHA256(script))
 */
function getScriptHash(script: Uint8Array): Uint8Array {
  const sha256Hash = sha256(script);
  return ripemd160(sha256Hash);
}

/**
 * Creates a NeoFS Owner ID (User ID) from a compressed public key.
 * 
 * @param compressedPublicKey - The 33-byte compressed ECDSA public key
 * @returns The 25-byte owner ID
 */
export function ownerIdFromPublicKey(compressedPublicKey: Uint8Array): Uint8Array {
  // Create verification script
  const script = createVerificationScript(compressedPublicKey);
  
  // Get script hash
  const scriptHash = getScriptHash(script);
  
  // Build owner ID: prefix + script hash + checksum
  const ownerId = new Uint8Array(25);
  ownerId[0] = NEO_ADDRESS_PREFIX;
  ownerId.set(scriptHash, 1);
  
  // Calculate checksum: first 4 bytes of double SHA256
  const checksumInput = ownerId.slice(0, 21);
  const checksum = doubleSha256(checksumInput).slice(0, 4);
  ownerId.set(checksum, 21);
  
  return ownerId;
}

/**
 * Validates an owner ID.
 * 
 * @param ownerId - The 25-byte owner ID to validate
 * @returns true if valid, false otherwise
 */
export function validateOwnerId(ownerId: Uint8Array): boolean {
  if (ownerId.length !== 25) {
    return false;
  }
  
  // Check prefix
  if (ownerId[0] !== NEO_ADDRESS_PREFIX) {
    return false;
  }
  
  // Verify checksum
  const checksumInput = ownerId.slice(0, 21);
  const expectedChecksum = doubleSha256(checksumInput).slice(0, 4);
  const actualChecksum = ownerId.slice(21, 25);
  
  for (let i = 0; i < 4; i++) {
    if (expectedChecksum[i] !== actualChecksum[i]) {
      return false;
    }
  }
  
  return true;
}
