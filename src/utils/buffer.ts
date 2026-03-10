/**
 * Cross-platform buffer utilities
 */

// Type declarations for browser globals
declare const window: any;
declare const btoa: (str: string) => string;
declare const atob: (str: string) => string;

/**
 * Converts hex string to Uint8Array
 */
export function hexToBytes(hex: string): Uint8Array {
  const bytes = new Uint8Array(hex.length / 2);
  for (let i = 0; i < hex.length; i += 2) {
    bytes[i / 2] = parseInt(hex.substr(i, 2), 16);
  }
  return bytes;
}

/**
 * Converts Uint8Array to hex string
 */
export function bytesToHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}

/**
 * Creates a Uint8Array from hex string (cross-platform Buffer.from alternative)
 */
export function fromHex(hex: string): Uint8Array {
  return hexToBytes(hex);
}

/**
 * Creates a Uint8Array from base64 string (cross-platform Buffer.from alternative)
 */
export function fromBase64(base64: string): Uint8Array {
  // Simple base64 decode implementation
  const binaryString = atob(base64);
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes;
}

/**
 * Converts Uint8Array to base64 string (cross-platform Buffer.toString alternative)
 */
export function toBase64(bytes: Uint8Array): string {
  let binary = '';
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

/**
 * Cross-platform random bytes generation
 */
export function randomBytes(length: number): Uint8Array {
  if (typeof window !== 'undefined' && window.crypto && window.crypto.getRandomValues) {
    // Browser/React Native
    return window.crypto.getRandomValues(new Uint8Array(length));
  } else if (typeof global !== 'undefined' && global.crypto && global.crypto.getRandomValues) {
    // React Native global crypto
    return global.crypto.getRandomValues(new Uint8Array(length));
  } else {
    // Fallback: use Math.random (not cryptographically secure, but works)
    const bytes = new Uint8Array(length);
    for (let i = 0; i < length; i++) {
      bytes[i] = Math.floor(Math.random() * 256);
    }
    return bytes;
  }
}

/**
 * Cross-platform UUID v4 generation
 */
export function generateUUID(): string {
  if (typeof window !== 'undefined' && window.crypto && window.crypto.randomUUID) {
    // Modern browsers
    return window.crypto.randomUUID();
  } else if (typeof global !== 'undefined' && global.crypto && global.crypto.randomUUID) {
    // React Native global crypto
    return global.crypto.randomUUID();
  } else {
    // Fallback implementation
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
      const r = Math.random() * 16 | 0;
      const v = c === 'x' ? r : (r & 0x3 | 0x8);
      return v.toString(16);
    });
  }
}
