// vendor
import React from 'react';
import { cleanup, render, screen, fireEvent } from '@testing-library/react';
// components
import InputText from '../InputText';

describe('InputText', () => {
  afterEach(cleanup)

  it('Renders InputText', () => {
    render(
      <InputText
        css="small"
        inputRef={{ current: <div />}}
        label="Add User"
        placeholder="Username"
        updateValue={jest.fn()}
      />
    );
    const labelElement = screen.getByText('Add User');
    const textboxElement = screen.getByRole('textbox');
    expect(labelElement).toBeInTheDocument();
    expect(textboxElement.placeholder).toBe('Username');
  });



  it('Text updates', () => {
    const updateValue = jest.fn();
    render(
      <InputText
        css="small"
        inputRef={{ current: <div />}}
        label="Add User"
        placeholder="Username"
        updateValue={updateValue}
      />
    );
    const textboxElement = screen.getByRole('textbox');
    fireEvent.keyUp(textboxElement, { target: { value: 'this is the input' }})

    expect(textboxElement.value).toBe('this is the input');
    expect(updateValue).toHaveBeenCalledTimes(1);
  });
})
