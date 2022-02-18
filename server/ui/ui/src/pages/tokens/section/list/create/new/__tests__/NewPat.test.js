// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event'
// components
import NewPat from '../NewPat';

const pat = {
  token: "r4nD0mStr1ngF0rTh3T0k3n",
  id: 2,
}

describe('NewPat', () => {
  beforeAll(() => {
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    global.document.execCommand = jest.fn();

    document.body.appendChild(modalElement)
  })

  it('Renders value in input', () => {
    render(
      <NewPat
        dismissPat={jest.fn()}
        pat={pat}
      />
    );
    const inputElement = screen.getByRole('textbox');
    expect(inputElement.value).toBe('r4nD0mStr1ngF0rTh3T0k3n');
  });

  it('Clicking copy copies to clipboard', () => {

    render(
      <NewPat
        dismissPat={jest.fn()}
        pat={pat}
      />
    );
    const buttons = screen.queryAllByRole(/button/);
    userEvent.click(buttons[0]);
    expect(global.document.execCommand).toHaveBeenCalledTimes(1);
  });


  it('Clicking dismiss shows tooltip', () => {

    render(
      <NewPat
        dismissPat={jest.fn()}
        pat={pat}
      />
    );
    const buttons = screen.queryAllByRole(/button/);
    userEvent.click(buttons[1]);

    const queryText = screen.getByText(/Are you sure/i)
    expect(queryText).toBeInTheDocument(1);
  });


  it('Clicking cancel keeps input value intact', () => {

    render(
      <NewPat
        dismissPat={jest.fn()}
        pat={pat}
      />
    );
    const buttons = screen.queryAllByRole(/button/);
    userEvent.click(buttons[1]);

    const newButtons = screen.queryAllByRole(/button/);
    userEvent.click(newButtons[2]);

    const inputElement = screen.getByRole('textbox');
    expect(inputElement.value).toBe('r4nD0mStr1ngF0rTh3T0k3n');
  });


  it('Clicking confirm clears input value', async() => {
    const dismissPat = jest.fn();
    render(
      <NewPat
        dismissPat={dismissPat}
        pat={pat}
      />
    );
    const buttons = screen.queryAllByRole(/button/);
    userEvent.click(buttons[1]);

    const newButtons = screen.queryAllByRole(/button/);
    await userEvent.click(newButtons[3]);

    expect(dismissPat).toHaveBeenCalledTimes(1);
  });

})
