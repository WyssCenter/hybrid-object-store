import React from 'react';
import { render, screen } from '@testing-library/react';
import Error from '../Error';

describe('Error Create Modal', ()=> {
  it('Renders as expected', () => {
    render(
      <Error
        errorMessage="This is an error"
        name="namespace-12"
        send={jest.fn}
      />
    );
    const errorElement = screen.getByText('Error Creating namespace-12');
    expect(errorElement).toBeInTheDocument();
  });
})
