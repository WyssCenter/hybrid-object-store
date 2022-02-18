// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// components
import Error from '../Error';

test('renders learn react link', () => {
  render(
    <Router>
      <Error />
    </Router>
  );
  const linkElement = screen.getByText(/Error/i);
  expect(linkElement).toBeInTheDocument();
});
