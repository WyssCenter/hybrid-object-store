// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
// components
import WarningButton from '../WarningButton';

describe('WarningButton', () => {
  it('Renders button', () => {
    render(
      <WarningButton
        click={jest.fn}
        disabled={false}
        text="Cancel"
      />
    );
    const linkElement = screen.getByText(/Cancel/);
    expect(linkElement).toBeInTheDocument();
  });


  it('Renders button disabled', () => {
    render(
      <WarningButton
        click={jest.fn}
        disabled={true}
        text="Cancel"
      />
    );
    expect(screen.getByRole('button')).toHaveAttribute('disabled')
  });


  it('Test Click', () => {
    const click = jest.fn();

    render(
      <WarningButton
        click={click}
        disabled={false}
        text="Cancel"
      />
    );

    userEvent.click(screen.getByRole('button'))
    expect(click).toHaveBeenCalledTimes(1);
  });


  it('Test Click disabled', () => {
    const click = jest.fn();

    render(
      <WarningButton
        click={click}
        disabled={true}
        text="Cancel"
      />
    );

    userEvent.click(screen.getByRole('button'))
    expect(click).toHaveBeenCalledTimes(0);
  });
});
