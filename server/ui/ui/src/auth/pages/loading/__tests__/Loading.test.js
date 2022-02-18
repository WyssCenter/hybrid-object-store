// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// components
import Loading from '../Loading';

test('renders learn react link', () => {
  render(
    <Router>
      <Loading />
    </Router>
  );
  const linkElement = screen.getByText(/Loading/i);
  expect(linkElement).toBeInTheDocument();
});
