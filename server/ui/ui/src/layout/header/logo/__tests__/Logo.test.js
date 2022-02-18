import React from 'react';
import { render, screen } from '@testing-library/react';
import Logo from '../Logo';

test('renders learn react link', () => {
  render(<Logo />);
  const linkElement = screen.getByAltText(/Gigantum HOS Logo/i);
  expect(linkElement).toBeInTheDocument();
});
