// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import Loader from '../Loader';

describe('Loader', () => {
  it('Renders Loader', () => {
    render(<Loader nested={false} />);

    const loader = screen.getByTestId('loader');
    expect(loader.className).toBe('Loader');
  });
});
