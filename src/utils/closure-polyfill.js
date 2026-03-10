/**
 * Closure Library polyfill for React Native
 * This provides the missing goog.object functions that protobuf files expect
 */

// Create global goog object if it doesn't exist
if (typeof global !== 'undefined' && !global.goog) {
  global.goog = {};
}

if (typeof window !== 'undefined' && !window.goog) {
  window.goog = {};
}

// Create goog.object if it doesn't exist
const goog = typeof global !== 'undefined' ? global.goog : window.goog;

if (!goog.object) {
  goog.object = {};
}

// Create global proto object that protobuf files expect
if (typeof global !== 'undefined' && !global.proto) {
  global.proto = {
    neo: {
      fs: {
        v2: {
          accounting: {},
          container: {},
          object: {},
          session: {},
          netmap: {},
          reputation: {}
        }
      }
    }
  };
  console.log('Created global.proto with structure');
}

if (typeof window !== 'undefined' && !window.proto) {
  window.proto = {
    neo: {
      fs: {
        v2: {
          accounting: {},
          container: {},
          object: {},
          session: {},
          netmap: {},
          reputation: {}
        }
      }
    }
  };
  console.log('Created window.proto with structure');
}

// Ensure proto is available everywhere
if (typeof global !== 'undefined') {
  global.proto = global.proto || {
    neo: {
      fs: {
        v2: {
          accounting: {},
          container: {},
          object: {},
          session: {},
          netmap: {},
          reputation: {}
        }
      }
    }
  };
}
if (typeof window !== 'undefined') {
  window.proto = window.proto || {
    neo: {
      fs: {
        v2: {
          accounting: {},
          container: {},
          object: {},
          session: {},
          netmap: {},
          reputation: {}
        }
      }
    }
  };
}

// Polyfill for goog.object.extend
if (!goog.object.extend) {
  goog.object.extend = function(target, source) {
    for (const key in source) {
      if (source.hasOwnProperty(key)) {
        target[key] = source[key];
      }
    }
    return target;
  };
}

// Polyfill for goog.object.setIfUndefined
if (!goog.object.setIfUndefined) {
  goog.object.setIfUndefined = function(obj, key, value) {
    if (!(key in obj)) {
      obj[key] = value;
    }
    return obj[key];
  };
}

// Polyfill for goog.object.clone
if (!goog.object.clone) {
  goog.object.clone = function(obj) {
    const clone = {};
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clone[key] = obj[key];
      }
    }
    return clone;
  };
}

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
  module.exports = goog;
}
