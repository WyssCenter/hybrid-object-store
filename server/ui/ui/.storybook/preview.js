import { addDecorator } from '@storybook/react';
import { withTests } from '@storybook/addon-jest';
import React from "react";
import { MemoryRouter } from "react-router";

import '../src/styles/critical.scss';


import results from '../.jest-test-results.json';

let modalRoot = document.createElement("div")
modalRoot.setAttribute("id", "modal")
document.querySelector("body").appendChild(modalRoot);

const user = JSON.stringify({ id_token: 'id_token' });
localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);


addDecorator(
  withTests({
    results,
    filesExt: '((\\.specs?)|(\\.tests?))?(\\.ts)?(\\.tsx)?$',
  })
);

addDecorator(story => {
    const render = story();

    const {pathname} = render.props.location ?  render.props.location : { pathname: '/'};
    return (
      <MemoryRouter initialEntries={[pathname]}>
        {story()}
      </MemoryRouter>
    )
  }
  );

export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  }
}
