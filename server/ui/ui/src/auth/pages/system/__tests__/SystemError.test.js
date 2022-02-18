// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// components
import SystemError from '../SystemError';

test('renders learn react link', () => {
  render(
    <Router>
      <SystemError />
    </Router>
  );
  const linkElement = screen.getByText(/SystemError/i);
  expect(linkElement).toBeInTheDocument();
});
