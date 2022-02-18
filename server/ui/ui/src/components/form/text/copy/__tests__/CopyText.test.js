// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
// components
import CopyText from '../CopyText';

describe('CopyText', () => {
  it('Renders button', () => {
    render(
      <CopyText
        text="copythis"
      />
    );
    const buttonElement = screen.getByRole('button');
    expect(buttonElement).toBeInTheDocument();
  });


  it('Renders text', () => {
    render(
      <CopyText
        text="copythis"
      />
    );
    expect(screen.getByRole('textbox').value).toBe('copythis');
  });


  it('Test Click', () => {
    global.document.execCommand = jest.fn();

    render(
      <CopyText
        text="copythis"
      />
    );

    userEvent.click(screen.getByRole('button'))
    expect(global.document.execCommand).toHaveBeenCalledTimes(1);
  });

});
