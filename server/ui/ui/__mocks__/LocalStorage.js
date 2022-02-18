class LocalStorageMock {
  constructor() {
    const user = JSON.stringify({ id_token: 'id_token' })
    this.store = {
      'oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer': user
    };
  }

  clear() {
    this.store = {};
  }

  getItem(key) {
    return this.store[key] || null;
  }

  setItem(key, value) {
    this.store[key] = String(value);
  }

  removeItem(key) {
    delete this.store[key];
  }
}

// Object.defineProperty(window, 'localStorage', {
//      value: LocalStorageMock
// });

global.window.localStorage = new LocalStorageMock();

global.localStorage = new LocalStorageMock();
