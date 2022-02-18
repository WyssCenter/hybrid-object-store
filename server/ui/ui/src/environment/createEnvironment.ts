//
const baseUrl = `${window.location.protocol}//${window.location.hostname}/core/v1`;
let authUrl = `${window.location.protocol}//${window.location.hostname}/auth/v1`;
const hostname = `${window.location.protocol}//${window.location.hostname}`;

const baseHeaders = {
  "Access-Control-Allow-Origin": "*",
  "Content-Type": 'application/json',
  "Origin": window.location.origin
}


fetch(`${baseUrl}/discover`, { headers: baseHeaders, method: 'GET' })
.then(response => response.json())
.then(data => authUrl = data.auth_service);


const getIdToken = () => {
  const bearerToken = JSON.parse(localStorage.getItem(`oidc.user:${authUrl}/.well-known/openid-configuration:HossServer`)).id_token;
  return bearerToken;
}
/**
* Method uses delete method with api
* @param {string} route
* @return {Promise}
*/
const del = (
  route: string,
  body = {},
  isAuthRoute = false,
): Promise<Response> => {
  const stringifiedBody = JSON.stringify(body);
  const url = isAuthRoute ? `${authUrl}/${route}` : `${baseUrl}/${route}`;
  const bearerToken = getIdToken();

  const config = {
    headers: {
      ...baseHeaders,
      Authorization: `Bearer ${bearerToken}`,
    },
    body: stringifiedBody,
    method: 'DELETE'
  };


  return fetch(url, config);
}

/**
* Method uses get method with api
* @param {string} route
* @param {boolean} isAuthRoute
* @return {Promise}
*/
const get = (route: string, isAuthRoute = false): Promise<Response> => {
  const url = isAuthRoute ? `${authUrl}/${route}` : `${baseUrl}/${route}`;
  const bearerToken = getIdToken();

  const headers = {
    headers: {
      ...baseHeaders,
      Authorization: `Bearer ${bearerToken}`,
    },
    method: 'GET',
  };

  return fetch(url, headers);
}


/**
* Method uses post method with api
* @param {string} route
* @return {Promise}
*/
const post = (
  route: string,
  body = {},
  isAuthRoute = false
) : Promise<Response> => {
  const url = isAuthRoute ? `${authUrl}/${route}` : `${baseUrl}/${route}`;
  const bearerToken = getIdToken();
  const stringifiedBody = JSON.stringify(body);

  const headers = {
    headers: {
      ...baseHeaders,
      Authorization: `Bearer ${bearerToken}`,
    },
    body: stringifiedBody,
    method: 'POST'
  };


  return fetch(url, headers);
}

/**
* Method uses put method with api
* @param {string} route
* @return {Promise}
*/
const put = (route: string, body = {}, isAuthRoute = false) : Promise<Response> => {
  const url = isAuthRoute ? `${authUrl}/${route}` : `${baseUrl}/${route}`;
  const bearerToken = getIdToken();
  const stringifiedBody = JSON.stringify(body);

  const config = {
    headers: {
      ...baseHeaders,
      Authorization: `Bearer ${bearerToken}`,
    },
    body: stringifiedBody,
    method: 'put'
  };


  return fetch(url, config);
}


/**
* Method uses get method with api
* @param {string} route
* @return {Promise}
*/
const getWellKnown = (route: string): Promise<Response> => {
  const baseUrl = `${window.location.protocol}//${window.location.hostname}/core/v1`;
  return fetch(`${baseUrl}/discover`, { headers: baseHeaders, method: 'GET' })
  .then(response => response.json())
  .then(data => {
    const url = `${data.auth_service}/${route}`;
    const config = {
      headers: {
        ...baseHeaders,
      },
      method: 'GET',
    };

    return fetch(url, config);
  });
}


const HossApi = {
  del,
  get,
  getWellKnown,
  post,
  put,
}

export {
  del,
  get,
  getWellKnown,
  post,
  put,
}




export default HossApi;
