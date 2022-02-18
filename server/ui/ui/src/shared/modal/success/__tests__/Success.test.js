import React from 'react';
import { render, screen } from '@testing-library/react';
import Success from '../Success';


describe('Success Create Modal', () => {
  it('renders learn react link', () => {
    render(
      <Success
        modalType="namespace"
        name="namespace-12"
      />
    );
    const linkElement = screen.getByText(/was successfully created/);
    expect(linkElement).toBeInTheDocument();
  });
})
