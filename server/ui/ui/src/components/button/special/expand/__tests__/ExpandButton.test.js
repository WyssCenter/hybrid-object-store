// vendor
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
// components
import ExpandButton from '../ExpandButton';

describe('ExpandButton', () => {
  it('Renders button', () => {
    render(
      <ExpandButton
        click={jest.fn}
        isExpanded={false}
        text="Cancel"
      />
    );
    const linkElement = screen.getByText(/Cancel/);
    expect(linkElement).toBeInTheDocument();
  });


  it('Renders button disabled', () => {
    render(
      <ExpandButton
        click={jest.fn}
        disabled
        isExpanded={false}
        text="Cancel"
      />
    );
    expect(screen.getByRole('button')).toHaveAttribute('disabled')
  });


  it('Test Click', () => {
    const click = jest.fn();

    render(
      <ExpandButton
        click={click}
        isExpanded={true}
        text="Cancel"
      />
    );

    userEvent.click(screen.getByRole('button'))
    expect(click).toHaveBeenCalledTimes(1);
  });


  it('Test Click disabled', () => {
    const click = jest.fn();

    render(
      <ExpandButton
        click={click}
        disabled={true}
        text="Cancel"
      />
    );

    userEvent.click(screen.getByRole('button'))
    expect(click).toHaveBeenCalledTimes(0);
  });
});
