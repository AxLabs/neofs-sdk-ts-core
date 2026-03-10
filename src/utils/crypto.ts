/**
 * React Native crypto utilities
 * This SDK is designed specifically for React Native
 */

/**
 * SHA256 hash using crypto-js
 */
export function sha256(data: Uint8Array | string): Uint8Array {
  // Always use crypto-js for React Native compatibility
  const CryptoJS = require('crypto-js');
  let hash;
  if (typeof data === 'string') {
    hash = CryptoJS.SHA256(data);
  } else {
    const wordArray = CryptoJS.lib.WordArray.create(data);
    hash = CryptoJS.SHA256(wordArray);
  }
  
  // Convert CryptoJS WordArray to Uint8Array
  const bytes = [];
  for (let i = 0; i < hash.sigBytes; i++) {
    bytes.push((hash.words[Math.floor(i / 4)] >>> (24 - (i % 4) * 8)) & 0xff);
  }
  return new Uint8Array(bytes);
}

/**
 * SHA512 hash using crypto-js
 */
export function sha512(data: Uint8Array | string): Uint8Array {
  // Always use crypto-js for React Native compatibility
  const CryptoJS = require('crypto-js');
  let hash;
  if (typeof data === 'string') {
    hash = CryptoJS.SHA512(data);
  } else {
    const wordArray = CryptoJS.lib.WordArray.create(data);
    hash = CryptoJS.SHA512(wordArray);
  }
  
  // Convert CryptoJS WordArray to Uint8Array
  const bytes = [];
  for (let i = 0; i < hash.sigBytes; i++) {
    bytes.push((hash.words[Math.floor(i / 4)] >>> (24 - (i % 4) * 8)) & 0xff);
  }
  return new Uint8Array(bytes);
}

/**
 * RIPEMD160 hash using crypto-js
 */
export function ripemd160(data: Uint8Array): Uint8Array {
  const CryptoJS = require('crypto-js');
  const wordArray = CryptoJS.lib.WordArray.create(data);
  const hash = CryptoJS.RIPEMD160(wordArray);
  
  // Convert CryptoJS WordArray to Uint8Array
  const bytes = [];
  for (let i = 0; i < hash.sigBytes; i++) {
    bytes.push((hash.words[Math.floor(i / 4)] >>> (24 - (i % 4) * 8)) & 0xff);
  }
  return new Uint8Array(bytes);
}

/**
 * Double SHA256 hash
 */
export function doubleSha256(data: Uint8Array): Uint8Array {
  return sha256(sha256(data));
}
