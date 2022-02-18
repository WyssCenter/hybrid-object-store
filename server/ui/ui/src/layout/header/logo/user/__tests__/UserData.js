
const context = {
  config: {
    userAuthManager: {
      getUser: () => new Promise((resolve) => {
        resolve({
          profile: {
            name: 'admin',
          }
        })
      })
    }
  }
}


export default context;
