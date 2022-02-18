// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import Tooltip from '../Tooltip';


describe('Tooltip', () => {
  it('renders learn react link', () => {
    render(<Tooltip />);
    const linkElement = screen.getByRole('presentation');
    expect(linkElement).toBeInTheDocument();
  });
});
