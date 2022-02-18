// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// components
import Authenticating from '../Authenticating';

test('renders learn react link', () => {
  render(
    <Router>
      <Authenticating />
    </Router>
  );
  const linkElement = screen.getByText(/Authenticating/i);
  expect(linkElement).toBeInTheDocument();
});
